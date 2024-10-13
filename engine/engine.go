package engine

import (
	"fmt"
	"io"
	"kato-studio/katoengine/utils"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/valyala/fasttemplate"
)

func SlipEngine(template string, json gjson.Result) string {
	var result string
	var included_components map[string]string = make(map[string]string)
	template = utils.CleanString(template)

	// Handle filling in variables
	init_data := fasttemplate.New(template, "{{.", "}}")
	result = init_data.ExecuteFuncString(func(w io.Writer, content string) (int, error) {
		if(json.Get(content).Exists()) {
			return w.Write([]byte(json.Get(content).String()))
		}
		return w.Write([]byte(""))
	})

	// Core templating functionality
	templ_func := fasttemplate.New(result, "{{#", "}}")
	result = templ_func.ExecuteFuncString(func(w io.Writer, content string) (int, error) {
		if(strings.HasPrefix(content, "include")){
			// Handle includes/imports
			// skip the first 8 characters and split the rest at the first space this is to ignore the include tag
			var comp_path, comp_name = utils.SplitAtRune(content[8:], ' ')
			if(comp_name == ""){
				utils.Error("failed to import component: "+content)
			}
			
			// find component at path and set it's content to 
			valid_path := "./view/"+strings.Replace(comp_path, "c/", "components/", 1)
			if(!strings.HasSuffix(valid_path, ".kato")){
				valid_path = valid_path+".kato"
			}

			file_bytes, err := os.ReadFile(valid_path)
			if err != nil {
				utils.Error("failed find component: "+valid_path)
			}
			
			included_components[comp_name] = string(file_bytes)
			return w.Write([]byte(""))

		}else if(strings.HasPrefix(content, "render")){
			// Skip rendering if not render tag as it needs to be rendered last
			return w.Write([]byte("{{#render "+content+"}}"))
		}else{
			// Handle templating functions
			utils.Print("TemplateFunctions: "+content)
			handle_content := TemplateFunctions(content)
			return w.Write([]byte(handle_content))
		}
	})


	// Handle rendering components
	render_func := fasttemplate.New(result, "{{#render", "}}")
	result = render_func.ExecuteFuncString(func(w io.Writer, content string) (int, error) {
		content = strings.Replace(content, "render <", "", 1)
		content = strings.Trim(content, " ")
		var comp_name, rest = utils.SplitAtRune(content, ' ')
		if(rest == ""){
			utils.Error("failed to render component: "+rest)
			return w.Write([]byte(comp_name))
		}

		
		hasChildren := strings.HasSuffix(strings.Trim(rest, " "), "</"+comp_name+">")

		utils.Print("-----------")
		utils.Debug("name: "+comp_name + "   children? "+fmt.Sprint(hasChildren))
		var remaining string
		var children string
		var props string
		var attrs string

		attrs, remaining = utils.SplitAt(rest, "props(")
		if(remaining != ""){
			// if passing props using json() syntax then split at the closing bracket
			props, remaining = utils.SplitAt(remaining, ")")
			
		}else{
			// if no props are passed then split at the closing bracket
			if(hasChildren){
				attrs, remaining = utils.SplitAt(rest, ">")
			}else{
				attrs, remaining = utils.SplitAt(rest, "/>")
			}
		}


		if(hasChildren){
			children, _ = utils.SplitAt(remaining, "</"+comp_name+">")
			if(children[0] == '>'){
				children = children[1:]
			}
		}else{
			// do shit
		}

		utils.Debug("attrs: "+attrs)
		utils.Debug("remaining: "+remaining)
		utils.Debug("props: "+props)
		utils.Debug("children: "+children)
		utils.Print("-----------")

		// var children string
		// var _closing string
		// if(hasChildren){
		// 	children, _ = utils.SplitAt(rest, "<"+comp_name+"/>")
		// }else{
		// 	children, _ = utils.SplitAt(rest, "/>")
		// }

		return w.Write([]byte("{{COMPONENT "+comp_name+"}}"))
	})

	return result
}

func RenderComponent(component string, attributes string, children string, json gjson.Result) string {
	// 
	
	return ""
}