// Copyright 2016 Google Inc. All Rights Reserved.
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

// Package app does all of the work necessary to create a Heapster
// APIServer by binding together the Master Metrics API.
// It can be configured and called directly or via the hyperkube framework.
package app

import (
	"fmt"

	"github.com/Stackdriver/heapster/metrics/options"
	metricsink "github.com/Stackdriver/heapster/metrics/sinks/metric"
	"k8s.io/apimachinery/pkg/util/wait"
	genericapiserver "k8s.io/apiserver/pkg/server"
	v1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

type HeapsterAPIServer struct {
	*genericapiserver.GenericAPIServer
	options    *options.HeapsterRunOptions
	metricSink *metricsink.MetricSink
	nodeLister v1listers.NodeLister
}

// Run runs the specified APIServer. This should never exit.
func (h *HeapsterAPIServer) RunServer() error {
	h.PrepareRun().Run(wait.NeverStop)
	return nil
}

func NewHeapsterApiServer(s *options.HeapsterRunOptions, metricSink *metricsink.MetricSink,
	nodeLister v1listers.NodeLister, podLister v1listers.PodLister) (*HeapsterAPIServer, error) {

	server, err := newAPIServer(s)
	if err != nil {
		return &HeapsterAPIServer{}, err
	}

	installMetricsAPIs(s, server, metricSink, nodeLister, podLister)

	return &HeapsterAPIServer{
		GenericAPIServer: server,
		options:          s,
		metricSink:       metricSink,
		nodeLister:       nodeLister,
	}, nil
}

func newAPIServer(s *options.HeapsterRunOptions) (*genericapiserver.GenericAPIServer, error) {
	if err := s.SecureServing.MaybeDefaultWithSelfSignedCerts("heapster.kube-system", nil, nil); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewConfig(legacyscheme.Codecs)

	if err := s.SecureServing.ApplyTo(&serverConfig.SecureServing, &serverConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	if !s.DisableAuthForTesting {
		if err := s.Authentication.ApplyTo(&serverConfig.Authentication, serverConfig.SecureServing, nil); err != nil {
			return nil, err
		}
		if err := s.Authorization.ApplyTo(&serverConfig.Authorization); err != nil {
			return nil, err
		}
	}

	return serverConfig.Complete(nil).New("heapster", genericapiserver.NewEmptyDelegate())
}
