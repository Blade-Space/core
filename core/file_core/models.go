package filecore

type FrontEnd struct {
	Include bool   `yaml:"include"`
	Repo    string `yaml:"repo"`
}

type ReposConfig struct {
	Name     string   `yaml:"name"`
	Versin   string   `yaml:"versin"`
	Port     string   `yaml:"port"`
	Repos    []string `yaml:"repos"`
	Packages []string `yaml:"packages"`
	FrontEnd FrontEnd `yaml:"front-end"`
}
