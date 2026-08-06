package mock

import (
	"context"

	"girder/internal/app"
	"girder/internal/domain/entity"
)

type MockInfrastructureService struct{}

func (MockInfrastructureService) RegisterNode(context.Context, app.RegisterNodeRequest) (entity.Node, error) {
	return entity.Node{}, nil
}

func (MockInfrastructureService) RemoveNode(context.Context, app.RemoveNodeRequest) (entity.NodeID, error) {
	return "", nil
}

func (MockInfrastructureService) GetNode(context.Context, app.GetNodeRequest) (entity.Node, error) {
	return entity.Node{}, nil
}

func (MockInfrastructureService) ListNodes(context.Context) ([]entity.Node, error) {
	return nil, nil
}
