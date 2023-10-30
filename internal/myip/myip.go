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
	"sync"
	"time"
)

// GetPublicIP concurrently discovers the public IP using multiple detectors.
// It waits until at least 'count' detectors return the same IP, and then returns this IP.
func GetPublicIP(localAddr string, count int, timeoutS int) (net.IP, error) {
	if timeoutS <= 0 {
		timeoutS = 10
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutS)*time.Second)
	defer cancel()

	ipChan := make(chan net.IP, len(Detectors))
	var wg sync.WaitGroup

	for _, detector := range Detectors {
		wg.Add(1)
		go func(det Detector) {
			defer wg.Done()
			ip, err := det.Discover(ctx, localAddr)
			if err == nil {
				ipChan <- ip
			}
		}(detector)
	}

	// Close ipChan once all detectors have finished
	go func() {
		wg.Wait()
		close(ipChan)
	}()

	ipCounts := make(map[string]int)
	for ip := range ipChan {
		ipStr := ip.String()
		ipCounts[ipStr]++
		if ipCounts[ipStr] >= count {
			return ip, nil
		}
	}

	return nil, context.DeadlineExceeded
}
