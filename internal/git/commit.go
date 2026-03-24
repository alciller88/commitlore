package git

import "time"

// Commit represents a single git commit with its metadata.
type Commit struct {
	Hash    string    `json:"hash"`
	Author  string    `json:"author"`
	Email   string    `json:"email"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
}

// LogOptions configures how commits are retrieved from a repository.
type LogOptions struct {
	Author string
	Since  time.Time
	Until  time.Time
	Limit  int
}
