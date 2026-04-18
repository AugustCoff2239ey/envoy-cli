package envset

// MaskValue replaces all but the last `revealSuffix` characters of a value
// with asterisks. If the value is shorter than or equal to revealSuffix, the
// entire value is masked.
//
// Example:
//
//	MaskValue("supersecret", 4) // => "*******cret"
//	MaskValue("hi", 4)          // => "**"

// IsSensitiveKey returns true if the key name suggests it holds a sensitive
// value (e.g. contains "SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE").

// MaskSensitive automatically masks values whose keys are considered sensitive
// according to IsSensitiveKey. The original EnvSet is not modified; a new map
// of key→masked-value is returned.
//
// Example:
//
//	masked := MaskSensitive(es, 4)
//	fmt.Println(masked["DB_PASSWORD"]) // "****word"

// MaskKeys masks only the explicitly supplied keys, regardless of whether they
// are considered sensitive. Unknown keys are silently ignored.
