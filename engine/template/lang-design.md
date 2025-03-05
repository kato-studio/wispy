# Wispy-Template

## Core Goals
- Plug and play modular
- No enteral dependencies (excluding other wispy modular packages)
- Support for component style "blocks" that can tag props and children

## Core Requirements
- Must be able to handle contains blocks I.E. "{% <TAG> ... %} ... {% <TAG> %}"
- Modular support for new tags
- Must be able to liquid style filters "{% .<value> | uppercase %}"
- Fast the engine should be able to parse a page in small page  750-1000 characters in less than 1ms
- Single pass parser & render function that does not require a AST 
- Render context for storing data such as page copy, site info, misc data. as well store for render time generate data such as slots or template defined variables

## Syntax Options Considered
- HTML style tags although the syntax aligning with the core use-case for the language was cool and felt it also bound it's usage to the html syntax and that goes against a core of modularity and flexibility as well as increasing parsing complexity
| <For for="x" as="y"> ... </For>
| <If> ... <If/>




markdown
# Wispy-Template

## Overview
Wispy-Template is a lightweight, modular template engine for Golang, inspired by Go’s templates and Liquid. It is designed with a focus on high performance, flexibility, and ease of extension. With Wispy-Template, you can create dynamic, component-based templates with minimal external dependencies.

## Core Goals
- **Plug and Play Modular:** Easily extend the engine with new tags and filters without modifying the core.
- **Minimal External Dependencies:** Apart from other Wispy modular packages, there is no need for additional dependencies.
- **Component Style Blocks:** Support for blocks that encapsulate props and children, making component-based design simple and efficient.

## Core Requirements
- **Block Handling:** 
  - Must support blocks with start and end delimiters (e.g., `{% TAG ... %} ... {% TAG %}`).
  - *Clarification:* Should we allow self-closing tags or strictly require both opening and closing blocks?
- **Modular Tag Support:** 
  - Developers can add custom tags easily.
- **Liquid Style Filters:** 
  - Enable usage of filters in the template with a pipe syntax (e.g., `{% .VALUE | uppercase %}`).
- **Performance:** 
  - The engine should be capable of parsing a standard template of 1500-2000 characters in under `500ųs`, with consistent performance scaling for larger pages
- **Single Pass Parsing & Rendering:** 
  - The system does not rely on an intermediate AST for parsing and rendering.
- ***note:*** *(although an AST structure may be useful in the future but should work independently from the core parser & render function)*
- **Render Context:** 
  - Provides a mechanism for storing static data (e.g., page copy, site info) and dynamic render-time data (e.g., slots, template-defined variables).

## Syntax Options Considered
### HTML Style Tags
- **Example:**
  ```html
    <div class="pt-22">
      <p>title: {% .title %}</p>
      Hello, {% .name | upcase %}!
    </div>
    {% partial "Boop" %}
  ```
- **Pros:** 
  - Familiar to developers with a web background.
- **Cons:** 
  - Binds the syntax to HTML, reducing flexibility.
  - Increases parsing complexity, decreasing performance scaling.

### Custom Delimiter Blocks
- **Example:**
  ```liquid
  {% <TAG> ... %} ... {% <TAG> %}
  ```
- **Pros:** 
  - More modular and less coupled with HTML.
  - Easier to extend with new custom tags.
- **Cons:** 
  - May require a learning curve for those used to HTML-like syntax.

## Examples

### 1. Applying a Filter
Demonstrates how filters can be applied to transform data.
```liquid
{{ "hello world" | uppercase }}
```
*Expected Output:* `HELLO WORLD`

## Extensibility and Future Features
- **Custom Tags & Filters:** 
  - The architecture should support easy registration and management of custom tags and filters.
- **Potential AST Integration:** 
  - Although the engine currently uses a single-pass parser, an AST may be introduced later for features like template introspection or optimization.
- **Nested Blocks:**
  - Define clear rules on how nested blocks behave and how scopes are managed.
  - *Clarification:* What is the desired behavior when blocks are nested? Should the inner blocks have independent contexts?

## FAQ and Clarifications
1. **Why a Single Pass Parser?**  
   A single pass parser improves performance by reducing overhead, but it may limit some advanced features compared to an AST-based system.
   
2. **How are Filters Processed?**  
   Filters are applied in a pipeline fashion, similar to Liquid. For example, chaining filters like `{{ value | trim | lowercase }}` would first trim the value and then convert it to lowercase.
   
3. **What Happens When a Tag is Not Recognized?**  
   It would be useful to define clear error-handling strategies. Should unrecognized tags fail silently, throw an error, or provide a fallback mechanism?
   
4. **How Can Developers Extend the Engine?**  
   Providing a registration mechanism for new tags and filters is essential. *Clarification:* Would you prefer a central registry or a plugin system where modules are loaded dynamically?

## Conclusion
Wispy-Template aims to combine the best aspects of Golang's template capabilities and Liquid’s flexibility into a fast, modular, and highly extensible engine. With its focus on performance and simplicity, it is well-suited for modern web applications and component-based development.

*Please review the questions and clarifications above. Your feedback will help refine the final documentation and ensure that Wispy-Template meets all your project requirements.*
