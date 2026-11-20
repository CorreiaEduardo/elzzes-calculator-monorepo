> For specs 003 in particular, I was worrying about boundaries, do I did a more elaborated prompt:

Implement spec @../docs/specs/003-calculator-ui/SPEC.md, but do not unnecessarily replace, restructure, or reconfigure existing infrastructure. Work with the current setup. And respect @../docs/ARCHITECTURE.md @../docs/PRINCIPLES.md @../docs/DESIGN.md @../docs/DECISIONS.md

The distinction between elzzes-ui-lib and elzzes-calculator-web is a fundamental architectural requirement:

## elzzes-ui-lib (@elzzes/ui)
This package owns reusable UI building blocks, including:
- UI components
- Component variants and sizes
- Shared visual styling and design tokens
- Reusable UI interaction patterns
- Components that can be used independently of the calculator application

Components should be reusable, configurable through props, and independent of calculator-specific business logic.

## elzzes-calculator-web

The web application owns:

- Calculator-specific composition and layouts
- Application state and business workflows
- Calculator-specific interaction logic
- API communication
- Routing, if required
- Application-specific content

The web application should consume reusable components from @elzzes/ui.

# Strict boundary rules

- Do not implement reusable UI components inside the web application when they belong in the UI library.
- Do not duplicate components or styling that already exists in @elzzes/ui.
- If a UI component is missing from the library but is needed by the application, evaluate whether it is a reusable UI primitive. If so, implement it in elzzes-ui-lib, export it through the library's public API, and consume it from the web application.
- Do not import internal UI library files directly. Use its public exports.
- Do not move application logic, API logic, calculator-specific state, or application composition into the UI library.
- The UI library must not depend on the web application.
- Avoid introducing circular dependencies.
- The UI library must remain independently reusable and must not know that it is being used by a calculator application.

## When deciding where something belongs, use this rule:
If it is a reusable UI building block, it belongs in @elzzes/ui. If it describes or implements the calculator application, it belongs in the web application.
Do not weaken or bypass this boundary to implement features faster.
