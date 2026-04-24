package envset

// Package-level documentation for the signature feature.
//
// Signature provides HMAC-SHA256-based integrity verification for EnvSet
// instances. A signature is computed over all key-value pairs (sorted
// deterministically by key) using a caller-supplied passphrase and stored
// inside the EnvSet's Meta map under a reserved key.
//
// Typical usage:
//
//	// Sign before persisting:
//	_ = envset.Sign(es, passphrase)
//	store.Save(es)
//
//	// Verify after loading:
//	es, _ = store.Load(name, env)
//	if err := envset.VerifySignature(es, passphrase); err != nil {
//	    log.Fatalf("integrity check failed: %v", err)
//	}
//
// Errors:
//   - ErrSignatureMismatch  – stored signature does not match recomputed value.
//   - ErrSignatureNotFound  – no signature has been stored on the EnvSet.
//   - ErrEmptyPassphrase    – passphrase argument was an empty string.
