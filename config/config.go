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
		Events  *bool
		Jvm     *bool
		Process *bool
	}
}
