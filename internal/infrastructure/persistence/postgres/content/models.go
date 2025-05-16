package content

import (
	"time"

	"github.com/gerthdala/webcrawler/internal/domain/content"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ContentModel is the database model for Content
type ContentModel struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key"`
	URL              string         `gorm:"uniqueIndex;not null"`
	Title            string
	Text             string         `gorm:"type:text"`
	HTML             string         `gorm:"type:text"`
	Summary          string
	Keywords         pq.StringArray `gorm:"type:text[]"`
	Classification   string
	Language         string
	ReadabilityScore float64
	WordCount        int
	SentenceCount    int
	VectorEmbedding  pq.Float32Array `gorm:"type:vector(384)"`  // Adjust vector dimension as needed
	CreatedAt        time.Time       `gorm:"index;not null"`
	UpdatedAt        time.Time       `gorm:"not null"`
}

// TableName returns the table name for the Content model
func (ContentModel) TableName() string {
	return "contents"
}

// ToDomain converts ContentModel to domain Content
func (m *ContentModel) ToDomain() *content.Content {
	c := &content.Content{
		ID:               m.ID,
		URL:              m.URL,
		Title:            m.Title,
		Text:             m.Text,
		HTML:             m.HTML,
		Summary:          m.Summary,
		Keywords:         []string(m.Keywords),
		Classification:   content.ContentType(m.Classification),
		Language:         m.Language,
		ReadabilityScore: m.ReadabilityScore,
		WordCount:        m.WordCount,
		SentenceCount:    m.SentenceCount,
		VectorEmbedding:  []float32(m.VectorEmbedding),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
	return c
}

// FromDomain converts domain Content to ContentModel
func ContentModelFromDomain(c *content.Content) *ContentModel {
	return &ContentModel{
		ID:               c.ID,
		URL:              c.URL,
		Title:            c.Title,
		Text:             c.Text,
		HTML:             c.HTML,
		Summary:          c.Summary,
		Keywords:         pq.StringArray(c.Keywords),
		Classification:   string(c.Classification),
		Language:         c.Language,
		ReadabilityScore: c.ReadabilityScore,
		WordCount:        c.WordCount,
		SentenceCount:    c.SentenceCount,
		VectorEmbedding:  pq.Float32Array(c.VectorEmbedding),
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

// NamedEntityModel is the database model for NamedEntity
type NamedEntityModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key"`
	ContentID  uuid.UUID      `gorm:"type:uuid;index;not null"`
	Text       string         `gorm:"index;not null"`
	Type       string         `gorm:"index;not null"`
	Count      int            `gorm:"not null"`
	Positions  pq.Int32Array  `gorm:"type:int[]"`
}

// TableName returns the table name for the NamedEntity model
func (NamedEntityModel) TableName() string {
	return "named_entities"
}

// ToDomain converts NamedEntityModel to domain NamedEntity
func (m *NamedEntityModel) ToDomain() content.NamedEntity {
	positions := make([]int, len(m.Positions))
	for i, pos := range m.Positions {
		positions[i] = int(pos)
	}
	
	return content.NamedEntity{
		ID:        m.ID,
		Text:      m.Text,
		Type:      content.EntityType(m.Type),
		Count:     m.Count,
		Positions: positions,
	}
}

// FromDomain converts domain NamedEntity to NamedEntityModel
func NamedEntityModelFromDomain(e content.NamedEntity, contentID uuid.UUID) *NamedEntityModel {
	positions := make(pq.Int32Array, len(e.Positions))
	for i, pos := range e.Positions {
		positions[i] = int32(pos)
	}
	
	return &NamedEntityModel{
		ID:        e.ID,
		ContentID: contentID,
		Text:      e.Text,
		Type:      string(e.Type),
		Count:     e.Count,
		Positions: positions,
	}
}

// TopicModel is the database model for Topic
type TopicModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key"`
	ContentID  uuid.UUID      `gorm:"type:uuid;index;not null"`
	Name       string         `gorm:"index;not null"`
	Keywords   pq.StringArray `gorm:"type:text[]"`
	Confidence float64        `gorm:"not null"`
}

// TableName returns the table name for the Topic model
func (TopicModel) TableName() string {
	return "topics"
}

// ToDomain converts TopicModel to domain Topic
func (m *TopicModel) ToDomain() content.Topic {
	return content.Topic{
		ID:         m.ID,
		Name:       m.Name,
		Keywords:   []string(m.Keywords),
		Confidence: m.Confidence,
	}
}

// FromDomain converts domain Topic to TopicModel
func TopicModelFromDomain(t content.Topic, contentID uuid.UUID) *TopicModel {
	return &TopicModel{
		ID:         t.ID,
		ContentID:  contentID,
		Name:       t.Name,
		Keywords:   pq.StringArray(t.Keywords),
		Confidence: t.Confidence,
	}
}

// SimilarContentModel is the database model for SimilarContent
type SimilarContentModel struct {
	ContentID       uuid.UUID `gorm:"type:uuid;index;not null"`
	SimilarToID     uuid.UUID `gorm:"type:uuid;index;not null"`
	SimilarityScore float64   `gorm:"not null"`
	CreatedAt       time.Time `gorm:"index;not null"`
}

// TableName returns the table name for the SimilarContent model
func (SimilarContentModel) TableName() string {
	return "similar_contents"
}

// ToDomain converts SimilarContentModel to domain SimilarContent
func (m *SimilarContentModel) ToDomain() content.SimilarContent {
	return content.SimilarContent{
		ContentID:       m.ContentID,
		SimilarToID:     m.SimilarToID,
		SimilarityScore: m.SimilarityScore,
		CreatedAt:       m.CreatedAt,
	}
}

// FromDomain converts domain SimilarContent to SimilarContentModel
func SimilarContentModelFromDomain(s content.SimilarContent) *SimilarContentModel {
	return &SimilarContentModel{
		ContentID:       s.ContentID,
		SimilarToID:     s.SimilarToID,
		SimilarityScore: s.SimilarityScore,
		CreatedAt:       s.CreatedAt,
	}
}