package runner

type Step struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}

type Pipeline struct {
	Steps []Step `yaml:"steps"`
}
type DeployConfig struct {
	RepoURL    string            `json:"repoUrl"`
	Env        map[string]string `json:"env,omitempty"`
	InstallCmd string            `json:"installCmd,omitempty"`
	BuildCmd   string            `json:"buildCmd,omitempty"`
	StartCmd   string            `json:"startCmd,omitempty"`
	ExposePort int               `json:"port,omitempty"`
	AppType    string            `json:"appType"`
}
