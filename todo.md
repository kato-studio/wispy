# Wispy TODO List
- [ ] **:** 

## Wispy Style 
- [ ] **:** [PROTOTYPE/REWORK] Rework the wispy style package to be based on the minimal pollen design system rather than larger Tailwind system
- [ ] **:** [PROTOTYPE] Prototype theme style variables based pollen base css vars
- [ ] **:** [PROTOTYPE] Theme style variables based pollen base css vars
- [ ] **:** [PROTOTYPE] Auto responsive grid utils
- [ ] **:** 
- [ ] **:** 

## Wispy Template  
- [X] **:** [FEATURE] Condition evaluation (e.g., `{{ if .Title }}` or `{{ if .Title > 10 }}`) 
-   -   new function `ResolveCondition`
- [X] **:** [FEATURE] Condition evaluation (e.g., `{{ each .item in ["key", "value"] }}` or `{{ each .item in .Items }}`)
- [X] **:** [BUG/FIX] Variable which are int's, bools, or floats not being resolved as strings when templating
-   -   new function `stringify`
- [X] **:** [BUG/FIX] Items not being properly assigned to `ctx.Props` and then fail to resolve
- [X] **:** [FEATURE] Create new util for finding closing tags that's able to handle nested tags of the same type
-   -   new function `SeedClosingHandleNested`
- [ ] **:** 
- [ ] **:** 
- [ ] **:** 

## Wispy Engine  
- [ ] **:** Support page layouts
- [ ] **:** Add methods for changing `ctx` instead of allowing direct variable setting
- [ ] **:** Support Toml or Json files for Page content data 
- [ ] **:** Support Component JS
- [ ] **:** Support Component CSS
- [ ] **:** Support Toml or Json files for Page content data 
- [ ] **:** 
- [ ] **:** 
- [ ] **:** 
