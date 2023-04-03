package apiassembler

type APIInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	API     string `yaml:"api"`
}

type APIYaml struct {
	Name    string   `yaml:"name"`
	API     string   `yaml:"api"`
	Link    string   `yaml:"link"`
	Version string   `yaml:"version"`
	Date    string   `yaml:"date"`
	Type    string   `yaml:"type"`
	Authors []string `yaml:"authors"`
}