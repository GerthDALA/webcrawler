package ml

import (
	"context"
	"fmt"
	"math"
	"sort"

	result "github.com/gerthdala/webcrawler/pkg/utils/result"
)

// SimilarityCalculator computes cosine similarities.
type SimilarityCalculator struct {}

func NewSimilarityCalculator() * SimilarityCalculator {
	return &SimilarityCalculator{}
}

func (sc *SimilarityCalculator) CalculateSimilarity(ctx context.Context, a, b []float32) result.Result[float64] {
	if len(a) != len(b) {
		return result.Err[float64](fmt.Errorf(
			"dimensions mismatch: len(a)=%d, len(b)=%d", len(a), len(b),
		))
	}

	var dot, normA, normB float64
	for i, ai := range a {
		bi := b[i]
		dot += float64(ai * bi)
		normA += float64(ai * ai)
		normB += float64(bi * bi)
	}
	if normA == 0 || normB == 0 {
		return result.Ok(0.0)
	}
	return result.Ok(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}

func (sc *SimilarityCalculator) FindMostSimilar(ctx context.Context, embedding []float32, embeddings [][]float32, limit int) result.Result[[]int] {
	n := len(embeddings)
	if n == 0 || limit <= 0 {
		return result.Ok([]int)
	}

	// Build a slice if (index, similarity)
	type pair struct {
		idx int
		sim float64
	}

	pairs := make([]pair, 0, n)

	for i, vec := range embeddings {
		if simRes := sc.CalculateSimilarity(ctx, embedding, vec); simRes.IsOk() {

			pairs = append(pairs, pair{i, simRes.Unwrap()})
		}
	}

	if limit < len(pairs) {
		sort.Slice(pairs, func(i, j int) bool  {
			return pairs[i].sim > pairs[j].sim
		})
		pairs = pairs[:limit]
	} else {
		// full sort
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].sim > pairs[j].sim
		})
	}

	resultIdxs := make([]int, len(pairs))
	
	for i, p := range pairs {
		resultIdxs[i] = p.idx
	}
	return result.Ok(resultIdxs)
}