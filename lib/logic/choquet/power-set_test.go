package choquet

import (
	"strings"
	"testing"
)

func TestPowerset(t *testing.T) {
	result := *PowerSet([]string{"A", "B", "C"})
	if len(result) != 7 {
		t.Errorf("invalid result lengt, got %d", len(result))
	}
	expected := map[string]bool{
		"A":     true,
		"B":     true,
		"C":     true,
		"A,B":   true,
		"A,C":   true,
		"B,C":   true,
		"A,B,C": true,
	}
	for _, r := range result {
		key := strings.Join(r, ",")
		_, ok := expected[key]
		if !ok {
			t.Errorf("key '%v' not found in result", key)
		}
	}
}
