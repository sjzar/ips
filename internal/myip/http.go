/*
 * Copyright (c) 2023 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package myip

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/pkg/errors"
)

// HTTPDetector is a struct for configuring HTTP queries.
// It contains the host to query for IP discovery.
type HTTPDetector struct {
	Host string
}

// Discover performs an HTTP query based on the HTTPDetector configuration and returns the discovered IP.
func (d *HTTPDetector) Discover(ctx context.Context, localAddr string) (net.IP, error) {
	transport := d.getTransport(localAddr)

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", d.Host, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return d.parseResponseBody(resp.Body)
}

// getTransport creates and returns a transport based on the local address configuration.
func (d *HTTPDetector) getTransport(localAddr string) http.RoundTripper {
	if len(localAddr) == 0 {
		return http.DefaultTransport
	}

	addr, err := net.ResolveTCPAddr("tcp", localAddr+":0")
	if err != nil {
		log.Warnf("Failed to resolve TCP address: %s", err)
		return http.DefaultTransport
	}

	dialer := &net.Dialer{LocalAddr: addr, Timeout: time.Second * 10}
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}
}

// parseResponseBody reads the HTTP response body and extracts the IP address.
func (d *HTTPDetector) parseResponseBody(body io.Reader) (net.IP, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(strings.TrimSpace(string(data)))
	if ip == nil {
		log.Debugf("Invalid IP in response: %s", string(data))
		return nil, errors.ErrDiscoveryFailed
	}
	return ip, nil
}
