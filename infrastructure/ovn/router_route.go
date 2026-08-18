package ovn

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

var ErrLogicalRouterStaticRouteAlreadyExists = fmt.Errorf("logical router static route already exists")
var ErrLogicalRouterStaticRouteNotFound = fmt.Errorf("logical router static route not found")


// AddStaticRoute adds a static route to the logical router.
//
// Girder does not enforce uniqueness on (ip_prefix, nexthop) beyond what
// OVN itself allows; calling this twice with the same prefix will add a
// second route, which OVN permits (ECMP). Use SetDefaultRoute if you want
// idempotent replace-style behavior for the default route.
func (c *Client) AddStaticRoute(
	ctx context.Context,
	router *LogicalRouter,
	route *LogicalRouterStaticRoute,
) (*LogicalRouterStaticRoute, error) {
	if router == nil {
		return nil, errors.New("logical router is nil")
	}
	if route == nil {
		return nil, errors.New("logical router static route is nil")
	}
	if route.IPPrefix == "" {
		return nil, errors.New("ip_prefix is required")
	}
	if route.Nexthop == "" {
		return nil, errors.New("nexthop is required")
	}

	route.UUID = uuid.NewString()

	ops, err := c.buildAddStaticRouteOps(router, route)
	if err != nil {
		return nil, fmt.Errorf("failed to build add static route op: %w", err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to add static route %q to router %q: %w", route.IPPrefix, router.Name, err)
	}

	return route, nil
}

// buildAddStaticRouteOps builds the Create+Mutate ops for adding route to
// router, without transacting. route.UUID must already be set by the caller.
func (c *Client) buildAddStaticRouteOps(router *LogicalRouter, route *LogicalRouterStaticRoute) ([]ovsdb.Operation, error) {
	createOps, err := c.nb.Create(route)
	if err != nil {
		return nil, fmt.Errorf("failed to build create op for static route %q: %w", route.IPPrefix, err)
	}

	mutateOps, err := c.nb.Where(router).Mutate(router,
		model.Mutation{
			Field:   &router.StaticRoutes,
			Mutator: "insert",
			Value:   []string{route.UUID},
		})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to build mutate op attaching static route %q to router %q: %w",
			route.IPPrefix, router.Name, err,
		)
	}

	return append(createOps, mutateOps...), nil
}

// buildRemoveStaticRouteOps builds the Delete+Mutate ops for removing
// routeUUID from router, without transacting.
func (c *Client) buildRemoveStaticRouteOps(router *LogicalRouter, routeUUID string) ([]ovsdb.Operation, error) {
	route := &LogicalRouterStaticRoute{UUID: routeUUID}

	deleteOps, err := c.nb.Where(route).Delete()
	if err != nil {
		return nil, fmt.Errorf("failed to build delete op for static route %q: %w", routeUUID, err)
	}

	mutateOps, err := c.nb.Where(router).Mutate(router,
		model.Mutation{
			Field:   &router.StaticRoutes,
			Mutator: "delete",
			Value:   []string{routeUUID},
		})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to build mutate op detaching static route %q from router %q: %w",
			routeUUID, router.Name, err,
		)
	}

	return append(deleteOps, mutateOps...), nil
}

// ListStaticRoutes returns all static routes currently attached to the router.
func (c *Client) ListStaticRoutes(ctx context.Context, router *LogicalRouter) ([]LogicalRouterStaticRoute, error) {
	if router == nil {
		return nil, errors.New("logical router is nil")
	}

	var all []LogicalRouterStaticRoute
	if err := c.nb.List(ctx, &all); err != nil {
		return nil, fmt.Errorf("failed to list logical router static routes: %w", err)
	}

	var routes []LogicalRouterStaticRoute
	for _, r := range all {
		if slices.Contains(router.StaticRoutes, r.UUID) {
			routes = append(routes, r)
		}
	}

	return routes, nil
}

// RemoveStaticRoute removes a static route (identified by its UUID) from the router.
func (c *Client) RemoveStaticRoute(ctx context.Context, router *LogicalRouter, routeUUID string) error {
	if router == nil {
		return errors.New("logical router is nil")
	}

	ops, err := c.buildRemoveStaticRouteOps(router, routeUUID)
	if err != nil {
		return fmt.Errorf("failed to build remove static route op: %w", err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return fmt.Errorf("failed to remove static route %q from router %q: %w", routeUUID, router.Name, err)
	}

	return nil
}

// ReplaceStaticRoute idempotently replaces the router's static route for
// route.IPPrefix, in a single transaction: any existing static route with
// the same ip_prefix is removed and the new one is inserted atomically.
//
// This has no notion of "default route" or any other special-cased prefix;
// it is a general replace-by-prefix operation. OVN itself has no concept
// of a "default route" distinct from an ordinary static route with
// ip_prefix "0.0.0.0/0" -- that meaning belongs to the API layer, not here.
func (c *Client) ReplaceStaticRoute(
	ctx context.Context,
	router *LogicalRouter,
	route *LogicalRouterStaticRoute,
) (*LogicalRouterStaticRoute, error) {
	if router == nil {
		return nil, errors.New("logical router is nil")
	}
	if route == nil {
		return nil, errors.New("logical router static route is nil")
	}
	if route.IPPrefix == "" {
		return nil, errors.New("ip_prefix is required")
	}
	if route.Nexthop == "" {
		return nil, errors.New("nexthop is required")
	}

	existing, err := c.ListStaticRoutes(ctx, router)
	if err != nil {
		return nil, fmt.Errorf("failed to list existing static routes: %w", err)
	}

	var ops []ovsdb.Operation

	for _, r := range existing {
		if r.IPPrefix != route.IPPrefix {
			continue
		}
		removeOps, err := c.buildRemoveStaticRouteOps(router, r.UUID)
		if err != nil {
			return nil, fmt.Errorf("failed to build remove op for existing route %q: %w", r.IPPrefix, err)
		}
		ops = append(ops, removeOps...)
		break
	}

	route.UUID = uuid.NewString()

	addOps, err := c.buildAddStaticRouteOps(router, route)
	if err != nil {
		return nil, fmt.Errorf("failed to build add op for route %q: %w", route.IPPrefix, err)
	}
	ops = append(ops, addOps...)

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to replace static route %q on router %q: %w", route.IPPrefix, router.Name, err)
	}

	return route, nil
}