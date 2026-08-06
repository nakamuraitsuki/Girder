package entity

import (
	"fmt"
	"strings"
	
	"girder/internal/domain/value"
)

type Route struct {
	ID          RouteID
	Name        string
	Destination value.CIDR
	NextHop     value.IPv4
	Metric      int
	Enabled     bool
}

func NewRoute(name string, destination value.CIDR, nextHop value.IPv4, metric int) (Route, error) {
	if strings.TrimSpace(name) == "" {
		return Route{}, fmt.Errorf("%w: route name", ErrEmptyName)
	}
	if destination.IsZero() || nextHop.IsZero() {
		return Route{}, fmt.Errorf("%w: route target", ErrInvalidResource)
	}
	if metric < 0 {
		return Route{}, fmt.Errorf("%w: route metric", ErrInvalidResource)
	}
	return Route{ID: NewRouteID(), Name: strings.TrimSpace(name), Destination: destination, NextHop: nextHop, Metric: metric, Enabled: true}, nil
}

func (r *Route) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: route name", ErrEmptyName)
	}
	r.Name = strings.TrimSpace(name)
	return nil
}

func (r *Route) UpdateDestination(destination value.CIDR) error {
	if destination.IsZero() {
		return fmt.Errorf("%w: route destination", ErrInvalidResource)
	}
	r.Destination = destination
	return nil
}

func (r *Route) UpdateNextHop(nextHop value.IPv4) error {
	if nextHop.IsZero() {
		return fmt.Errorf("%w: route next hop", ErrInvalidResource)
	}
	r.NextHop = nextHop
	return nil
}

func (r *Route) SetMetric(metric int) error {
	if metric < 0 {
		return fmt.Errorf("%w: route metric", ErrInvalidResource)
	}
	r.Metric = metric
	return nil
}

func (r *Route) Enable() {
	r.Enabled = true
}

func (r *Route) Disable() {
	r.Enabled = false
}
