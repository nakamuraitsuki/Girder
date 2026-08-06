package app

import (
	"context"

	"girder/internal/domain/entity"
)

type InfrastructureService interface {
	RegisterNode(context.Context, RegisterNodeRequest) (entity.Node, error)
	RemoveNode(context.Context, RemoveNodeRequest) (entity.NodeID, error)
	GetNode(context.Context, GetNodeRequest) (entity.Node, error)
	ListNodes(context.Context) ([]entity.Node, error)
}
