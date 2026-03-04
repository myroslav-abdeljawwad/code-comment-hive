package search

import (
	"sort"
	"strings"

	"github.com/myproject/code-comment-hive/internal/indexer"
	"github.com/myproject/code-comment-hive/internal/parser"
)

// Result represents a single search hit in the comment index.
type Result struct {
	Repo    string
	File    string
	Line    int
	Snippet string // snippet of the comment text containing the query
}

// Searcher performs queries over an index created by internal/indexer.
// It is intentionally lightweight and does not maintain its own state.
type Searcher struct {
	idx *indexer.Indexer
}

// NewSearcher creates a new Searcher for the provided Indexer.
func NewSearcher(idx *indexer.Indexer) *Searcher {
	return &Searcher{idx: idx}
}

// Perform executes a case‑insensitive substring search over all stored comments.
// It returns up to maxResults hits, sorted deterministically by Repo/File/Line.
// If maxResults is zero or negative, all matches are returned.
func (s *Searcher) Perform(query string, maxResults int) ([]Result, error) {
	if s.idx == nil {
		return nil, ErrNilIndexer
	}
	if strings.TrimSpace(query) == "" {
		return nil, ErrEmptyQuery
	}

	all := s.idx.All()
	var hits []Result

	lowerQuery := strings.ToLower(query)
	for _, c := range all {
		text := c.Text
		if strings.Contains(strings.ToLower(text), lowerQuery) {
			snippet := snippetForMatch(text, query)
			hits = append(hits, Result{
				Repo:    c.Repo,
				File:    c.File,
				Line:    c.Line,
				Snippet: snippet,
			})
		}
	}

	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].Repo != hits[j].Repo {
			return hits[i].Repo < hits[j].Repo
		}
		if hits[i].File != hits[j].File {
			return hits[i].File < hits[j].File
		}
		return hits[i].Line < hits[j].Line
	})

	if maxResults > 0 && len(hits) > maxResults {
		hits = hits[:maxResults]
	}

	return hits, nil
}

// snippetForMatch returns a short excerpt of text surrounding the first occurrence
// of query. It is capped to 80 characters total for readability.
func snippetForMatch(text, query string) string {
	idx := strings.Index(strings.ToLower(text), strings.ToLower(query))
	if idx == -1 {
		return ""
	}
	start := idx - 20
	if start < 0 {
		start = 0
	}
	end := idx + len(query) + 20
	if end > len(text) {
		end = len(text)
	}
	snip := strings.TrimSpace(text[start:end])
	if len(snip) > 80 {
		snip = snip[:77] + "..."
	}
	return snip
}

// Errors returned by Searcher.
var (
	ErrNilIndexer = &SearchError{"indexer is nil"}
	ErrEmptyQuery = &SearchError{"query string is empty"}
)

// SearchError represents an error in the search package.
type SearchError struct {
	msg string
}

func (e *SearchError) Error() string { return e.msg }

// ---------------------------------------------------------------------------
// The following helper functions are used by tests and other packages.
// ---------------------------------------------------------------------------

// For convenience, expose a top‑level function that performs a search using
// an existing indexer. It is equivalent to NewSearcher(idx).Perform.
func Search(idx *indexer.Indexer, query string, maxResults int) ([]Result, error) {
	return NewSearcher(idx).Perform(query, maxResults)
}

// ---------------------------------------------------------------------------
// Example usage in README:
//
//   idx := indexer.New()
//   // ... populate idx ...
//   results, err := search.Search(idx, "TODO", 10)
//   if err != nil {
//       log.Fatal(err)
//   }
//   for _, r := range results {
//       fmt.Printf("%s:%d %q\n", r.File, r.Line, r.Snippet)
//   }
//
// ---------------------------------------------------------------------------

/*
Author: Myroslav Mokhammad Abdeljawwad
Version: 1.0.0 (search module)
*/