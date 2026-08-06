package api

import "context"

type MetadataStore interface {
	UpsertResource(context.Context, UpsertResourceRequest) (Resource, error)
	DeleteResource(context.Context, DeleteResourceRequest) error
	GetResource(context.Context, GetResourceRequest) (Resource, error)
	ListResources(context.Context, ListResourcesRequest) ([]Resource, error)
	UpsertRelation(context.Context, UpsertRelationRequest) (Relation, error)
	DeleteRelation(context.Context, DeleteRelationRequest) error
	ListRelations(context.Context, ListRelationsRequest) ([]Relation, error)
	SetAnnotation(context.Context, SetAnnotationRequest) (Annotation, error)
	DeleteAnnotation(context.Context, DeleteAnnotationRequest) error
}

type UpsertResourceRequest struct {
	Resource Resource
}

type DeleteResourceRequest struct {
	ResourceID string
}

type GetResourceRequest struct {
	ResourceID string
}

type ListResourcesRequest struct {
	Kind   string
	Driver string
}

type UpsertRelationRequest struct {
	Relation Relation
}

type DeleteRelationRequest struct {
	RelationID string
}

type ListRelationsRequest struct {
	ResourceID string
}

type SetAnnotationRequest struct {
	Annotation Annotation
}

type DeleteAnnotationRequest struct {
	ResourceID string
	Key        string
}
