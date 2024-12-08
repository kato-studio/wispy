package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kato-studio/wispy/engine"
)

func TestParseTemplates(t *testing.T) {

	files, err := os.ReadDir(engine.ROOT_DIR)

	if err != nil {
		t.Fatalf("failed to read sites directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			siteDir := filepath.Join(engine.ROOT_DIR, file.Name())
			templates, err := template.ParseGlob(filepath.Join(siteDir, "*"+engine.EXT))
			if err != nil {
				t.Errorf("failed to parse templates for site %s: %v", file.Name(), err)
				continue
			}

			log.Printf("Parsed templates for site %s: %v", file.Name(), templates.DefinedTemplates())
		}
	}
}
