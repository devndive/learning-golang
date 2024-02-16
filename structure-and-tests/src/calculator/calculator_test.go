package calculator

import (
	"testing"
)

func TestAdd(t *testing.T) {
	actual := Add(1, 2);

	if actual != 3 {
		t.Errorf("Expected 3, got %d", actual);
	}
}
