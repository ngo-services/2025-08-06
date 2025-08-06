package config

type Config struct {
	Port            string
	AllowedTypes    map[string]struct{}
	MaxFilesPerTask int
	MaxActiveTasks  int
}

func New() *Config {
	return &Config{
		Port: ":8080",
		AllowedTypes: map[string]struct{}{
			".pdf":  {},
			".jpeg": {},
		},
		MaxFilesPerTask: 3,
		MaxActiveTasks:  3,
	}
}
