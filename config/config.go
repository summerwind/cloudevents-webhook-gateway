package config

type Config struct {
	Listen       string              `json:"listen"`
	TLS          *TLSConfig          `json:"tls"`
	GitHub       *GitHubConfig       `json:"github"`
	Alertmanager *AlertmanagerConfig `json:"alertmanager"`
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

type AlertmanagerConfig struct {
	Path    string `json:"path"`
	Backend string `json:"backend"`
}

func New() *Config {
	return &Config{
		Listen: "0.0.0.0:24381",
		TLS:    &TLSConfig{},
		GitHub: &GitHubConfig{
			Path: "/github",
		},
		Alertmanager: &AlertmanagerConfig{
			Path: "/alertmanager",
		},
	}
}
