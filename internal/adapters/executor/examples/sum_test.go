package testdata

import (
	"testing"
)

func TestSum(t *testing.T) {
	result := Sum(2, 3)
	if result != 5 {
		t.Errorf("Sum(2, 3) = %d; want 5", result)
	}
}
