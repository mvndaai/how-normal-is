package domain

import (
	"sort"
	"strings"
	"unicode"
)

type Trait struct {
	Name  string `json:"name"`
	Group string `json:"group,omitempty"`
}

func NormalizeTraitName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	var b strings.Builder
	lastDash := false
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
			lastDash = false
		case unicode.IsSpace(r) || r == '-' || r == '_':
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	return out
}

func SuggestTraits(input string, known []Trait, limit int) []Trait {
	needle := NormalizeTraitName(input)
	if needle == "" || limit <= 0 {
		return nil
	}
	matches := make([]Trait, 0, limit)
	for _, trait := range known {
		norm := NormalizeTraitName(trait.Name)
		if norm == needle || strings.Contains(norm, needle) || strings.HasPrefix(norm, needle) {
			matches = append(matches, trait)
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return NormalizeTraitName(matches[i].Name) < NormalizeTraitName(matches[j].Name)
	})
	if len(matches) > limit {
		return matches[:limit]
	}
	return matches
}

func ValidateExclusiveTraitGroups(selected []Trait) (conflictingGroups []string) {
	groupCounts := map[string]int{}
	for _, t := range selected {
		g := strings.TrimSpace(strings.ToLower(t.Group))
		if g == "" {
			continue
		}
		groupCounts[g]++
	}
	for group, c := range groupCounts {
		if c > 1 {
			conflictingGroups = append(conflictingGroups, group)
		}
	}
	sort.Strings(conflictingGroups)
	return conflictingGroups
}
