// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Logstashbeat LogstashbeatConfig
}

type LogstashbeatConfig struct {
	Period string   `config:"period"`
	URLs   []string `yaml:"urls"`
	Node   struct {
		Stats struct {
			Events  *bool `json:"events"`
			Jvm     *bool `json:"jvm"`
			Process *bool `json:"process"`
			Mem     *bool `json:"mem"`
		} `json:"stats"`
		Pipeline *bool `json:"pipeline"`
		Jvm      *bool `json:"jvm"`
	} `json:"node"`
}
