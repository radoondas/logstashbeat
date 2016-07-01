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
		Stats struct {
			events  bool
			jvm     bool
			process bool
			mem     bool
		}
		pipeline bool
		jvm      bool
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

	if bt.beatConfig.Logstashbeat.Node.Stats.Events != nil {
		bt.Node.Stats.events = *bt.beatConfig.Logstashbeat.Node.Stats.Events
	} else {
		bt.Node.Stats.events = true
	}

	if bt.beatConfig.Logstashbeat.Node.Stats.Jvm != nil {
		bt.Node.Stats.jvm = *bt.beatConfig.Logstashbeat.Node.Stats.Jvm
	} else {
		bt.Node.Stats.jvm = true
	}

	if bt.beatConfig.Logstashbeat.Node.Stats.Process != nil {
		bt.Node.Stats.process = *bt.beatConfig.Logstashbeat.Node.Stats.Process
	} else {
		bt.Node.Stats.process = true
	}

	if bt.beatConfig.Logstashbeat.Node.Stats.Mem != nil {
		bt.Node.Stats.mem = *bt.beatConfig.Logstashbeat.Node.Stats.Mem
	} else {
		bt.Node.Stats.mem = true
	}

	if bt.beatConfig.Logstashbeat.Node.Pipeline != nil {
		bt.Node.pipeline = *bt.beatConfig.Logstashbeat.Node.Pipeline
	} else {
		bt.Node.pipeline = true
	}

	if bt.beatConfig.Logstashbeat.Node.Jvm != nil {
		bt.Node.jvm = *bt.beatConfig.Logstashbeat.Node.Jvm
	} else {
		bt.Node.jvm = true
	}

	if !bt.Node.Stats.events && !bt.Node.Stats.jvm && !bt.Node.Stats.process && !bt.Node.Stats.mem && !bt.Node.pipeline && !bt.Node.jvm {
		return errors.New("Invalid configuration. Nothing to request! Check your configuration file.")
	}

	logp.Debug(selector, "Init configuration logstashbeat")
	logp.Debug(selector, "Period %v\n", bt.period)
	logp.Debug(selector, "URLs %v", bt.urls)
	logp.Debug(selector, "NodeStats Events statistics %t\n", bt.Node.Stats.events)
	logp.Debug(selector, "NodeStats JVM statistics %t\n", bt.Node.Stats.jvm)
	logp.Debug(selector, "NodeStats Process statistics %t\n", bt.Node.Stats.process)
	logp.Debug(selector, "NodeStats Mem statistics %t\n", bt.Node.Stats.mem)
	logp.Debug(selector, "Node JVM statistics %t\n", bt.Node.jvm)
	logp.Debug(selector, "Node Pipeline statistics %t\n", bt.Node.pipeline)

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

				if bt.Node.Stats.jvm {
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
							"jvm":        jvm.JVM,
						}
						logp.Debug(selectorDetail, "Published Event detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				if bt.Node.Stats.events {
					logp.Debug(selector, "Node/stats/events for url: %v", u)
					e, err := bt.GetNodeStatsEvents(*u)
					if err != nil {
						logp.Err("Error reading Node/stats/events metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/stats/events metrics detail: %+v", e)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "nodeStats",
							"url":        u.String(),
							"events":     e.Events,
						}
						logp.Debug(selectorDetail, "Published Event detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				if bt.Node.Stats.process {
					logp.Debug(selector, "Node/stats/process for url: %v", u)
					p, err := bt.GetNodeStatsProcess(*u)
					if err != nil {
						logp.Err("Error reading Node/stats/process metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/stats/process metrics detail: %+v", p)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "nodeStats",
							"url":        u.String(),
							"process":    p.Process,
						}
						logp.Debug(selectorDetail, "Published Process detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				if bt.Node.Stats.mem {
					logp.Debug(selector, "Node/stats/mem for url: %v", u)
					m, err := bt.GetNodeStatsMem(*u)
					if err != nil {
						logp.Err("Error reading Node/stats/mem metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/stats/mem metrics detail: %+v", m)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "nodeStats",
							"url":        u.String(),
							"mem":        m.Mem,
						}
						logp.Debug(selectorDetail, "Published Mem detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				if bt.Node.pipeline {
					logp.Debug(selector, "Node/pipeline for url: %v", u)
					p, err := bt.GetNodePipeline(*u)
					if err != nil {
						logp.Err("Error reading Node/pipeline metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/pipeline metrics detail: %+v", p)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "node",
							"url":        u.String(),
							"pipeline":   p.Pipeline,
						}
						logp.Debug(selectorDetail, "Published Pipeline detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

				if bt.Node.jvm {
					logp.Debug(selector, "Node/jvm for url: %v", u)
					j, err := bt.GetNodeJVM(*u)
					if err != nil {
						logp.Err("Error reading Node/jvm metrics: %v", err)
					} else {
						logp.Debug(selectorDetail, "Node/jvm metrics detail: %+v", j)

						event := common.MapStr{
							"@timestamp": common.Time(time.Now()),
							"type":       "node",
							"url":        u.String(),
							"jvm":        j.Jvm,
						}
						logp.Debug(selectorDetail, "Published JVM detail: %+v", event)
						bt.client.PublishEvent(event)
					}
				}

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
