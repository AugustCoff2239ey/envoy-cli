// Package envset — depend.go
//
// # Key Dependency Tracking
//
// The dependency feature lets you declare that one environment variable
// depends on one or more others being present. This is useful for
// documenting and enforcing configuration contracts.
//
// ## Functions
//
//   - AddDependency(es, key, dep)  — declare that key depends on dep
//   - RemoveDependency(es, key, dep) — remove a declared dependency
//   - GetDependencies(es, key)     — list all deps for a key
//   - CheckDependencies(es)        — return keys referenced as deps but absent
//
// ## Storage
//
// Dependencies are stored in EnvSet.Meta under the key "dep:<KEY>" as a
// comma-separated list of dependency key names.
package envset
