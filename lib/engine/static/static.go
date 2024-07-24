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
	page_file, err := os.Create("./build/" + path + "/" + strings.ReplaceAll(name, "+page.kato", "index.html"))
	//
	if err != nil {
		utils.Fatal(fmt.Sprint(err))
	}
	defer page_file.Close()
	page_file.WriteString(pageContent)
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
			page_bytes, err := os.ReadFile(path + "/" + page.Name())

			if err != nil {
				utils.Fatal(fmt.Sprint(err))
			}

			// render page
			file_data := template.Render(string(page_bytes), data, components.Store())
			// write page to file
			WritePageFile(page.Name(), path, file_data)
		}
	}
}

func LoadAllComponents(componentsDir string) {
	components := store.GlobalByteMap()
	current_path := componentsDir
	folders_and_files, err := os.ReadDir(componentsDir)
	if err != nil {
		utils.Fatal(fmt.Sprint(err))
	}

	for _, thing := range folders_and_files {
		// if page is a folder
		if thing.IsDir() {
			LoadAllComponents(current_path + thing.Name() + "/")
		} else if strings.Contains(thing.Name(), ".kato") {
			path := current_path + thing.Name()
			clean_path := strings.Replace(current_path + thing.Name(),"./view/components","@",1)
			component_bytes, err := os.ReadFile(path)
			utils.Print(clean_path)
			if err != nil {
				utils.Fatal(fmt.Sprint(err))
			}
			//
			has_set := components.SafeSet(clean_path, component_bytes)
			fmt.Println("Loading component: ", thing.Name())
			if has_set == nil {
				utils.Fatal(fmt.Sprintf("Component the name {%s} already exists", thing.Name()))
			}
		}
	}
}
