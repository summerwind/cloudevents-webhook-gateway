package config

type Config struct {
	Listen        string               `json:"listen"`
	TLS           *TLSConfig           `json:"tls"`
	GitHub        *GitHubConfig        `json:"github"`
	DockerHub     *DockerHubConfig     `json:"dockerhub"`
	Alertmanager  *AlertmanagerConfig  `json:"alertmanager"`
	AnchoreEngine *AnchoreEngineConfig `json:"anchore-engine"`
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

type DockerHubConfig struct {
	Path    string `json:"path"`
	Backend string `json:"backend"`
}

type AlertmanagerConfig struct {
	Path    string `json:"path"`
	Backend string `json:"backend"`
}

type AnchoreEngineConfig struct {
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
		DockerHub: &DockerHubConfig{
			Path: "/dockerhub",
		},
		Alertmanager: &AlertmanagerConfig{
			Path: "/alertmanager",
		},
		AnchoreEngine: &AnchoreEngineConfig{
			Path: "/anchore-engine",
		},
	}
}
