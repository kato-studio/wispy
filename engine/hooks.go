package engine

import (
	"fmt"
	"os"
	"strings"
)

// -----===================-----
func fileClosure(dir_path string, ctx RenderCTX) error {
	renderedContents := RenderPage(dir_path, ctx)
	//
	if renderedContents == "" {
		fmt.Println("[warn]: Could not render the page: ", dir_path)
		return nil
	}
	output_dir := strings.Replace(dir_path, "./sites", "./.build/static", 1)
	output_dir = strings.Replace(output_dir, "+page.hstm", "", 1)
	//
	output_path := strings.Replace(dir_path, "./sites", "./.build/static", 1)
	output_path = strings.Replace(output_path, "+page.hstm", "index.html", 1)
	//
	dir_err := os.MkdirAll(output_dir, 0755)
	err := os.WriteFile(output_path, []byte(renderedContents), 0644)
	if err != nil {
		fmt.Println("[Error]: Could not write the file: ", output_path)
		if dir_err != nil {
			fmt.Println("[Error]: Could not create the directory: ", output_dir)
		} else {
			fmt.Println("[Error]: Could not write the file: ", output_path)
		}
		fmt.Println("----------")
		return err
	}
	return nil
}

func dirClosure(dir_path string, ctx RenderCTX) error {
	folder_items, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println("[Error]: Could not read the directory: ", dir_path)
		return err
	}
	for _, item := range folder_items {
		this_path := dir_path + "/" + item.Name()
		if item.IsDir() {
			dirClosure(this_path, ctx)
		} else {
			if item.Name() != "+page.hstm" {
				continue
			}
			fileClosure(this_path, ctx)
		}
	}
	return nil
}

func RenderAllSites(sitesDir string, ctx RenderCTX) error {
	files, err := os.ReadDir(sitesDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		dirClosure(sitesDir+"/"+file.Name()+"/pages", ctx)
	}

	return nil
}
