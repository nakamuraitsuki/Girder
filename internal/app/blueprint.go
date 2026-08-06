package app

import (
	"context"

	"girder/internal/domain/entity"
)

type BlueprintService interface {
	ApplyBlueprint(context.Context, ApplyBlueprintRequest) (entity.Blueprint, error)
	ValidateBlueprint(context.Context, ValidateBlueprintRequest) (ValidationResponse, error)
	ImportBlueprint(context.Context, ImportBlueprintRequest) (entity.Blueprint, error)
	ExportBlueprint(context.Context, ExportBlueprintRequest) (entity.Blueprint, error)
	DestroyBlueprint(context.Context, DestroyBlueprintRequest) (entity.BlueprintID, error)
}
