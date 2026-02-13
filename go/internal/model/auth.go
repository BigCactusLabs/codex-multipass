package model

// Auth represents the structure of the authentication file (auth.json).
// We use map[string]any to be flexible and preserve all fields,
// as we don't strictly know the schema of the auth token and want to allow forward compatibility.
type Auth map[string]any
