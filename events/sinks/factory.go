// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sinks

import (
	"fmt"

	"github.com/Stackdriver/heapster/common/flags"
	"github.com/Stackdriver/heapster/events/core"
	"github.com/Stackdriver/heapster/events/sinks/elasticsearch"
	"github.com/Stackdriver/heapster/events/sinks/gcl"
	"github.com/Stackdriver/heapster/events/sinks/honeycomb"
	"github.com/Stackdriver/heapster/events/sinks/influxdb"
	"github.com/Stackdriver/heapster/events/sinks/kafka"
	"github.com/Stackdriver/heapster/events/sinks/log"
	"github.com/Stackdriver/heapster/events/sinks/riemann"

	"github.com/golang/glog"
)

type SinkFactory struct {
}

func (this *SinkFactory) Build(uri flags.Uri) (core.EventSink, error) {
	switch uri.Key {
	case "gcl":
		return gcl.CreateGCLSink(&uri.Val)
	case "log":
		return logsink.CreateLogSink()
	case "influxdb":
		return influxdb.CreateInfluxdbSink(&uri.Val)
	case "elasticsearch":
		return elasticsearch.NewElasticSearchSink(&uri.Val)
	case "kafka":
		return kafka.NewKafkaSink(&uri.Val)
	case "riemann":
		return riemann.CreateRiemannSink(&uri.Val)
	case "honeycomb":
		return honeycomb.NewHoneycombSink(&uri.Val)
	default:
		return nil, fmt.Errorf("Sink not recognized: %s", uri.Key)
	}
}

func (this *SinkFactory) BuildAll(uris flags.Uris) []core.EventSink {
	result := make([]core.EventSink, 0, len(uris))
	for _, uri := range uris {
		sink, err := this.Build(uri)
		if err != nil {
			glog.Errorf("Failed to create %v sink: %v", uri, err)
			continue
		}
		result = append(result, sink)
	}
	return result
}

func NewSinkFactory() *SinkFactory {
	return &SinkFactory{}
}
