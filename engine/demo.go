package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func TestHtml(res http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n\nRunning TestHtml")
	var Full_time = time.Now()
	var output []byte = make([]byte, 0)
	var file []byte = make([]byte, 0)
	var err error
	var r Render

	var attFuncMap = make(map[string]AttributeFunc)
	attFuncMap["class"] = func(name, value string) (bool, string, []error) {
		fmt.Println(value)
		return false, "", []error{}
	}

	var opFuncMap = make(map[string]OperationFunc)
	var pageCtx = NewCtx(TemplateCtx{
		Page: TemplatePageCtx{
			Title:   "ABC Example",
			Head:    "<link rel=\"stylesheet\" href=\"/css/style.css\">",
			Meta:    "<meta name=\"description\" content=\"ABC Example\">",
			Styles:  "",
			Scripts: "",
		},
		Data: map[string]any{},
	})
	r = InitEngine(pageCtx, attFuncMap, opFuncMap)
	//
	fmt.Println("Init Duration: ", time.Since(Full_time))

	var OsRead_Duration = time.Now()
	file, err = os.ReadFile("./test.html")
	fmt.Println("os-Read Duration: ", time.Since(OsRead_Duration))
	// Failed to html file
	if err != nil {
		// res.Write([]byte(err.Error()))
		res.Write([]byte("Failed to read file"))
		res.Write([]byte(err.Error()))
		return
	}

	Html_time := time.Now()
	output, errs := r.Render(file, map[string]any{
		"foo":    "bar",
		"whitty": "whimsy",
	})
	//
	fmt.Println("Render Duration: ", time.Since(Html_time))
	//
	if len(errs) > 0 {
		for _, err := range errs {
			res.Write([]byte(err.Error()))
		}
		return
	}
	//
	fmt.Println("Scan Duration: ", time.Since(Full_time))
	//
	// Headers
	res.Header().Add("content-type", "text/html")
	// res.Header().Add("content-type", "text/html")
	// res.Header().Add("content-type", "application/json")
	res.Write(output)
}

func TestTreeGen(res http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n\nRunning TestTreeGen()")
	var Full_time = time.Now()
	var output []byte = make([]byte, 0)
	var file []byte = make([]byte, 0)
	var err error
	var r Render

	//
	fmt.Println("Init Duration: ", time.Since(Full_time))
	var OsRead_Duration = time.Now()
	file, err = os.ReadFile("./test.html")
	fmt.Println("os-Read Duration: ", time.Since(OsRead_Duration))
	// Failed to html file
	if err != nil {
		fmt.Println("[Error] ", err.Error())
		res.Write([]byte("Failed to read file"))
		res.Write([]byte(err.Error()))
		return
	}

	Html_time := time.Now()
	imports := make(map[string][]byte)
	nodeTree, errs := r.BuildNodeTree(file, imports)

	//
	fmt.Println("Tree Gen Duration: ", time.Since(Html_time))
	//
	if errs != nil {
		fmt.Printf(err.Error())
		res.Write([]byte(err.Error()))
	}

	fmt.Println("Scan Duration: (no marshal)", time.Since(Full_time))
	output, err = json.Marshal(nodeTree)
	fmt.Println("Scan Duration: (after marshal)", time.Since(Full_time))
	if err != nil {
		fmt.Printf(err.Error())
		res.Write([]byte(err.Error()))
	}

	//
	// Headers
	res.Header().Add("content-type", "application/json")
	res.Write(output)
}
