package envset

// Package-level documentation for the protect feature.
//
// Protect allows individual keys within an EnvSet to be marked as protected.
// Protected keys cannot be overwritten or deleted without first removing
// their protection via UnprotectKey.
//
// Typical usage:
//
//	// Protect a sensitive key
//	_ = envset.ProtectKey(e, "PROD_DB_PASSWORD")
//
//	// Check protection before mutation
//	if envset.IsProtected(e, key) {
//		return errors.New("cannot modify protected key")
//	}
//
//	// List all protected keys
//	keys := envset.ProtectedKeys(e)
//
// Protection state is stored in the EnvSet Meta map and is persisted
// alongside the rest of the EnvSet when saved via the Store.
