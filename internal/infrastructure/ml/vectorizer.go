package ml

import (
	"sync"

	"github.com/james-bowman/nlp"
)

//TextVectorizer implements content.TextVectorizer interface
// It uses TF-IDF model under the hood and normalizes output to unut length

type TextVectorizer struct {
	vectorizer  *nlp.TfidfVectorizer
	initOnce    sync.Once
	mu          sync.RWMutex
	dimensions  int
	minDocFreq  int
	maxFeatures int
}

// Config for building a TextVectorizer.
type TextVectorizerConfig struct {
	// If zero, automatically use vocabulary size.
	Dimensions int
	// Minimum number of documents a term must appear in.
	MinDocFreq int
	// Maximum number of terms (by highest TF-IDF) to keep.
	MaxFeatures int
}

// NewTextVectorizer constructs a configurable TF-IDF vectorizer.
func NewTextVectorizer(cfg TextVectorizerConfig) *TextVectorizer {
	// Pass options into the vectorizer for filtering and feature‐selection.
	opts := []nlp.Option{
		nlp.WithMinDocCount(cfg.MinDocFreq),
	}
	if cfg.MaxFeatures > 0 {
		opts = append(opts, nlp.WithMaxFeatures(cfg.MaxFeatures))
	}
	return &TextVectorizer{
		vectorizer:  nlp.NewTfidfVectorizer(opts...),
		dimensions:  cfg.Dimensions,
		minDocFreq:  cfg.MinDocFreq,
		maxFeatures: cfg.MaxFeatures,
	}
}
