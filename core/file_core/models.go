package filecore

type ReposConfig struct {
	Name     string   `yaml:"name"`
	Versin   string   `yaml:"versin"`
	Port     string   `yaml:"port"`
	Repos    []string `yaml:"repos"`
	Packages []string `yaml:"packages"`
}
