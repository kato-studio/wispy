package atomicstyle_test

import (
	"log"
	"testing"

	"github.com/kato-studio/wispy/atomicstyle"
)

func TestServer(t *testing.T) {

	var result = atomicstyle.Begin()

	log.Printf("result \n %s", result)
}
