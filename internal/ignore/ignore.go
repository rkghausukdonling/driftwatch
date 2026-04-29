// Package ignore provides functionality for loading and evaluating
// .driftwatchignore rules to suppress known or accepted drift.
package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Rules holds a set of ignore patterns loaded from a .driftwatchignore file.
type Rules struct {
	patterns []pattern
}

type pattern struct {
	resourceType string // empty means wildcard
	idPrefix     string
}

// Load reads ignore rules from the given file path.
// Each non-blank, non-comment line must be in the form:
//
//	<type>/<id-prefix>   — match by type and id prefix
//	*/<id-prefix>        — match any type with the given id prefix
//	<type>/*             — match all resources of a given type
func Load(path string) (*Rules, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Rules{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var patterns []pattern
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "/", 2)
		if len(parts) != 2 {
			continue // malformed — skip
		}
		rType := strings.TrimSpace(parts[0])
		idPfx := strings.TrimSpace(parts[1])
		if rType == "*" {
			rType = ""
		}
		if idPfx == "*" {
			idPfx = ""
		}
		patterns = append(patterns, pattern{resourceType: rType, idPrefix: idPfx})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Rules{patterns: patterns}, nil
}

// Match reports whether the given resource type and id should be ignored.
func (r *Rules) Match(resourceType, id string) bool {
	for _, p := range r.patterns {
		typeMatch := p.resourceType == "" || p.resourceType == resourceType
		idMatch := p.idPrefix == "" || strings.HasPrefix(id, p.idPrefix)
		if typeMatch && idMatch {
			return true
		}
	}
	return false
}

// Empty returns true when no patterns are loaded.
func (r *Rules) Empty() bool {
	return len(r.patterns) == 0
}
