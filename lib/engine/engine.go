package engine

import (
	"fmt"
	"kato-studio/katoengine/lib/engine/extract"
	"kato-studio/katoengine/lib/engine/template"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
)

func Render(raw_content []byte, data interface{}, components map[string][]byte) string {

	content := utils.CleanString(string(raw_content))

	target := extract.ContentScanner(content)

	utils.Debug("Template Targets")
	for _, target := range target.Components {
		name := extract.ComponentName(content[target:])
		utils.Debug(name)
	}
	utils.Debug("---------------")
	for _, target := range target.TemplateOps {
		name := extract.OperationName(content[target:])
		utils.Debug(name)
	}
	utils.Debug("---------------")

	// operationReplaceStore
	opStore := store.SmallIntStore()
	// check if there are any operations
	opsArrayLen := len(target.TemplateOps) - 1
	if opsArrayLen > 0 {
		opStore.Set(0, content[:target.TemplateOps[0]])
		index := 0
		for _, location := range target.TemplateOps {
			result, endIndex := RenderOperations(content, data, location, opStore)
			opStore.Set(location, result)

			if opsArrayLen > index+1 {
				opStore.Set(endIndex, content[endIndex:endIndex])
			} else if opsArrayLen == index+1 {
				opStore.Set(endIndex, content[endIndex:])
			}

			index++
		}
	}

	result := "<title>Ummm???</title>"
	sortedOpsMap := opStore.Sorted()
	utils.Debug("Sorted Ops Map")
	utils.Print(fmt.Sprint(sortedOpsMap))
	for key, _ := range sortedOpsMap {
		utils.Debug("result: " + result)
		result += sortedOpsMap[key]
	}

	return result
}

func RenderOperations(content string, data interface{}, location int, opStore store.SmallIntStoreInstance) (string, int) {
	// Get the operation name
	opName := extract.OperationName(content[location:])
	// _ is nestedCount, we don't need it here
	// but it represents the number of nested operations of the
	operationEndIndex, nestedCount := extract.FindOperationEndTag(content[location:], opName)

	operationContent := content[location : location+operationEndIndex]
	utils.Debug("Operation " + opName + " Content:")
	utils.Print(operationContent)

	if nestedCount > 0 {
		return "BOOP!!! NESTED!!!", operationEndIndex
	} else {
		result := template.OperationParser(operationContent, opName, data)
		opStore.Set(location, result)
		return result, operationEndIndex
	}

}

func RenderComponent(component string, data interface{}) string {
	//
	return ""
}
