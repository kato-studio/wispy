package structure

import (
	"crypto/sha256"
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"
)

type AssetType string

const (
	CSS AssetType = "CSS"
	JS            = "JS"
)

type Asset struct {
	Path        string
	Type        AssetType
	Priority    int  // Lower numbers load first
	Async       bool // JS only
	Defer       bool // JS only
	Module      bool // JS module type
	Nomodule    bool
	Content     string // For inline assets
	Preload     bool
	Integrity   string
	Crossorigin string
	Media       string
	IsInline    bool
	Condition   string // Conditional loading
}

type AssetRegistry struct {
	assets   map[AssetType][]*Asset
	seen     map[string]struct{} // For deduplication
	depGraph map[string][]string
}

func (r *AssetRegistry) Add(asset *Asset, dependencies ...string) error {
	// Initialize data structures if needed
	if r.assets == nil {
		r.assets = make(map[AssetType][]*Asset)
		r.seen = make(map[string]struct{})
	}

	// Validate asset
	if asset == nil {
		return fmt.Errorf("cannot add nil asset")
	}

	if !asset.IsInline && asset.Path == "" {
		return fmt.Errorf("external asset must have a path")
	}

	// Generate unique key
	key := r.generateAssetKey(asset)

	// Check for duplicates
	if _, exists := r.seen[key]; exists {
		fmt.Println("Dup Found skipping: ", key)
		return nil // No error, just skip
	}

	// Set defaults
	r.applyAssetDefaults(asset)

	// Handle dependencies
	if len(dependencies) > 0 {
		asset.Priority = r.calculateDependentPriority(asset, dependencies)
	}

	// Add to registry
	r.assets[asset.Type] = append(r.assets[asset.Type], asset)
	r.seen[key] = struct{}{}

	// Maintain separate dependency graph if needed
	if r.depGraph == nil {
		r.depGraph = make(map[string][]string)
	}
	r.depGraph[key] = dependencies

	return nil
}

// Helper methods
func (r *AssetRegistry) generateAssetKey(asset *Asset) string {
	if asset.IsInline {
		return fmt.Sprintf("inline:%s:%x", asset.Type, sha256.Sum256([]byte(asset.Content)))
	}
	return fmt.Sprintf("external:%s:%s", asset.Type, filepath.Clean(asset.Path))
}

func (r *AssetRegistry) applyAssetDefaults(asset *Asset) {
	if asset.Priority == 0 {
		switch asset.Type {
		case CSS:
			asset.Priority = 100
		case JS:
			if asset.Module {
				asset.Priority = 200
			} else {
				asset.Priority = 300
			}
		}
	}
}

func (r *AssetRegistry) calculateDependentPriority(asset *Asset, deps []string) int {
	minPriority := math.MaxInt32
	for _, dep := range deps {
		if existing, ok := r.findAssetByKey(dep); ok {
			if existing.Priority < minPriority {
				minPriority = existing.Priority
			}
		}
	}

	if minPriority != math.MaxInt32 {
		return minPriority - 10
	}
	return asset.Priority
}

func (r *AssetRegistry) findAssetByKey(key string) (*Asset, bool) {
	for _, assets := range r.assets {
		for _, asset := range assets {
			if r.generateAssetKey(asset) == key {
				return asset, true
			}
		}
	}
	return nil, false
}

func (r *AssetRegistry) Render(t AssetType) string {
	var output strings.Builder

	// Sort assets by priority
	sort.Slice(r.assets[t], func(i, j int) bool {
		return r.assets[t][i].Priority < r.assets[t][j].Priority
	})

	for _, asset := range r.assets[t] {
		if t == CSS {
			if asset.IsInline {
				output.WriteString("<style>")
				output.WriteString(asset.Content)
				output.WriteString("</style>")
			} else {
				output.WriteString(fmt.Sprintf(`<link href="%s" rel="stylesheet"  type="text/css">`, asset.Path))
			}
		} else if t == JS {
			attrs := ""
			if asset.Async {
				attrs += " async"
			}
			if asset.Defer {
				attrs += " defer"
			}
			if asset.Module {
				attrs += " type=\"module\""
			}
			//
			if asset.IsInline {
				output.WriteString(fmt.Sprintf(`<script %s>`, attrs))
				output.WriteString(asset.Content)
			} else {
				output.WriteString(fmt.Sprintf(`<script src="%s"%s>`, asset.Path, attrs))
			}
			output.WriteString("</script>")
		}
	}

	return output.String()
}
