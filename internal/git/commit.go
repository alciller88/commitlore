package git

import "time"

// Commit represents a single git commit with its metadata.
type Commit struct {
	Hash    string
	Author  string
	Email   string
	Date    time.Time
	Message string
}

// LogOptions configures how commits are retrieved from a repository.
type LogOptions struct {
	Author string
	Since  time.Time
	Until  time.Time
	Limit  int
}
