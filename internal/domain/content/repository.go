package content

import (
	"context"

	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/google/uuid"
)

// ContentRepository handles Content storage and retrieval
type ContentRepository interface {
	// Save stores a Content
	Save(ctx context.Context, content *Content) result.Result[*Content]

	// FindByID finds a Content by its ID
	FindByID(ctx context.Context, id uuid.UUID) result.Result[*Content]

	// FindByURL finds a Content by its URL
	FindByURL(ctx context.Context, url string) result.Result[*Content]

	// FindByTopic finds Content by topic
	FindByTopic(ctx context.Context, topic string, limit int) result.Result[[]Content]

	// FindByEntityType finds Content by entity type
	FindByEntityType(ctx context.Context, entityType EntityType, limit int) result.Result[[]Content]

	// FindByContentType finds Content by content type
	FindByContentType(ctx context.Context, contentType ContentType, limit int) result.Result[[]Content]

	// FindSimilar finds Content similar to the given content ID
	FindSimilar(ctx context.Context, contentID uuid.UUID, limit int) result.Result[[]Content]

	// FindNearest finds Content nearest to the given vector embedding
	FindNearest(ctx context.Context, embedding []float32, limit int) result.Result[[]Content]

	// Search searches Content by text
	Search(ctx context.Context, query string, limit int) result.Result[[]Content]

	// CountByContentType counts Content by content type
	CountByContentType(ctx context.Context, contentType ContentType) result.Result[int]

	// DeleteOlderThan deletes Content older than the given duration
	DeleteOlderThan(ctx context.Context, days int) result.Result[int]
}

// NamedEntityRepository handles NamedEntity storage and retrieval
type NamedEntityRepository interface {
	// Save stores a NamedEntity
	Save(ctx context.Context, entity NamedEntity, contentID uuid.UUID) result.Result[NamedEntity]

	// FindByID finds a NamedEntity by its ID
	FindByID(ctx context.Context, id uuid.UUID) result.Result[NamedEntity]

	// FindByText finds a NamedEntity by its text
	FindByText(ctx context.Context, text string) result.Result[[]NamedEntity]

	// FindByType finds NamedEntities by type
	FindByType(ctx context.Context, entityType EntityType, limit int) result.Result[[]NamedEntity]

	// FindMostFrequent finds the most frequent NamedEntities
	FindMostFrequent(ctx context.Context, entityType EntityType, limit int) result.Result[[]NamedEntity]

	// FindByContentID finds NamedEntities by content ID
	FindByContentID(ctx context.Context, contentID uuid.UUID) result.Result[[]NamedEntity]
}

// TopicRepository handles Topic storage and retrieval
type TopicRepository interface {
	// Save stores a Topic
	Save(ctx context.Context, topic Topic, contentID uuid.UUID) result.Result[Topic]

	// FindByID finds a Topic by its ID
	FindByID(ctx context.Context, id uuid.UUID) result.Result[Topic]

	// FindByName finds a Topic by its name
	FindByName(ctx context.Context, name string) result.Result[[]Topic]

	// FindMostConfident finds the most confident Topics
	FindMostConfident(ctx context.Context, limit int) result.Result[[]Topic]

	// FindByContentID finds Topics by content ID
	FindByContentID(ctx context.Context, contentID uuid.UUID) result.Result[[]Topic]

	// FindMostPopular finds the most popular Topics
	FindMostPopular(ctx context.Context, limit int) result.Result[[]Topic]
}

// SimilarContentRepository handles SimilarContent storage and retrieval
type SimilarContentRepository interface {
	// Save stores a SimilarContent
	Save(ctx context.Context, similarContent SimilarContent) result.Result[SimilarContent]

	// FindByContentID finds SimilarContent by content ID
	FindByContentID(ctx context.Context, contentID uuid.UUID, limit int) result.Result[[]SimilarContent]

	// FindBySimilarToID finds SimilarContent by similar to ID
	FindBySimilarToID(ctx context.Context, similarToID uuid.UUID, limit int) result.Result[[]SimilarContent]

	// FindMostSimilar finds the most similar content
	FindMostSimilar(ctx context.Context, limit int) result.Result[[]SimilarContent]

	// DeleteByContentID deletes SimilarContent by content ID
	DeleteByContentID(ctx context.Context, contentID uuid.UUID) result.Result[int]
}
