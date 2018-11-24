package config

type Config struct {
	Listen string        `json:"listen"`
	TLS    *TLSConfig    `json:"tls"`
	GitHub *GitHubConfig `json:"github"`
}

type TLSConfig struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

type GitHubConfig struct {
	Path    string `json:"path"`
	Backend string `json:"backend"`
	Secret  string `json:"secret"`
}

func New() *Config {
	return &Config{
		Listen: "0.0.0.0:14381",
		TLS:    &TLSConfig{},
		GitHub: &GitHubConfig{
			Path: "/github",
		},
	}
}
