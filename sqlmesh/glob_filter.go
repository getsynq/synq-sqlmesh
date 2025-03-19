package sqlmesh

import (
	"path/filepath"
	"strings"

	doublestart "github.com/bmatcuk/doublestar/v4"
)

type GlobFilter interface {
	Match(string) (bool, error)
}

func NewExcludeEverythingGlobFilter() GlobFilter {
	return &GlobFilterImpl{
		includePatterns: []string{},
	}
}

func NewGlobFilter(includePattern string, excludePattern string) GlobFilter {

	var includePatterns []string
	var excludePatterns []string

	if len(includePattern) > 0 {
		includePatterns = strings.Split(includePattern, ",")
	}
	if len(excludePattern) > 0 {
		excludePatterns = strings.Split(excludePattern, ",")
	}

	return &GlobFilterImpl{
		includePatterns: includePatterns,
		excludePatterns: excludePatterns,
	}
}

type GlobFilterImpl struct {
	includePatterns []string
	excludePatterns []string
}

func (g GlobFilterImpl) Match(s string) (bool, error) {
	match := false
	for _, pattern := range g.includePatterns {
		ok, err := doublestart.Match(pattern, s)
		if err != nil {
			return false, err
		}
		if ok {
			match = true
			break
		}
	}
	if !match {
		return false, nil
	}

	for _, pattern := range g.excludePatterns {
		ok, err := filepath.Match(pattern, s)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}

	return true, nil
}
