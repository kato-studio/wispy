package static

import (
	"fmt"
	"kato-studio/katoengine/lib/engine/template"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func WritePageFile(name string, path string, pageContent string) {
	// ensure folders exists for page
	os.MkdirAll("./build/"+path, os.ModePerm)

	if path[:5] == "./view" {
		path = path[5:]
	}

	// write page content to file
	pageFile, err := os.Create("./build/" + path + "/" + strings.ReplaceAll(name, "+page.kato", "index.html"))
	//
	if err != nil {
		utils.Fatal(fmt.Sprint(err))
	}
	defer pageFile.Close()
	pageFile.WriteString(pageContent)
}

func RenderFolder(path string, data gjson.Result) {
	pages, err := os.ReadDir(path)
	if err != nil {
		utils.Fatal(fmt.Sprint(err))
	}
	// setup component storage/cache
	components := store.GlobalByteMap()

	for _, page := range pages {
		// if page is a folder
		if page.IsDir() {
			RenderFolder(path+"/"+page.Name(), data)
		} else if strings.Contains(page.Name(), ".kato") {
			pageBytes, err := os.ReadFile(path + "/" + page.Name())

			if err != nil {
				utils.Fatal(fmt.Sprint(err))
			}

			// render page
			fileData := template.Render(string(pageBytes), data, components.Store())
			// write page to file
			WritePageFile(page.Name(), path, fileData)
		}
	}
}

func LoadAllComponents(componentsDir string) {
	components := store.GlobalByteMap()

	foldersAndFiles, err := os.ReadDir(componentsDir)
	if err != nil {
		utils.Fatal(fmt.Sprint(err))
	}

	utils.Debug(fmt.Sprintf("Loading components from %s", componentsDir))
	utils.Debug(fmt.Sprintf("Found %d components", len(foldersAndFiles)))

	for _, thing := range foldersAndFiles {
		// if page is a folder
		if thing.IsDir() {
		} else if strings.Contains(thing.Name(), ".kato") {
			componentBytes, err := os.ReadFile(componentsDir + "/" + thing.Name())
			if err != nil {
				utils.Fatal(fmt.Sprint(err))
			}
			hasSet := components.SafeSet(thing.Name(), componentBytes)
			if hasSet == nil {
				utils.Fatal(fmt.Sprintf("Component the name {%s} already exists", thing.Name()))
			}
		}
	}
}
