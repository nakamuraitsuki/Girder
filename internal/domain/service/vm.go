package service

import (
	"fmt"

	"girder/internal/domain/entity"
)

type VMAllocationService struct{}

func (VMAllocationService) Assign(vm *entity.VM, node *entity.Node) error {
	if vm == nil || node == nil {
		return fmt.Errorf("%w: vm or node", entity.ErrInvalidResource)
	}
	if !node.CanHost(vm.ResourceRequest()) {
		return fmt.Errorf("%w: vm allocation", entity.ErrInsufficientCapacity)
	}
	if err := node.Reserve(vm.ResourceRequest()); err != nil {
		return err
	}
	return vm.AssignNode(node.ID)
}

func (VMAllocationService) Detach(vm *entity.VM, node *entity.Node) error {
	if vm == nil || node == nil {
		return fmt.Errorf("%w: vm or node", entity.ErrInvalidResource)
	}
	if err := node.Release(vm.ResourceRequest()); err != nil {
		return err
	}
	vm.DetachNode()
	return nil
}