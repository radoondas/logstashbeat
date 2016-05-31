package beater

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/radoondas/logstashbeat/config"
)

const selector = "logstashbeat"
const selectorDetail = "json"

type Logstashbeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
	client     publisher.Client

	urls []*url.URL

	Node struct {
		events  bool
		jvm     bool
		process bool
	}
}

// Creates beater
func New() *Logstashbeat {
	return &Logstashbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Logstashbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := b.RawConfig.Unpack(&bt.beatConfig)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Logstashbeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Logstashbeat.Period == "" {
		bt.beatConfig.Logstashbeat.Period = "1s"
	}

	bt.client = b.Publisher.Connect()

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Logstashbeat.Period)
	if err != nil {
		return err
	}

	//define default URL if none provided
	var urlConfig []string
	if bt.beatConfig.Logstashbeat.URLs != nil {
		urlConfig = bt.beatConfig.Logstashbeat.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1:9600"}
	}

	bt.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid Logstash url: %v", err)
			return err
		}
		bt.urls[i] = u
	}

	if bt.beatConfig.Logstashbeat.Node.Events != nil {
		bt.Node.events = *bt.beatConfig.Logstashbeat.Node.Events
	} else {
		bt.Node.events = true
	}

	if bt.beatConfig.Logstashbeat.Node.Jvm != nil {
		bt.Node.jvm = *bt.beatConfig.Logstashbeat.Node.Jvm
	} else {
		bt.Node.jvm = true
	}

	if bt.beatConfig.Logstashbeat.Node.Process != nil {
		bt.Node.process = *bt.beatConfig.Logstashbeat.Node.Process
	} else {
		bt.Node.process = true
	}

	if !bt.Node.events && !bt.Node.jvm && !bt.Node.process {
		return errors.New("Invalid configuration. Nothing to request! Check your configuration file.")
	}

	logp.Debug(selector, "Init configuration logstashbeat")
	logp.Debug(selector, "Period %v\n", bt.period)
	logp.Debug(selector, "URLs %v", bt.urls)
	logp.Debug(selector, "Node Events statistics %t\n", bt.Node.events)
	logp.Debug(selector, "Node JVM statistics %t\n", bt.Node.jvm)
	logp.Debug(selector, "Node Process statistics %t\n", bt.Node.process)

	return nil
}

func (bt *Logstashbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run elasticbeat")

	//for each url
	for _, u := range bt.urls {
		go func(u *url.URL) {
			ticker := time.NewTicker(bt.period)
			defer ticker.Stop()

			for {
				select {
				case <-bt.done:
					goto GotoFinish
				case <-ticker.C:
				}

				timerStart := time.Now()

				if bt.Node.jvm {
					logp.Debug(selector, "Node/stats/jvm for url: %v", u)
					jvm, err := bt.GetNodeStatsJVM(*u)
					if err != nil {
						logp.Err("Error reading Node/stats/jvm metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/stats/jvm metrics detail: %+v", jvm)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "nodeStats",
							"url":        u.String(),
							"jvm":        jvm,
						}
						logp.Debug(selectorDetail, "Published Event detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				if bt.Node.events {
					logp.Debug(selector, "Node/stats/events for url: %v", u)
					ev, err := bt.GetNodeStatsEvents(*u)
					if err != nil {
						logp.Err("Error reading Node/stats/events metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/stats/events metrics detail: %+v", ev)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "nodeStats",
							"url":        u.String(),
							"events":     ev,
						}
						logp.Debug(selectorDetail, "Published Event detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				//if bt.Node.process {}

				timerEnd := time.Now()
				duration := timerEnd.Sub(timerStart)
				if duration.Nanoseconds() > bt.period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}(u)
	}

	<-bt.done
	return nil
}

func (bt *Logstashbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Logstashbeat) Stop() {
	close(bt.done)
}
