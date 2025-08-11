// repository/tag_repository.go
package repository

import (
	"context"
	"time"
	"write_base/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagRepositoryImpl struct {
	collection *mongo.Collection
}

type TagDTO struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Status    string    `bson:"status"`
	CreatedBy string    `bson:"created_by"`
	CreatedAt time.Time `bson:"created_at"`
}

func NewTagRepository(db *mongo.Database) domain.TagRepository {
	return &TagRepositoryImpl{
		collection: db.Collection("tags"),
	}
}

// Convert domain Tag to DTO
func toTagDTO(tag *domain.Tag) *TagDTO {
	return &TagDTO{
		ID:        tag.ID,
		Name:      tag.Name,
		Status:    string(tag.Status),
		CreatedBy: tag.CreatedBy,
		CreatedAt: tag.CreatedAt,
	}
}

// Convert DTO to domain Tag
func toDomainTag(dto *TagDTO) *domain.Tag {
	return &domain.Tag{
		ID:        dto.ID,
		Name:      dto.Name,
		Status:    domain.TagStatus(dto.Status),
		CreatedBy: dto.CreatedBy,
		CreatedAt: dto.CreatedAt,
	}
}

func (r *TagRepositoryImpl) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tagDTO := toTagDTO(tag)
	tagDTO.CreatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, tagDTO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrTagAlreadyExists
		}
		return nil, domain.ErrInternalServer
	}
	
	// Return the created tag with updated timestamps
	tag.CreatedAt = tagDTO.CreatedAt
	return tag, nil
}

func (r *TagRepositoryImpl) Update(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	filter := bson.M{"_id": tag.ID}
	update := bson.M{"$set": bson.M{
		"status": string(tag.Status),
	}}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	return tag, nil
}

func (r *TagRepositoryImpl) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	var dto TagDTO
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrTagNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return toDomainTag(&dto), nil
}

func (r *TagRepositoryImpl) GetByID(ctx context.Context, ID string) (*domain.Tag, error) {
	var dto TagDTO
	err := r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrTagNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return toDomainTag(&dto), nil
}

func (r *TagRepositoryImpl) List(ctx context.Context, filter domain.TagFilter) ([]domain.Tag, error) {
	query := bson.M{}
	if filter.Status != "" {
		query["status"] = string(filter.Status)
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var tags []domain.Tag
	for cursor.Next(ctx) {
		var dto TagDTO
		if err := cursor.Decode(&dto); err != nil {
			return nil, domain.ErrInternalServer
		}
		tags = append(tags, *toDomainTag(&dto))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, domain.ErrInternalServer
	}
	
	return tags, nil
}

func (r *TagRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return domain.ErrInternalServer
	}
	return nil
}