package backend

import (
	"context"

	"girder/api"
)

type MetadataAPIImpl struct {
	Store api.MetadataStore
}

var _ api.MetadataStore = (*MetadataAPIImpl)(nil)

func NewMetadataAPIImpl(store api.MetadataStore) *MetadataAPIImpl {
	return &MetadataAPIImpl{Store: store}
}

func (b *MetadataAPIImpl) UpsertResource(ctx context.Context, request api.UpsertResourceRequest) (api.Resource, error) {
	return b.Store.UpsertResource(ctx, request)
}

func (b *MetadataAPIImpl) DeleteResource(ctx context.Context, request api.DeleteResourceRequest) error {
	return b.Store.DeleteResource(ctx, request)
}

func (b *MetadataAPIImpl) GetResource(ctx context.Context, request api.GetResourceRequest) (api.Resource, error) {
	return b.Store.GetResource(ctx, request)
}

func (b *MetadataAPIImpl) ListResources(ctx context.Context, request api.ListResourcesRequest) ([]api.Resource, error) {
	return b.Store.ListResources(ctx, request)
}

func (b *MetadataAPIImpl) UpsertRelation(ctx context.Context, request api.UpsertRelationRequest) (api.Relation, error) {
	return b.Store.UpsertRelation(ctx, request)
}

func (b *MetadataAPIImpl) DeleteRelation(ctx context.Context, request api.DeleteRelationRequest) error {
	return b.Store.DeleteRelation(ctx, request)
}

func (b *MetadataAPIImpl) ListRelations(ctx context.Context, request api.ListRelationsRequest) ([]api.Relation, error) {
	return b.Store.ListRelations(ctx, request)
}

func (b *MetadataAPIImpl) SetAnnotation(ctx context.Context, request api.SetAnnotationRequest) (api.Annotation, error) {
	return b.Store.SetAnnotation(ctx, request)
}

func (b *MetadataAPIImpl) DeleteAnnotation(ctx context.Context, request api.DeleteAnnotationRequest) error {
	return b.Store.DeleteAnnotation(ctx, request)
}
