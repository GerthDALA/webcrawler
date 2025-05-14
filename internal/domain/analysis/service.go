package analysis

import (
	"context"

	"github.com/gerthdala/webcrawler/internal/domain/content"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
)

// TextVectorizer generates vector embeddings for text
type TextVectorizer interface {
	// Vectorize generates a vector embedding for text
	Vectorize(ctx context.Context, text string) result.Result[[]float32]
}

// TopicModeler performs topic modeling on text
type TopicModeler interface {
	// ExtractTopics extracts topics from text
	ExtractTopics(ctx context.Context, text string, numTopics int) result.Result[[]content.Topic]
	
	// TrainModel trains the topic model
	TrainModel(ctx context.Context, texts []string) result.Result[bool]
	
	// GetTopicKeywords gets keywords for a topic
	GetTopicKeywords(ctx context.Context, topicID int, numKeywords int) result.Result[[]string]
}

// NamedEntityRecognizer performs named entity recognition
type NamedEntityRecognizer interface {
	// ExtractEntities extracts named entities from text
	ExtractEntities(ctx context.Context, text string) result.Result[[]content.NamedEntity]
}

// ContentClassifier classifies content
type ContentClassifier interface {
	// Classify classifies content
	Classify(ctx context.Context, text string) result.Result[content.ContentType]
	
	// GetConfidence gets the confidence of classification
	GetConfidence(ctx context.Context, text string, contentType content.ContentType) result.Result[float64]
}

// TextSummarizer summarizes text
type TextSummarizer interface {
	// Summarize summarizes text
	Summarize(ctx context.Context, text string, maxLength int) result.Result[string]
}

// KeywordExtractor extracts keywords from text
type KeywordExtractor interface {
	// ExtractKeywords extracts keywords from text
	ExtractKeywords(ctx context.Context, text string, numKeywords int) result.Result[[]string]
}

// LanguageDetector detects language of text
type LanguageDetector interface {
	// DetectLanguage detects the language of text
	DetectLanguage(ctx context.Context, text string) result.Result[string]
}

// ReadabilityAnalyzer analyzes readability of text
type ReadabilityAnalyzer interface {
	// AnalyzeReadability analyzes readability of text
	AnalyzeReadability(ctx context.Context, text string) result.Result[float64]
	
	// CountWords counts words in text
	CountWords(ctx context.Context, text string) result.Result[int]
	
	// CountSentences counts sentences in text
	CountSentences(ctx context.Context, text string) result.Result[int]
}

// SimilarityCalculator calculates similarity between content
type SimilarityCalculator interface {
	// CalculateSimilarity calculates similarity between two pieces of content
	CalculateSimilarity(ctx context.Context, a, b []float32) result.Result[float64]
	
	// FindMostSimilar finds the most similar content to the given content
	FindMostSimilar(ctx context.Context, embedding []float32, embeddings [][]float32, limit int) result.Result[[]int]
}

// AnalysisService orchestrates content analysis
type AnalysisService struct {
	vectorizer      TextVectorizer
	topicModeler    TopicModeler
	entityRecognizer NamedEntityRecognizer
	classifier      ContentClassifier
	summarizer      TextSummarizer
	keywordExtractor KeywordExtractor
	languageDetector LanguageDetector
	readabilityAnalyzer ReadabilityAnalyzer
	similarityCalculator SimilarityCalculator
}

// NewAnalysisService creates a new AnalysisService
func NewAnalysisService(
	vectorizer TextVectorizer,
	topicModeler TopicModeler,
	entityRecognizer NamedEntityRecognizer,
	classifier ContentClassifier,
	summarizer TextSummarizer,
	keywordExtractor KeywordExtractor,
	languageDetector LanguageDetector,
	readabilityAnalyzer ReadabilityAnalyzer,
	similarityCalculator SimilarityCalculator,
) *AnalysisService {
	return &AnalysisService{
		vectorizer:      vectorizer,
		topicModeler:    topicModeler,
		entityRecognizer: entityRecognizer,
		classifier:      classifier,
		summarizer:      summarizer,
		keywordExtractor: keywordExtractor,
		languageDetector: languageDetector,
		readabilityAnalyzer: readabilityAnalyzer,
		similarityCalculator: similarityCalculator,
	}
}