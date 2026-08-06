package service

import (
	"context"

	app "girder/internal/app"
	"girder/internal/domain/entity"
)

type BlueprintServiceImpl struct {
	BlueprintDriver      any
	InfrastructureDriver any
	ComputeDriver        any
	NetworkDriver        any
}

func (BlueprintServiceImpl) ApplyBlueprint(context.Context, app.ApplyBlueprintRequest) (entity.Blueprint, error) {
	return entity.Blueprint{}, nil
}

func (BlueprintServiceImpl) ValidateBlueprint(context.Context, app.ValidateBlueprintRequest) (app.ValidationResponse, error) {
	return app.ValidationResponse{}, nil
}

func (BlueprintServiceImpl) ImportBlueprint(context.Context, app.ImportBlueprintRequest) (entity.Blueprint, error) {
	return entity.Blueprint{}, nil
}

func (BlueprintServiceImpl) ExportBlueprint(context.Context, app.ExportBlueprintRequest) (entity.Blueprint, error) {
	return entity.Blueprint{}, nil
}

func (BlueprintServiceImpl) DestroyBlueprint(context.Context, app.DestroyBlueprintRequest) (entity.BlueprintID, error) {
	return "", nil
}
