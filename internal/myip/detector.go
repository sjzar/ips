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
)

type Detector interface {
	Discover(ctx context.Context, localAddr string) (ip net.IP, err error)
}

var Detectors = []Detector{
	// STUN servers
	&STUNDetector{Host: "stun:stun.l.google.com:19302"},
	&STUNDetector{Host: "stun:stun1.l.google.com:19302"},
	&STUNDetector{Host: "stun:stun2.l.google.com:19302"},
	&STUNDetector{Host: "stun:stun3.l.google.com:19302"},
	&STUNDetector{Host: "stun:stun4.l.google.com:19302"},
	&STUNDetector{Host: "stun:stun.aa.net.uk:3478"},
	&STUNDetector{Host: "stun:stun.hoiio.com:3478"},
	&STUNDetector{Host: "stun:stun.acrobits.cz:3478"},
	&STUNDetector{Host: "stun:stun.voip.blackberry.com:3478"},
	&STUNDetector{Host: "stun:stun.sip.us:3478"},
	&STUNDetector{Host: "stun:stun.antisip.com:3478"},
	&STUNDetector{Host: "stun:stun.avigora.fr:3478"},
	&STUNDetector{Host: "stun:stun.linphone.org:3478"},
	&STUNDetector{Host: "stun:stun.voipgate.com:3478"},
	&STUNDetector{Host: "stun:stun.cope.es:3478"},
	&STUNDetector{Host: "stun:stun.bluesip.net:3478"},
	&STUNDetector{Host: "stun:stun.solcon.nl:3478"},
	&STUNDetector{Host: "stun:stun.uls.co.za:3478"},

	// HTTP servers
	&HTTPDetector{Host: "http://inet-ip.info/ip"},
	&HTTPDetector{Host: "http://whatismyip.akamai.com/"},
	&HTTPDetector{Host: "https://ipecho.net/plain"},
	&HTTPDetector{Host: "https://eth0.me/"},
	&HTTPDetector{Host: "https://ifconfig.me/ip"},
	&HTTPDetector{Host: "https://checkip.amazonaws.com/"},
	&HTTPDetector{Host: "https://wgetip.com/"},
	&HTTPDetector{Host: "https://ip.tyk.nu/"},
	&HTTPDetector{Host: "https://l2.io/ip"},
	&HTTPDetector{Host: "https://api.ipify.org/"},
	&HTTPDetector{Host: "https://myexternalip.com/raw"},
	&HTTPDetector{Host: "https://icanhazip.com"},
	&HTTPDetector{Host: "https://ifconfig.io/ip"},
	&HTTPDetector{Host: "https://ifconfig.co/ip"},
	&HTTPDetector{Host: "https://ipinfo.io/ip"},
	&HTTPDetector{Host: "https://wtfismyip.com/text"},

	// DNS servers
	&DNSDetector{Domain: "myip.opendns.com.", Server: "resolver1.opendns.com:53", QueryType: "A"},
	&DNSDetector{Domain: "myip.opendns.com.", Server: "resolver2.opendns.com:53", QueryType: "A"},
	&DNSDetector{Domain: "myip.opendns.com.", Server: "resolver3.opendns.com:53", QueryType: "A"},
	&DNSDetector{Domain: "myip.opendns.com.", Server: "resolver4.opendns.com:53", QueryType: "A"},
	&DNSDetector{Domain: "whoami.akamai.net.", Server: "ns1-1.akamaitech.net:53", QueryType: "A"},
	&DNSDetector{Domain: "whoami.ultradns.net.", Server: "pdns1.ultradns.net:53", QueryType: "A"},
	&DNSDetector{Domain: "o-o.myaddr.l.google.com.", Server: "ns1.google.com:53", QueryType: "TXT"},
}
