package config

type Honeypot struct {
	Protocol string `mapstructure:"protocol" yaml:"protocol"`
	Port     string `mapstructure:"port" yaml:"port"`
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	Fragile  bool   `mapstructure:"fragile" yaml:"fragile"`
}

type HoneypotConfig struct {
	Name      string     `mapstructure:"name" yaml:"name"`
	ShoutUrls []string   `mapstructure:"shout_urls" yaml:"shout_urls"`
	Honeypots []Honeypot `mapstructure:"honeypots" yaml:"honeypots"`
}
