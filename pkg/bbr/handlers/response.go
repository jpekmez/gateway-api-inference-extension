/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"fmt"
	"time"

	basepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	eppb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

// HandleResponseHeaders handles response headers.
func (s *Server) HandleResponseHeaders(headers *eppb.HttpHeaders, reqCtx *RequestContext) ([]*eppb.ProcessingResponse, error) {
	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_ResponseHeaders{
				ResponseHeaders: &eppb.HeadersResponse{
					Response: &eppb.CommonResponse{
						HeaderMutation: &eppb.HeaderMutation{
							SetHeaders: s.generateResponseHeaders(reqCtx),
						},
					},
				},
			},
		},
	}, nil
}

// HandleResponseBody handles response bodies.
func (s *Server) HandleResponseBody(body *eppb.HttpBody) ([]*eppb.ProcessingResponse, error) {
	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_ResponseBody{
				ResponseBody: &eppb.BodyResponse{},
			},
		},
	}, nil
}

// HandleResponseTrailers handles response trailers.
func (s *Server) HandleResponseTrailers(trailers *eppb.HttpTrailers) ([]*eppb.ProcessingResponse, error) {
	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_ResponseTrailers{
				ResponseTrailers: &eppb.TrailersResponse{},
			},
		},
	}, nil
}

func (s *Server) generateResponseHeaders(reqCtx *RequestContext) []*basepb.HeaderValueOption {
	headers := []*basepb.HeaderValueOption{}

	// Calculate and add latency if request timestamp is available
	if reqCtx != nil && !reqCtx.RequestReceivedTimestamp.IsZero() {
		latencyMs := time.Since(reqCtx.RequestReceivedTimestamp).Milliseconds()
		headers = append(headers, &basepb.HeaderValueOption{
			Header: &basepb.HeaderValue{
				Key:      "x-bbr-latency-ms",
				RawValue: []byte(fmt.Sprintf("%d", latencyMs)),
			},
		})
	}

	return headers
}
