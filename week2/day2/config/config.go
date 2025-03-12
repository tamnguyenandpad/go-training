package config

type Config struct {
	// Worker pool configuration
	PoolSize              int  // the number of worker in the pool - default 2
	MaxJobs               int  // the maximum number of jobs in the queue - default 2
	WorkerPoolNonBlocking bool // non-blocking mode - default false

	// Greeting service configuration
	BannedNames map[string]struct{} // the list of banned names - default empty

	// HTTP server configuration
	HTTPPort int // the port number for the HTTP server - default 8080
}
