package service

import (
	"fmt"

	"girder/internal/domain/entity"
)

type NodeAllocationService struct{}

func (NodeAllocationService) Reserve(node *entity.Node, request entity.NodeAllocation) error {
	if node == nil {
		return fmt.Errorf("%w: node", entity.ErrInvalidResource)
	}
	return node.Reserve(request)
}

func (NodeAllocationService) Release(node *entity.Node, request entity.NodeAllocation) error {
	if node == nil {
		return fmt.Errorf("%w: node", entity.ErrInvalidResource)
	}
	return node.Release(request)
}

func (NodeAllocationService) SelectNode(nodes []entity.Node, request entity.NodeAllocation) (*entity.Node, error) {
	var selected *entity.Node
	for i := range nodes {
		node := &nodes[i]
		if !node.CanHost(request) {
			continue
		}
		if selected == nil {
			selected = node
			continue
		}
		current := selected.RemainingCapacity()
		candidate := node.RemainingCapacity()
		if candidate.CPU > current.CPU || candidate.MemoryMiB > current.MemoryMiB || candidate.DiskGiB > current.DiskGiB {
			selected = node
		}
	}
	if selected == nil {
		return nil, fmt.Errorf("%w: no eligible node", entity.ErrNotFound)
	}
	return selected, nil
}