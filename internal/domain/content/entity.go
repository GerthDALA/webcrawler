package content

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeArctile  ContentType = "article"
	ContentTypeBlog     ContentType = "blog"
	ContentTypeDoc      ContentType = "documentation"
	ContentTypeProduct  ContentType = "product"
	ContentTypeHomePage ContentType = "homepage"
	ContentTypeOther    ContentType = "other"
)

type Content struct {
	ID               uuid.UUID
	URL              string
	Title            string
	Text             string
	HTML             string
	Summary          string
	Keywords         []string
	NamedEntities    []NamedEntity
	Classification   ContentType
	Language         string
	ReadabilityScore float64
	WordCount        int
	SentenceCount    int
	VectorEmbedding  []float32
	Topics           []Topic
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewContent creates a new Content entity
func NewContent(url, title, text, html string) *Content {
	return &Content{
		ID:        uuid.New(),
		URL:       url,
		Title:     title,
		Text:      text,
		HTML:      html,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c *Content) AddNamedEntities(entities []NamedEntity) {

	c.NamedEntities = entities
	c.UpdatedAt = time.Now()

}

func (c *Content) SetVectorEmbedding(embedding []float32) {

	c.VectorEmbedding = embedding
	c.UpdatedAt = time.Now()
}

func (c *Content) AddTopics(topics []Topic) {
	c.Topics = topics
	c.UpdatedAt = time.Now()
}

// SetClassification sets the content classification
func (c *Content) SetClassification(contentType ContentType) {
	c.Classification = contentType
	c.UpdatedAt = time.Now()
}

// SetSummary sets the content summary
func (c *Content) SetSummary(summary string) {
	c.Summary = summary
	c.UpdatedAt = time.Now()
}

// SetKeywords sets the content keywords
func (c *Content) SetKeywords(keywords []string) {
	c.Keywords = keywords
	c.UpdatedAt = time.Now()
}

// SetLanguage sets the content language
func (c *Content) SetLanguage(language string) {
	c.Language = language
	c.UpdatedAt = time.Now()
}

// SetReadabilityScore sets the content readability score
func (c *Content) SetReadabilityScore(score float64) {
	c.ReadabilityScore = score
	c.UpdatedAt = time.Now()
}

// SetWordCount sets the content word count
func (c *Content) SetWordCount(count int) {
	c.WordCount = count
	c.UpdatedAt = time.Now()
}

// SetSentenceCount sets the content sentence count
func (c *Content) SetSentenceCount(count int) {
	c.SentenceCount = count
	c.UpdatedAt = time.Now()
}

type EntityType string

const (
	EntityTypePerson       EntityType = "person"
	EntityTypeOrganization EntityType = "organization"
	EntityTypeLocation     EntityType = "location"
	EntityTypeDate         EntityType = "date"
	EntityTypeProduct      EntityType = "product"
	EntityTypeEvent        EntityType = "event"
	EntityTypeOther        EntityType = "other"
)

type NamedEntity struct {
	ID        uuid.UUID
	Text      string
	Type      string
	Count     int
	Positions []int
}

type Topic struct {
	ID         uuid.UUID
	Name       string
	Keywords   []string
	Confidence float64
}

// NewTopic creates a new Topic
func NewTopic(name string, keywords []string, confidence float64) Topic {
	return Topic{
		ID:         uuid.New(),
		Name:       name,
		Keywords:   keywords,
		Confidence: confidence,
	}
}

type SimilarContent struct {
	ContentID       uuid.UUID
	SimilarToID     uuid.UUID
	SimilarityScore float64
	CreatedAt       time.Time
}

// NewSimilarContent creates a new SimilarContent
func NewSimilarContent(contentID, similarToID uuid.UUID, similarityScore float64) SimilarContent {
	return SimilarContent{
		ContentID:       contentID,
		SimilarToID:     similarToID,
		SimilarityScore: similarityScore,
		CreatedAt:       time.Now(),
	}
}
