package envset

// Expire provides key-level TTL (time-to-live) management for EnvSet variables.
//
// Use SetExpiry to attach an expiration timestamp to any existing key.
// Expiry metadata is stored in the EnvSet's Meta map under a reserved prefix.
//
// Example:
//
//	err := SetExpiry(e, "TEMP_TOKEN", time.Now().Add(24*time.Hour))
//
// Use IsExpired to check whether a specific key has passed its deadline:
//
//	expired, err := IsExpired(e, "TEMP_TOKEN")
//
// Use PurgeExpired to sweep and remove all expired keys at once:
//
//	purged, err := PurgeExpired(e)
//
// Expiry is stored in RFC3339 UTC format and is compatible with Store
// serialization, so expiry information persists across save/load cycles.
