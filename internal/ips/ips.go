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

package ips

import (
	"github.com/sjzar/ips/format"
)

// Manager is a command-line tool for IP operations.
type Manager struct {
	// Conf holds the common configurations.
	Conf *Config

	// IPv4 and IPv6 are the IP readers for their respective IP versions.
	ipv4 format.Reader
	ipv6 format.Reader
}

// NewManager initializes and returns a new Manager instance.
func NewManager(conf *Config) *Manager {
	return &Manager{
		Conf: conf,
	}
}
