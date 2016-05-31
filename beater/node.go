package beater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const nodeStatsEventsURI = "/_node/stats/events?pretty"
const nodeStatsJVMURI = "/_node/stats/jvm?pretty"

//{
//  "events" : {
//    "in" : 0,
//    "filtered" : 0,
//    "out" : 0
//  }
//}
type NodeStatsEvents struct {
	Events struct {
		In     uint64 `json:"in"`
		Filter uint64 `json:"filter"`
		Out    uint64 `json:"out"`
	} `json:"events"`
}

//{
//  "timestamp" : 1464210557679,
//  "uptime_in_millis" : 1129085,
//  "mem" : {
//  "mem" : {
//    "heap_used_in_bytes" : 155156864,
//    "heap_used_percent" : 7,
//    "heap_committed_in_bytes" : 247332864,
//    "heap_max_in_bytes" : 2077753344,
//    "non_heap_used_in_bytes" : 151577520,
//    "non_heap_committed_in_bytes" : 159006720,
//      "pools" : {
//        "survivor" : {
//          "peak_used_in_bytes" : 4259840,
//          "used_in_bytes" : 4535200,
//          "peak_max_in_bytes" : 34865152,
//          "max_in_bytes" : 69730304
//        },
//        "old" : {
//          "peak_used_in_bytes" : 57657080,
//          "used_in_bytes" : 104508352,
//          "peak_max_in_bytes" : 724828160,
//          "max_in_bytes" : 1449656320
//        },
//        "young" : {
//          "peak_used_in_bytes" : 34078720,
//          "used_in_bytes" : 46871568,
//          "peak_max_in_bytes" : 279183360,
//          "max_in_bytes" : 558366720
//      }
//    }
//  }
//}
type NodeStatsJVM struct {
	Timestamp        uint64 `json:"timestamp"`
	Uptime_in_millis uint64 `json:"uptime_in_millis"`
	Mem              struct {
		Heap_used_in_bytes          uint64 `json:"heap_used_in_bytes"`
		Heap_used_percent           uint64 `json:"heap_used_percent"`
		Heap_committed_in_bytes     uint64 `json:"heap_committed_in_bytes"`
		Heap_max_in_bytes           uint64 `json:"heap_max_in_bytes"`
		Non_heap_used_in_bytes      uint64 `json:"non_heap_used_in_bytes"`
		Non_heap_committed_in_bytes uint64 `json:"non_heap_committed_in_bytes"`
		Pools                       struct {
			Young struct {
				Used_in_bytes      uint64 `json:"used_in_bytes"`
				Max_in_bytes       uint64 `json:"max_in_bytes"`
				Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
				Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
			} `json:"young"`
			Survivor struct {
				Used_in_bytes      uint64 `json:"used_in_bytes"`
				Max_in_bytes       uint64 `json:"max_in_bytes"`
				Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
				Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
			} `json:"survivor"`
			Old struct {
				Used_in_bytes      uint64 `json:"used_in_bytes"`
				Max_in_bytes       uint64 `json:"max_in_bytes"`
				Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
				Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
			} `json:"old"`
		} `json:"pools"`
	} `json:"mem"`
}

func (bt *Logstashbeat) GetNodeStatsEvents(u url.URL) (NodeStatsEvents, error) {
	events := NodeStatsEvents{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + nodeStatsEventsURI)
	if err != nil {
		return events, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return events, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return events, err
	}

	err = json.Unmarshal([]byte(body), &events)
	if err != nil {
		return events, err
	}

	return events, nil
}

func (bt *Logstashbeat) GetNodeStatsJVM(u url.URL) (NodeStatsJVM, error) {
	nodeStatsJVM := NodeStatsJVM{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + nodeStatsJVMURI)
	if err != nil {
		return nodeStatsJVM, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nodeStatsJVM, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nodeStatsJVM, err
	}

	err = json.Unmarshal([]byte(body), &nodeStatsJVM)
	if err != nil {
		return nodeStatsJVM, err
	}

	return nodeStatsJVM, nil
}
