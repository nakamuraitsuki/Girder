package service

import (
	"context"

	app "girder/internal/app"
	"girder/internal/domain/entity"
)

type InfrastructureServiceImpl struct {
	NodeDriver any
}

func (InfrastructureServiceImpl) RegisterNode(context.Context, app.RegisterNodeRequest) (entity.Node, error) {
	return entity.Node{}, nil
}

func (InfrastructureServiceImpl) RemoveNode(context.Context, app.RemoveNodeRequest) (entity.NodeID, error) {
	return "", nil
}

func (InfrastructureServiceImpl) GetNode(context.Context, app.GetNodeRequest) (entity.Node, error) {
	return entity.Node{}, nil
}

func (InfrastructureServiceImpl) ListNodes(context.Context) ([]entity.Node, error) {
	return nil, nil
}
