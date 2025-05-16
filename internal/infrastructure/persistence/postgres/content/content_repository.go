package content

import (
	"context"
	"fmt"
	"time"

	"github.com/gerthdala/webcrawler/internal/domain/content"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ContentRepository implements content.ContentRepository using PostgreSQL
type ContentRepository struct {
	db *gorm.DB
}

// NewContentRepository creates a new ContentRepository
func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{
		db: db,
	}
}

// Save stores a Content
func (r *ContentRepository) Save(ctx context.Context, contentData *content.Content) result.Result[*content.Content] {
	tx := r.db.WithContext(ctx)
	model := ContentModelFromDomain(contentData)

	if err := tx.Create(model).Error; err != nil {
		return result.Err[*content.Content](fmt.Errorf("failed to save Content: %w", err))
	}

	// Save associated entities
	for _, entity := range contentData.NamedEntities {
		entityModel := NamedEntityModelFromDomain(entity, contentData.ID)
		if err := tx.Create(entityModel).Error; err != nil {
			return result.Err[*content.Content](fmt.Errorf("failed to save NamedEntity: %w", err))
		}
	}

	for _, topic := range contentData.Topics {
		topicModel := TopicModelFromDomain(topic, contentData.ID)
		if err := tx.Create(topicModel).Error; err != nil {
			return result.Err[*content.Content](fmt.Errorf("failed to save Topic: %w", err))
		}
	}

	return result.Ok(contentData)
}

// FindByID finds a Content by its ID
func (r *ContentRepository) FindByID(ctx context.Context, id uuid.UUID) result.Result[*content.Content] {
	tx := r.db.WithContext(ctx)
	var model ContentModel

	if err := tx.Where("id = ?", id).First(&model).Error; err != nil {
		return result.Err[*content.Content](fmt.Errorf("failed to find Content by ID: %w", err))
	}

	// Get domain object
	contentObject := model.ToDomain()

	// Load associated entities
	var entityModels []NamedEntityModel
	if err := tx.Where("content_id = ?", id).Find(&entityModels).Error; err != nil {
		return result.Err[*content.Content](fmt.Errorf("failed to load NamedEntities: %w", err))
	}

	for _, entityModel := range entityModels {
		contentObject.NamedEntities = append(contentObject.NamedEntities, entityModel.ToDomain())
	}

	var topicModels []TopicModel
	if err := tx.Where("content_id = ?", id).Find(&topicModels).Error; err != nil {
		return result.Err[*content.Content](fmt.Errorf("failed to load Topics: %w", err))
	}

	for _, topicModel := range topicModels {
		contentObject.Topics = append(contentObject.Topics, topicModel.ToDomain())
	}

	return result.Ok(contentObject)
}

// FindByURL finds a Content by its URL
func (r *ContentRepository) FindByURL(ctx context.Context, url string) result.Result[*content.Content] {
	tx := r.db.WithContext(ctx)
	var model ContentModel

	if err := tx.Where("url = ?", url).First(&model).Error; err != nil {
		return result.Err[*content.Content](fmt.Errorf("failed to find Content by URL: %w", err))
	}

	return r.FindByID(ctx, model.ID)
}

// FindByTopic finds Content by topic
func (r *ContentRepository) FindByTopic(ctx context.Context, topic string, limit int) result.Result[[]content.Content] {
	tx := r.db.WithContext(ctx)
	var contentIDs []uuid.UUID

	// Find content IDs with the given topic
	if err := tx.Model(&TopicModel{}).
		Where("name LIKE ?", "%"+topic+"%").
		Distinct("content_id").
		Limit(limit).
		Pluck("content_id", &contentIDs).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find Content by topic: %w", err))
	}

	if len(contentIDs) == 0 {
		return result.Ok([]content.Content{})
	}

	// Find the contents
	var models []ContentModel
	if err := tx.Where("id IN ?", contentIDs).Find(&models).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find Content by IDs: %w", err))
	}

	// Convert to domain objects
	contents := make([]content.Content, len(models))
	for i, model := range models {
		content := model.ToDomain()

		// Load associated entities (optimized)
		var entityModels []NamedEntityModel
		tx.Where("content_id = ?", content.ID).Find(&entityModels)
		for _, entityModel := range entityModels {
			content.NamedEntities = append(content.NamedEntities, entityModel.ToDomain())
		}

		var topicModels []TopicModel
		tx.Where("content_id = ?", content.ID).Find(&topicModels)
		for _, topicModel := range topicModels {
			content.Topics = append(content.Topics, topicModel.ToDomain())
		}

		contents[i] = *content
	}

	return result.Ok(contents)
}

// FindByEntityType finds Content by entity type
func (r *ContentRepository) FindByEntityType(ctx context.Context, entityType content.EntityType, limit int) result.Result[[]content.Content] {
	tx := r.db.WithContext(ctx)
	var contentIDs []uuid.UUID

	// Find content IDs with the given entity type
	if err := tx.Model(&NamedEntityModel{}).
		Where("type = ?", string(entityType)).
		Distinct("content_id").
		Limit(limit).
		Pluck("content_id", &contentIDs).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find Content by entity type: %w", err))
	}

	if len(contentIDs) == 0 {
		return result.Ok([]content.Content{})
	}

	// Find the contents
	var models []ContentModel
	if err := tx.Where("id IN ?", contentIDs).Find(&models).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find Content by IDs: %w", err))
	}

	// Convert to domain objects with minimal loading
	contents := make([]content.Content, len(models))
	for i, model := range models {
		contents[i] = *model.ToDomain()
	}

	return result.Ok(contents)
}

// FindByContentType finds Content by content type
func (r *ContentRepository) FindByContentType(ctx context.Context, contentType content.ContentType, limit int) result.Result[[]content.Content] {
	tx := r.db.WithContext(ctx)
	var models []ContentModel

	if err := tx.Where("classification = ?", string(contentType)).
		Order("created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find Content by content type: %w", err))
	}

	// Convert to domain objects with minimal loading
	contents := make([]content.Content, len(models))
	for i, model := range models {
		contents[i] = *model.ToDomain()
	}

	return result.Ok(contents)
}

// FindSimilar finds Content similar to the given content ID
func (r *ContentRepository) FindSimilar(ctx context.Context, contentID uuid.UUID, limit int) result.Result[[]content.Content] {
	tx := r.db.WithContext(ctx)
	var similarIDs []uuid.UUID

	// Find similar content IDs
	if err := tx.Model(&SimilarContentModel{}).
		Where("content_id = ?", contentID).
		Order("similarity_score DESC").
		Limit(limit).
		Pluck("similar_to_id", &similarIDs).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find similar Content IDs: %w", err))
	}

	if len(similarIDs) == 0 {
		return result.Ok([]content.Content{})
	}

	// Find the contents
	var models []ContentModel
	if err := tx.Where("id IN ?", similarIDs).Find(&models).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find Content by IDs: %w", err))
	}

	// Convert to domain objects with minimal loading
	contents := make([]content.Content, len(models))
	for i, model := range models {
		contents[i] = *model.ToDomain()
	}

	return result.Ok(contents)
}

// FindNearest finds Content nearest to the given vector embedding
func (r *ContentRepository) FindNearest(ctx context.Context, embedding []float32, limit int) result.Result[[]content.Content] {
	tx := r.db.WithContext(ctx)
	var models []ContentModel

	// Using PostgreSQL's vector similarity search with pgvector extension
	// This assumes the pgvector extension is installed and the vector_embedding column is a vector type
	if err := tx.Raw(`
		SELECT * FROM contents
		ORDER BY vector_embedding <-> ?
		LIMIT ?
	`, pq.Float32Array(embedding), limit).Scan(&models).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to find nearest Content: %w", err))
	}

	// Convert to domain objects with minimal loading
	contents := make([]content.Content, len(models))
	for i, model := range models {
		contents[i] = *model.ToDomain()
	}

	return result.Ok(contents)
}

// Search searches Content by text
func (r *ContentRepository) Search(ctx context.Context, query string, limit int) result.Result[[]content.Content] {
	tx := r.db.WithContext(ctx)
	var models []ContentModel

	// Using PostgreSQL's full-text search
	if err := tx.Where("to_tsvector('english', title || ' ' || text) @@ plainto_tsquery('english', ?)", query).
		Order("created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return result.Err[[]content.Content](fmt.Errorf("failed to search Content: %w", err))
	}

	// Convert to domain objects with minimal loading
	contents := make([]content.Content, len(models))
	for i, model := range models {
		contents[i] = *model.ToDomain()
	}

	return result.Ok(contents)
}

// CountByContentType counts Content by content type
func (r *ContentRepository) CountByContentType(ctx context.Context, contentType content.ContentType) result.Result[int] {
	tx := r.db.WithContext(ctx)
	var count int64

	if err := tx.Model(&ContentModel{}).
		Where("classification = ?", string(contentType)).
		Count(&count).Error; err != nil {
		return result.Err[int](fmt.Errorf("failed to count Content by content type: %w", err))
	}

	return result.Ok(int(count))
}

// DeleteOlderThan deletes Content older than the given duration
func (r *ContentRepository) DeleteOlderThan(ctx context.Context, days int) result.Result[int] {
	tx := r.db.WithContext(ctx)
	cutoff := time.Now().AddDate(0, 0, -days)

	// Delete associated entities first
	var contentIDs []uuid.UUID
	if err := tx.Model(&ContentModel{}).
		Where("created_at < ?", cutoff).
		Pluck("id", &contentIDs).Error; err != nil {
		return result.Err[int](fmt.Errorf("failed to get old Content IDs: %w", err))
	}

	if len(contentIDs) > 0 {
		if err := tx.Where("content_id IN ?", contentIDs).Delete(&NamedEntityModel{}).Error; err != nil {
			return result.Err[int](fmt.Errorf("failed to delete NamedEntities: %w", err))
		}

		if err := tx.Where("content_id IN ?", contentIDs).Delete(&TopicModel{}).Error; err != nil {
			return result.Err[int](fmt.Errorf("failed to delete Topics: %w", err))
		}

		if err := tx.Where("content_id IN ? OR similar_to_id IN ?", contentIDs, contentIDs).Delete(&SimilarContentModel{}).Error; err != nil {
			return result.Err[int](fmt.Errorf("failed to delete SimilarContents: %w", err))
		}
	}

	// Delete contents
	resultD := tx.Where("created_at < ?", cutoff).Delete(&ContentModel{})
	if resultD.Error != nil {
		return result.Err[int](fmt.Errorf("failed to delete old Contents: %w", resultD.Error))
	}

	return result.Ok(int(result.RowsAffected))
}
