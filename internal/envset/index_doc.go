package envset

// Index builds a positional and metadata index over an EnvSet.
//
// Each key in the EnvSet is assigned an IndexEntry containing:
//   - Position: zero-based sorted position of the key
//   - Group:    the group name stored in Meta under "group:<key>"
//   - Tags:     comma-separated tags stored in Meta under "tags:<key>"
//
// Use IndexByGroup to retrieve all entries belonging to a specific group,
// ordered by their original position.
//
// Example:
//
//	idx, err := Index(es)
//	if err != nil { ... }
//	appEntries := IndexByGroup(idx, "app")
//	for _, e := range appEntries {
//	    fmt.Printf("[%d] %s (tags: %v)\n", e.Position, e.Key, e.Tags)
//	}
package envset
