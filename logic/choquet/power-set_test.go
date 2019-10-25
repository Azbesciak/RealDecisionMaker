package choquet

import (
	"testing"
)

func TestPowerset(t *testing.T) {
	result := PowerSet([]string{"A", "B", "C"})
	if len(result) != 7 {
		t.Errorf("invalid result lengt, got %d", len(result))
	}
}
