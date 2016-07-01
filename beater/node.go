package beater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const nodeStatsEventsURI = "/_node/stats/events"
const nodeStatsJVMURI = "/_node/stats/jvm"
const nodeStatsProcessURI = "/_node/stats/process"
const nodeStatsMemURI = "/_node/stats/mem"

const nodePipelineURI = "/_node/pipeline"
const nodeJVMURI = "/_node/jvm"

//const nodeOsURI = "/_node/os"

//{
//	"host" : "localhost",
//	"version" : "5.0.0-alpha4",
//	"http_address" : "127.0.0.1:9600",
//	"events" : {
//		"in" : 3,
//		"filtered" : 3,
//		"out" : 3
//	}
//}
type NodeStatsEvents struct {
	Events struct {
		In       uint64 `json:"in"`
		Filtered uint64 `json:"filtered"`
		Out      uint64 `json:"out"`
	} `json:"events"`
}

//{
//	"host" : "localhost",
//	"version" : "5.0.0-alpha4",
//	"http_address" : "127.0.0.1:9600",
//	"jvm" : {
//		"threads" : {
//			"count" : 19,
//			"peak_count" : 22
//		}
//	}
//}
type NodeStatsJVM struct {
	JVM struct {
		Threads struct {
			Count      uint64 `json:"count"`
			Peak_count uint64 `json:"peak_count"`
		} `json:"threads"`
	} `json:"jvm"`
}

//{
//	"host" : "localhost",
//	"version" : "5.0.0-alpha4",
//	"http_address" : "127.0.0.1:9600",
//	"process" : {
//		"open_file_descriptors" : 45,
//		"peak_open_file_descriptors" : 49,
//		"max_file_descriptors" : 4096,
//		"mem" : {
//			"total_virtual_in_bytes" : 4709322752
//		},
//		"cpu" : {
//			"total_in_millis" : 57440000000,
//			"percent" : 0
//		}
//	}
//}
type NodeStatsProcess struct {
	Process struct {
		Peak_open_file_descriptors uint64 `json:"peak_open_file_descriptors"`
		Max_file_descriptors       uint64 `json:"max_file_descriptors"`
		Open_file_descriptors      uint64 `json:"open_file_descriptors"`
		Mem                        struct {
			Total_virtual_in_bytes uint64 `json:"total_virtual_in_bytes"`
		} `json:"mem"`
		Cpu struct {
			Total_in_millis uint64 `json:"total_in_millis"`
			Percent         uint64 `json:"percent"`
		} `json:"cpu"`
	} `json:"process"`
}

//{
//	"host" : "localhost",
//	"version" : "5.0.0-alpha4",
//	"http_address" : "127.0.0.1:9600",
//	"mem" : {
//		"heap_used_in_bytes" : 276561168,
//		"heap_used_percent" : 13,
//		"heap_committed_in_bytes" : 519045120,
//		"heap_max_in_bytes" : 2077753344,
//		"non_heap_used_in_bytes" : 164040800,
//		"non_heap_committed_in_bytes" : 173449216,
//		"pools" : {
//			"survivor" : {
//				"peak_used_in_bytes" : 8912896,
//				"used_in_bytes" : 11325176,
//				"peak_max_in_bytes" : 34865152,
//				"max_in_bytes" : 69730304,
//				"committed_in_bytes" : 17825792
//			},
//			"old" : {
//				"peak_used_in_bytes" : 99454200,
//				"used_in_bytes" : 147187720,
//				"peak_max_in_bytes" : 724828160,
//				"max_in_bytes" : 1449656320,
//				"committed_in_bytes" : 357957632
//			},
//			"young" : {
//				"peak_used_in_bytes" : 71630848,
//				"used_in_bytes" : 118048272,
//				"peak_max_in_bytes" : 279183360,
//				"max_in_bytes" : 558366720,
//				"committed_in_bytes" : 143261696
//			}
//		}
//	}
//}
type NodeStatsMem struct {
	Mem struct {
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
				Committed_in_bytes uint64 `json:"committed_in_bytes"`
			} `json:"young"`
			Survivor struct {
				Used_in_bytes      uint64 `json:"used_in_bytes"`
				Max_in_bytes       uint64 `json:"max_in_bytes"`
				Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
				Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
				Committed_in_bytes uint64 `json:"committed_in_bytes"`
			} `json:"survivor"`
			Old struct {
				Used_in_bytes      uint64 `json:"used_in_bytes"`
				Max_in_bytes       uint64 `json:"max_in_bytes"`
				Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
				Peak_max_in_bytes  uint64 `json:"peak_max_in_bytes"`
				Committed_in_bytes uint64 `json:"committed_in_bytes"`
			} `json:"old"`
		} `json:"pools"`
	} `json:"mem"`
}

//{
//	"host" : "localhost",
//	"version" : "5.0.0-alpha4",
//	"http_address" : "127.0.0.1:9600",
//	"pipeline" : {
//		"workers" : 1,
//		"batch_size" : 125,
//		"batch_delay" : 5
//	}
//}
type NodePipeline struct {
	Pipeline struct {
		Workers     uint64 `json:"workers"`
		Batch_size  uint64 `json:"batch_size"`
		Batch_delay uint64 `json:"batch_delay"`
	} `json:"pipeline"`
}

//{
//	"host" : "localhost",
//	"version" : "5.0.0-alpha4",
//	"http_address" : "127.0.0.1:9600",
//	"jvm" : {
//		"pid" : 20351,
//		"version" : "1.8.0_74",
//		"vm_name" : "Java HotSpot(TM) 64-Bit Server VM",
//		"vm_version" : "1.8.0_74",
//		"vm_vendor" : "Oracle Corporation",
//		"start_time_in_millis" : 1467367397367,
//		"mem" : {
//			"heap_init_in_bytes" : 268435456,
//			"heap_max_in_bytes" : 1038876672,
//			"non_heap_init_in_bytes" : 2555904,
//			"non_heap_max_in_bytes" : 0
//		}
//	}
//}
type NodeJVM struct {
	Jvm struct {
		Start_time_in_millis uint64 `json:"start_time_in_millis"`
		Mem                  struct {
			Heap_init_in_bytes     uint64 `json:"heap_init_in_bytes"`
			Heap_max_in_bytes      uint64 `json:"heap_max_in_bytes"`
			Non_heap_init_in_bytes uint64 `json:"non_heap_init_in_bytes"`
			Non_heap_max_in_bytes  uint64 `json:"non_heap_max_in_bytes"`
		} `json:"mem"`
	} `json:"jvm"`
}

//OS is not implemented
//{
//	"os": {
//		"name": "Mac OS X",
//		"arch": "x86_64",
//		"version": "10.11.2",
//		"available_processors": 8
//	}
//}

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

func (bt *Logstashbeat) GetNodeStatsProcess(u url.URL) (NodeStatsProcess, error) {
	nodeStatsProcess := NodeStatsProcess{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + nodeStatsProcessURI)

	if err != nil {
		return nodeStatsProcess, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nodeStatsProcess, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nodeStatsProcess, err
	}

	err = json.Unmarshal([]byte(body), &nodeStatsProcess)
	if err != nil {
		return nodeStatsProcess, err
	}

	return nodeStatsProcess, nil
}

func (bt *Logstashbeat) GetNodeStatsMem(u url.URL) (NodeStatsMem, error) {
	nodeStatsMem := NodeStatsMem{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + nodeStatsMemURI)

	if err != nil {
		return nodeStatsMem, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nodeStatsMem, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nodeStatsMem, err
	}

	err = json.Unmarshal([]byte(body), &nodeStatsMem)
	if err != nil {
		return nodeStatsMem, err
	}

	return nodeStatsMem, nil
}

func (bt *Logstashbeat) GetNodeJVM(u url.URL) (NodeJVM, error) {
	nodeJVM := NodeJVM{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + nodeJVMURI)

	if err != nil {
		return nodeJVM, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nodeJVM, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nodeJVM, err
	}

	err = json.Unmarshal([]byte(body), &nodeJVM)
	if err != nil {
		return nodeJVM, err
	}

	return nodeJVM, nil
}

func (bt *Logstashbeat) GetNodePipeline(u url.URL) (NodePipeline, error) {
	nodePipeline := NodePipeline{}

	res, err := http.Get(TrimSuffix(u.String(), "/") + nodePipelineURI)

	if err != nil {
		return nodePipeline, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nodePipeline, fmt.Errorf("HTTP%s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nodePipeline, err
	}

	err = json.Unmarshal([]byte(body), &nodePipeline)
	if err != nil {
		return nodePipeline, err
	}

	return nodePipeline, nil
}
