# Removed, to be replaced with DaisyUI 5.0! after some testing!
## backedup to branch `backup-POC-style-engine`

### component css
https://github.com/saadeghi/daisyui/tree/master/packages/daisyui

## https://daisyui.com/docs/v5
  - Install
  - Core Improvement
  - Build and integration improvements 
  - Design System Improvements
  - Themes and styling
  - New components


## Build and integration improvements 
TLDR – Import only the parts you need.
Micro CSS files are now available for no-build projects.
Native CSS nesting reduces CSS size.
It's ESM compatible and has dependency-free class name prefixing.

Native CSS nesting

CSS nesting is now supported on all browsers. daisyUI 5 uses CSS nesting which prevents duplication of CSS rules and results smaller CSS size in your browser!
ESM compatibility

daisyUI 5 is now ESM (ECMAScript Module) compatible. Which means you can import and use specific parts of the library in JS if you need to.
Dependency-free class name prefixing

daisyUI 5 can now prefix class names without a dependency.
Micro CSS files for No-Build projects

For server-side rendered projects (Rails, Django, PHP, etc) or projects that don't have a JS build step (HTMX, Alpine.js, WordPress, etc), it's now possible to use specific parts of daisyUI without including the entire library or even without Tailwind CSS.

For example if you only want to use daisyUI toggle component, include a tiny CSS file that only contains the styles for the toggle component:
Before

Not possible

After

https://cdn.jsdelivr.net/npm/daisyui@5/components/toggle.css


daisyui@5 html tags
```html
<link href="https://cdn.jsdelivr.net/npm/daisyui@5" rel="stylesheet" type="text/css" />
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
```


## Ideas / Thoughts
build out custom golang engine for on-demand bundling and building out utility glasses tailwind cdn currently provides this would include supporting dynamic classes "[...]"