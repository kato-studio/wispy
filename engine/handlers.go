package engine

import (
	"kato-studio/katoengine/utils"
	"strings"
)

// TemplateFunctions processes the template based on the tag type
func TemplateFunctions(content string) string {
    utils.Print("tag: "+content)
    switch {
    case strings.HasPrefix(content, "if"):
        return handleIf(content)
    case strings.HasPrefix(content, "each"):
        return handleEach(content)
    case strings.HasPrefix(content, "for"):
        return handleFor(content)
    default:
        return content
    }
}

// handleIf processes the if tag
func handleIf(rest string) string {
    // Implement if logic here
    // Example: evaluate the condition and return the appropriate content
    return "<!-- If logic -->"
}

// handleEach processes the each tag
func handleEach(rest string) string {
    // Implement each logic here
    // Example: iterate over the collection and return the appropriate content
    return "<!-- Each logic -->"
}

// handleFor processes the for tag
func handleFor(rest string) string {
    // Implement for logic here
    // Example: iterate over the collection and return the appropriate content
    return "<!-- For logic -->"
}


// // INSTRUCTIONS
// // 1. Create a function called TemplateFunctions that takes in a string and 
// // returns a string depending on the template tag type a different handle should be called to achived output
// // write all handles to make this a fully fucntioning template engine that can handle includes, components, if, each and for loops
// // and return the output of the template

// // where these handles will be called from 
// ```func TemplateFunctions(tag string, rest string) string{
// 	var firstRun = tag[0] 
// 	// TODO: add all the handles here

// }```

// exmaple static template  code ```
// <!-- /Parent.kato -->
// <div @root class="flex">
//   <h2>parent here</h2>
//   <small> child here</small>
//   <div>
//     {{#slot}}
//   </div>
// </div>


// <!-- /Footer.kato -->
// <div @root>
//   <h2>Footer</h2>
//   <p>{{company.name}}</p>
//   <div>
//     {{#for link in .links}}
//       <a href="{{.link.url}}">{{.link.text}}</a>
//     {{/for}}
//   </div>
// <div>

// <!-- /Header.kato -->
// <div @root class="flex">
//   <h2>Header</h2>
//   <div>
//     {{#for link in .links}}
//       <a href="{{link.url}}">{{link.text}}</a>
//     {{/for}}
//   </div>
// </div>


// <!-- /+Page.kato -->
// <script>
//   console.log('hello from page')
// </script>

// {{#include 'comps/Header.kato' (links:{{.links}}) }}
//   <Header foo="bar"/>
// {{/include}}

// <h1>Home page</h1>
// {{#include 'comps/Parent.kato' (links:{{.links}}, company:{{.company}}) }}
//   <Parent>
//     <ul>
//       <li>Item 1</li>
//       <li>Item 2</li>
//       <li>Item 3</li>
//     </ul>
//   </Parent>
// {{/include}}

// {{#include 'comps/Footer.kato' (links:{{.links}}, company:{{.company}}) }}
//   <Footer foo="bar"/>
// {{/include}}
// ```


// example output RESULT SHOULD BE THIS! ```
// <script>
//   console.log('hello from page')
// </script>

// <div foo="bar" class="flex">
//   <h2>Header</h2>
//   <div>
// 		<a href="https://google.com">Google</a>
// 		<a href="https://facebook.com">Facebook</a>
// 		<a href="https://twitter.com">Twitter</a>
//   </div>
// </div>

// <h1>Home page</h1>

// <div @root class="flex">
//   <h2>parent here</h2>
//   <small> child here</small>
//   <div>
//     <ul>
//       <li>Item 1</li>
//       <li>Item 2</li>
//       <li>Item 3</li>
//     </ul>
//   </div>
// </div>

// <div foo="bar">
//   <h2>Footer</h2>
//   <p>xyz inc</p>
//   <div>
// 		<a href="https://google.com">Google</a>
// 		<a href="https://facebook.com">Facebook</a>
// 		<a href="https://twitter.com">Twitter</a>
//   </div>
// <div>
// ```