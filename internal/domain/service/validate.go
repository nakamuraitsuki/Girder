package service

import (
	"fmt"
	"sort"

	"girder/internal/domain/entity"
)

type ImageValidationService struct{}

func (ImageValidationService) Validate(image entity.Image) error {
	if image.ID == "" || image.Name == "" {
		return fmt.Errorf("%w: image", entity.ErrInvalidResource)
	}
	if image.SizeBytes < 0 {
		return fmt.Errorf("%w: image size", entity.ErrInvalidResource)
	}
	return nil
}

type RouteValidationService struct{}

func (RouteValidationService) Validate(route entity.Route) error {
	if route.ID == "" || route.Name == "" || route.Destination.IsZero() || route.NextHop.IsZero() {
		return fmt.Errorf("%w: route", entity.ErrInvalidResource)
	}
	return nil
}

type ACLValidationService struct{}

func (ACLValidationService) Validate(acl entity.ACL) error {
	if acl.ID == "" || acl.Name == "" {
		return fmt.Errorf("%w: acl", entity.ErrInvalidResource)
	}
	priorities := make([]int, 0, len(acl.Rules))
	for _, rule := range acl.Rules {
		priorities = append(priorities, rule.Priority)
	}
	sort.Ints(priorities)
	for i := 1; i < len(priorities); i++ {
		if priorities[i] == priorities[i-1] {
			return fmt.Errorf("%w: acl rule priority", entity.ErrDuplicateItem)
		}
	}
	return nil
}

type NetworkValidationService struct{}

func (NetworkValidationService) ValidatePort(port entity.Port) error {
	if port.ID == "" || port.Name == "" || port.MACAddress.IsZero() {
		return fmt.Errorf("%w: port", entity.ErrInvalidResource)
	}
	attachments := 0
	if port.AttachedVMID != nil {
		attachments++
	}
	if port.AttachedSwitchID != nil {
		attachments++
	}
	if port.AttachedRouterID != nil {
		attachments++
	}
	if attachments > 1 {
		return fmt.Errorf("%w: port attachment", entity.ErrInvalidRelation)
	}
	return nil
}

func (NetworkValidationService) ValidateSwitch(sw entity.Switch) error {
	if sw.ID == "" || sw.Name == "" {
		return fmt.Errorf("%w: switch", entity.ErrInvalidResource)
	}
	return nil
}

func (NetworkValidationService) ValidateRouter(router entity.Router) error {
	if router.ID == "" || router.Name == "" {
		return fmt.Errorf("%w: router", entity.ErrInvalidResource)
	}
	return nil
}

type BlueprintValidationService struct{}

func (BlueprintValidationService) Validate(blueprint entity.Blueprint) error {
	if blueprint.ID == "" || blueprint.Name == "" || blueprint.ProjectID == "" || blueprint.OwnerUserID == "" {
		return fmt.Errorf("%w: blueprint", entity.ErrInvalidResource)
	}
	return nil
}

type PermissionValidationService struct{}

func (PermissionValidationService) CanManageProject(user entity.User, project entity.Project) error {
	if user.ID == "" || project.ID == "" {
		return fmt.Errorf("%w: permission target", entity.ErrInvalidResource)
	}
	if user.ID == project.OwnerUserID {
		return nil
	}
	if _, ok := project.MemberIDs[user.ID]; ok {
		return nil
	}
	return fmt.Errorf("%w: user is not a project member", entity.ErrInvalidRelation)
}
