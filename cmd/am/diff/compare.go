package diff

import (
	"fmt"
	"sort"
)

type FieldChange struct {
	Field string
	From  interface{}
	To    interface{}
}

type DiffResult struct {
	Added   int
	Removed int
	Changed int
	Lines   []string
}

func diffObjects(from, to map[string]interface{}, compareFields []string) []FieldChange {
	var changes []FieldChange
	for _, field := range compareFields {
		fromVal := fmt.Sprintf("%v", from[field])
		toVal := fmt.Sprintf("%v", to[field])
		if fromVal != toVal {
			changes = append(changes, FieldChange{Field: field, From: from[field], To: to[field]})
		}
	}
	return changes
}

func compareResources(fromItems, toItems []map[string]interface{}, keyField string, compareFields []string) DiffResult {
	result := DiffResult{}

	fromMap := make(map[string]map[string]interface{}, len(fromItems))
	for _, item := range fromItems {
		key := fmt.Sprintf("%v", item[keyField])
		fromMap[key] = item
	}

	toMap := make(map[string]map[string]interface{}, len(toItems))
	for _, item := range toItems {
		key := fmt.Sprintf("%v", item[keyField])
		toMap[key] = item
	}

	// Collect added+changed keys from toMap (sorted for deterministic output)
	toKeys := make([]string, 0, len(toMap))
	for k := range toMap {
		toKeys = append(toKeys, k)
	}
	sort.Strings(toKeys)

	for _, key := range toKeys {
		toItem := toMap[key]
		fromItem, exists := fromMap[key]
		if !exists {
			result.Added++
			result.Lines = append(result.Lines, formatDiffLine("+", keyField, key, toItem))
			continue
		}
		changes := diffObjects(fromItem, toItem, compareFields)
		if len(changes) > 0 {
			result.Changed++
			for _, c := range changes {
				result.Lines = append(result.Lines, fmt.Sprintf("~ %-20s %-30s %v -> %v", key, c.Field, c.From, c.To))
			}
		}
	}

	// Collect removed keys from fromMap (sorted for deterministic output)
	fromKeys := make([]string, 0, len(fromMap))
	for k := range fromMap {
		fromKeys = append(fromKeys, k)
	}
	sort.Strings(fromKeys)

	for _, key := range fromKeys {
		if _, exists := toMap[key]; !exists {
			result.Removed++
			result.Lines = append(result.Lines, formatDiffLine("-", keyField, key, fromMap[key]))
		}
	}

	return result
}

func formatDiffLine(prefix string, _ string, key string, _ map[string]interface{}) string {
	return fmt.Sprintf("%s %-30s", prefix, key)
}
