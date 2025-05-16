package ml

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/gerthdala/webcrawler/internal/domain/content"
	result "github.com/gerthdala/webcrawler/pkg/utils/result"
	"github.com/jdkato/prose/v2"
)

// NamedEntityRecognizer implements analysis.NamedEntityRecognizer.
type NamedEntityRecognizer struct{}

// NewNamedEntityRecognizer constructs a NamedEntityRecognizer.
func NewNamedEntityRecognizer() *NamedEntityRecognizer {
    return &NamedEntityRecognizer{}
}

// ExtractEntities runs NER on text and returns a slice of NamedEntity,
// each with all zero-based start positions and sorted by earliest occurrence.
func (ner *NamedEntityRecognizer) ExtractEntities(
	ctx context.Context,
	text string,
) result.Result[[]content.NamedEntity] {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return result.Err[[]content.NamedEntity](fmt.Errorf("NER initialization failed: %w", err))
	}

	labelMap := map[string]content.EntityType{
		"PERSON":  content.EntityTypePerson,
		"ORG":     content.EntityTypeOrganization,
		"GPE":     content.EntityTypeLocation,
		"LOC":     content.EntityTypeLocation,
		"DATE":    content.EntityTypeDate,
		"TIME":    content.EntityTypeDate,
		"PRODUCT": content.EntityTypeProduct,
		"EVENT":   content.EntityTypeEvent,
	}

	type info struct {
		typ       content.EntityType
		positions []int
	}

	entitiesByText := make(map[string]*info)

	for _, ent := range doc.Entities() {
		if isFilteredEntity(ent.Text) {
			continue
		}

		etype := labelMap[ent.Label]
		if etype == "" {
			etype = content.EntityTypeOther
		}

		// Collect all Unicode-safe positions
		for _, pos := range findAllPositionsRunes(text, ent.Text) {
			if _, ok := entitiesByText[ent.Text]; !ok {
				entitiesByText[ent.Text] = &info{typ: etype}
			}
			entitiesByText[ent.Text].positions = append(entitiesByText[ent.Text].positions, pos)
		}
	}

	results := make([]content.NamedEntity, 0, len(entitiesByText))
	for txt, inf := range entitiesByText {
		if len(inf.positions) == 0 {
			continue
		}
		ne := content.NewNamedEntity(txt, inf.typ, inf.positions)
		results = append(results, ne)
	}

	// Sort by first occurrence
	sort.Slice(results, func(i, j int) bool {
		return results[i].Positions[0] < results[j].Positions[0]
	})

	return result.Ok(results)
}


// findAllPositions returns every start index of substr in s.
func finAllPositions(s, substr string) [] int {
	var positions []int
	offset := 0
	for {
		idx := strings.Index(s[offset:], substr)
		if idx < 0 {
			break
		}
		absolute := offset + idx
		positions = append(positions, absolute)
		offset = absolute + len(substr)
	}
	return positions
}
// findAllPositionsRunes finds all start positions of substr in text (rune-aware)
func findAllPositionsRunes(text, substr string) []int {
	var positions []int
	runes := []rune(text)
	target := []rune(substr)
	textLen, targetLen := len(runes), len(target)

	for i := 0; i <= textLen-targetLen; i++ {
		match := true
		for j := 0; j < targetLen; j++ {
			if runes[i+j] != target[j] {
				match = false
				break
			}
		}
		if match {
			positions = append(positions, i)
			i += targetLen - 1 // move to next non-overlapping start
		}
	}
	return positions
}


func isFilteredEntity(entity string) bool {
	entity = strings.TrimSpace(strings.ToLower(entity))
	stopwords := map[string]bool{
		"the": true, "at": true, "of": true, "a": true, "in": true,
	}
	return len(entity) <= 2 || stopwords[entity]
}

