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
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockDetector is a mock implementation of the Detector interface
type MockDetector struct {
	IP  net.IP
	Err error
}

func (d *MockDetector) Discover(ctx context.Context, localAddr string) (net.IP, error) {
	return d.IP, d.Err
}

func TestMyIP(t *testing.T) {
	type Config struct {
		LocalAddr   string
		MyIPCount   int
		MyIPTimeout int
	}

	tests := []struct {
		name          string
		detectors     []Detector
		config        Config
		expectedIP    net.IP
		expectedError error
	}{
		{
			name: "Successfully retrieve IP with majority detectors",
			detectors: []Detector{
				&MockDetector{IP: net.ParseIP("192.168.1.1"), Err: nil},
				&MockDetector{IP: net.ParseIP("192.168.1.1"), Err: nil},
				&MockDetector{IP: net.ParseIP("10.0.0.1"), Err: nil},
			},
			config:        Config{MyIPCount: 2, MyIPTimeout: 5},
			expectedIP:    net.ParseIP("192.168.1.1"),
			expectedError: nil,
		},
		{
			name: "Error if no majority IP found",
			detectors: []Detector{
				&MockDetector{IP: net.ParseIP("192.168.1.1"), Err: nil},
				&MockDetector{IP: net.ParseIP("10.0.0.1"), Err: nil},
				&MockDetector{IP: net.ParseIP("172.16.0.1"), Err: nil},
			},
			config:        Config{MyIPCount: 2, MyIPTimeout: 5},
			expectedIP:    nil,
			expectedError: context.DeadlineExceeded,
		},
		// ... you can add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Detectors = tt.detectors
			ip, err := GetPublicIP(tt.config.LocalAddr, tt.config.MyIPCount, tt.config.MyIPTimeout)
			assert.Equal(t, tt.expectedIP, ip)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
