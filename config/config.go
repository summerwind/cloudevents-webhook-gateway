package config

type Config struct {
	Listen        string        `json:"listen"`
	TLS           *TLSConfig    `json:"tls"`
	GitHub        *GitHubConfig `json:"github"`
	DockerHub     *ProxyConfig  `json:"dockerhub"`
	Alertmanager  *ProxyConfig  `json:"alertmanager"`
	AnchoreEngine *ProxyConfig  `json:"anchore-engine"`
	Clair         *ProxyConfig  `json:"clair"`
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

type ProxyConfig struct {
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
		DockerHub: &ProxyConfig{
			Path: "/dockerhub",
		},
		Alertmanager: &ProxyConfig{
			Path: "/alertmanager",
		},
		AnchoreEngine: &ProxyConfig{
			Path: "/anchore-engine",
		},
		Clair: &ProxyConfig{
			Path: "/clair",
		},
	}
}
