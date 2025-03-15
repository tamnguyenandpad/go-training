package config

type Config struct {
	NumWorkers    int
	QueueCap      int
	IsNonBlocking bool
	BannedNames   map[string]struct{}
	HTTPPort      int
}
