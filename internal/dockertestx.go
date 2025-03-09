package internal

// GetEnvValue searches the given slice of environment variable strings for the specified key
// and returns its value. If the key is not found, it returns an empty string.
func GetEnvValue(env []string, key string) string {
	prefix := key + "="
	for _, v := range env {
		if len(v) >= len(prefix) && v[:len(prefix)] == prefix {
			return v[len(prefix):]
		}
	}
	return ""
}
