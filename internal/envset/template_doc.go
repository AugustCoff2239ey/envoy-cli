// Package envset — template rendering
//
// RenderTemplate substitutes {{KEY}} placeholders in a template string with
// values from an EnvSet. Placeholders follow the pattern:
//
//	{{UPPER_SNAKE_CASE_KEY}}
//
// Example:
//
//	e, _ := envset.New("app", "production")
//	e.Vars["DB_HOST"] = "db.example.com"
//	e.Vars["DB_PORT"] = "5432"
//
//	res, err := envset.RenderTemplate(e, "postgres://{{DB_HOST}}:{{DB_PORT}}/mydb")
//	// res.Rendered  => "postgres://db.example.com:5432/mydb"
//	// res.Unresolved => []
//
// If a placeholder key is not present in the EnvSet it is left unchanged and
// its name is added to TemplateResult.Unresolved so callers can warn the user.
//
// Helper functions:
//
//	ExtractPlaceholders(tmpl)         — list all {{KEY}} names in a template
//	TemplateComplete(e, tmpl)         — true when every placeholder is satisfied
//	MissingPlaceholders(e, tmpl)      — keys referenced but absent from e
//	SuggestPlaceholders(e, unresolved)— fuzzy suggestions for typo detection
package envset
