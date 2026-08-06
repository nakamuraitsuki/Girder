package mock

import (
	"context"

	"girder/internal/app"
	"girder/internal/domain/entity"
)

type MockBlueprintService struct{}

func (MockBlueprintService) ApplyBlueprint(context.Context, app.ApplyBlueprintRequest) (entity.Blueprint, error) {
	return entity.Blueprint{}, nil
}

func (MockBlueprintService) ValidateBlueprint(context.Context, app.ValidateBlueprintRequest) (app.ValidationResponse, error) {
	return app.ValidationResponse{}, nil
}

func (MockBlueprintService) ImportBlueprint(context.Context, app.ImportBlueprintRequest) (entity.Blueprint, error) {
	return entity.Blueprint{}, nil
}

func (MockBlueprintService) ExportBlueprint(context.Context, app.ExportBlueprintRequest) (entity.Blueprint, error) {
	return entity.Blueprint{}, nil
}

func (MockBlueprintService) DestroyBlueprint(context.Context, app.DestroyBlueprintRequest) (entity.BlueprintID, error) {
	return "", nil
}
