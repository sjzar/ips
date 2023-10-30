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
	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/internal/myip"
)

func (m *Manager) MyIP() (string, error) {
	ip, err := myip.GetPublicIP(m.Conf.LocalAddr, m.Conf.MyIPCount, m.Conf.MyIPTimeoutS)
	if err != nil {
		log.Debugf("myip.GetPublicIP error: %v", err)
		return "", err
	}

	return m.ParseText(ip.String())
}
