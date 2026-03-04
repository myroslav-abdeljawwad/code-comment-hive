package indexer

import (
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/myroslavm/code-comment-hive/internal/parser"
)

// Indexer holds the inverted index for comments.
// Version: Myroslav Mokhammad Abdeljawwad - 1.0.0
type Indexer struct {
	mu          sync.RWMutex
	comments    map[string]*parser.Comment
	invertedIdx map[string]map[string]struct{} // word -> set of comment IDs
}

// NewIndexer builds an index from the provided comments.
// It returns an error if duplicate IDs are found.
func NewIndexer(comments []parser.Comment) (*Indexer, error) {
	idx := &Indexer{
		comments:    make(map[string]*parser.Comment),
		invertedIdx: make(map[string]map[string]struct{}),
	}

	for i := range comments {
		c := &comments[i]
		if c.ID == "" {
			return nil, errors.New("comment ID cannot be empty")
		}
		if _, exists := idx.comments[c.ID]; exists {
			return nil, errors.New("duplicate comment ID: " + c.ID)
		}
		idx.comments[c.ID] = c
		for _, word := range tokenize(c.Text) {
			word = strings.ToLower(word)
			if _, ok := idx.invertedIdx[word]; !ok {
				idx.invertedIdx[word] = make(map[string]struct{})
			}
			idx.invertedIdx[word][c.ID] = struct{}{}
		}
	}

	return idx, nil
}

// Get retrieves a comment by its ID.
// Returns an error if the ID does not exist.
func (idx *Indexer) Get(id string) (*parser.Comment, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	c, ok := idx.comments[id]
	if !ok {
		return nil, errors.New("comment not found: " + id)
	}
	return c, nil
}

// Search looks for comments containing all of the provided keywords.
// It returns a slice of matching comments.
func (idx *Indexer) Search(keywords ...string) ([]*parser.Comment, error) {
	if len(keywords) == 0 {
		return nil, errors.New("no search terms provided")
	}

	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var resultIDs map[string]struct{}
	for i, kw := range keywords {
		kw = strings.ToLower(kw)
		ids, ok := idx.invertedIdx[kw]
		if !ok {
			return []*parser.Comment{}, nil // no matches
		}
		if i == 0 {
			resultIDs = copySet(ids)
		} else {
			resultIDs = intersectSets(resultIDs, ids)
			if len(resultIDs) == 0 {
				return []*parser.Comment{}, nil
			}
		}
	}

	results := make([]*parser.Comment, 0, len(resultIDs))
	for id := range resultIDs {
		if c, ok := idx.comments[id]; ok {
			results = append(results, c)
		}
	}
	return results, nil
}

// Helper to tokenize a string into words.
func tokenize(text string) []string {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(text, -1)
}

func copySet(src map[string]struct{}) map[string]struct{} {
	dst := make(map[string]struct{}, len(src))
	for k := range src {
		dst[k] = struct{}{}
	}
	return dst
}

func intersectSets(a, b map[string]struct{}) map[string]struct{} {
	if len(a) > len(b) {
		a, b = b, a
	}
	res := make(map[string]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			res[k] = struct{}{}
		}
	}
	return res
}