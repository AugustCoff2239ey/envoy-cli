// Package envset provides rollback support for EnvSet state management.
//
// The RollbackStack allows you to save and restore snapshots of an EnvSet's
// variable state, enabling safe experimentation and change recovery.
//
// # Usage
//
//	stack := envset.NewRollbackStack(10)
//
//	// Save current state before making changes
//	stack.Push(es, "before bulk update")
//
//	// Make changes...
//	es.Vars["KEY"] = "new_value"
//
//	// Restore previous state if needed
//	stack.Pop(es)
//
// # Stack Behaviour
//
// The stack is capped at a configurable maximum depth. When the cap is
// reached, the oldest entry is evicted (FIFO eviction, LIFO restore).
//
// Snapshots are independent copies — mutations to the EnvSet after a Push
// do not affect the saved state.
//
// # CLI Commands
//
//	envoy-cli rollback push <name> <env> -m "message"
//	envoy-cli rollback pop  <name> <env>
//	envoy-cli rollback peek <name> <env>
package envset
