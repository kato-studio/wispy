package atomicstyle_test

import (
	"log"
	"testing"

	"github.com/kato-studio/wispy/atomicstyle"
)

func TestBegin(t *testing.T) {
	result := atomicstyle.Begin()
	log.Printf("Generated CSS: \n%s", result)

	if result == "" {
		t.Errorf("Expected generated CSS, but got empty string")
	}
}
