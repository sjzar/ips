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

//
// import (
// 	"math"
// 	"math/rand"
// 	"net"
// 	"time"
//
// 	"github.com/tonobo/mtr/pkg/icmp"
// )
//
// const (
// 	ICMP_TIMEOUT_MS   = 800
// 	ICMP_INTERVAL_MS  = 100
// 	ICMP_MAXHOPS      = 64
// 	ICMP_HOP_SLEEP_NS = 1
// )
//
// // GetHops returns the number of hops to reach the destination address.
// // Need to run as root.
// func GetHops(localAddr, address string) (int, bool) {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	seq := r.Intn(math.MaxUint16)
// 	id := r.Intn(math.MaxUint16) & 0xffff
//
// 	ipAddr := net.IPAddr{IP: net.ParseIP(address)}
//
// 	for i := 1; i <= 5; i++ {
// 		time.Sleep(ICMP_INTERVAL_MS * time.Millisecond)
//
// 		for ttl := 1; ttl < ICMP_MAXHOPS; ttl++ {
// 			seq++
// 			time.Sleep(ICMP_HOP_SLEEP_NS * time.Nanosecond)
// 			var hopReturn icmp.ICMPReturn
// 			if ipAddr.IP.To4() != nil {
// 				hopReturn, _ = icmp.SendDiscoverICMP(localAddr, &ipAddr, ttl, id, ICMP_TIMEOUT_MS*time.Millisecond, seq)
// 			} else {
// 				hopReturn, _ = icmp.SendDiscoverICMPv6(localAddr, &ipAddr, ttl, id, ICMP_TIMEOUT_MS*time.Millisecond, seq)
// 			}
// 			if hopReturn.Addr == address {
// 				return ttl, true
// 			}
// 		}
// 	}
// 	return 0, false
// }
