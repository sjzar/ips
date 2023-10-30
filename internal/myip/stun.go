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
	"strconv"

	"github.com/pion/stun/v2"
	log "github.com/sirupsen/logrus"
)

// STUNDetector is used to discover IP address via STUN protocol.
// It contains the host of the STUN server.
type STUNDetector struct {
	Host string
}

// Discover performs a STUN query to the configured host and returns the discovered IP.
func (d *STUNDetector) Discover(ctx context.Context, localAddr string) (net.IP, error) {
	// Parse a STUN URI
	u, err := stun.ParseURI(d.Host)
	if err != nil {
		return nil, err
	}

	var dialer *net.Dialer
	if localAddr != "" {
		dialer = &net.Dialer{
			LocalAddr: &net.UDPAddr{IP: net.ParseIP(localAddr)},
		}
	} else {
		dialer = &net.Dialer{}
	}

	conn, err := dialer.DialContext(ctx, u.Proto.String(), net.JoinHostPort(u.Host, strconv.Itoa(u.Port)))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	client, err := stun.NewClient(conn)
	if err != nil {
		return nil, err
	}

	// Building binding request with random transaction id.
	message, err := stun.Build(stun.TransactionID, stun.BindingRequest)
	if err != nil {
		return nil, err
	}

	var ip net.IP
	var callbackErr error

	// Callback function to handle the STUN server response
	f := func(res stun.Event) {
		if res.Error != nil {
			callbackErr = res.Error
			return
		}
		// Decode the XOR-MAPPED-ADDRESS attribute from the message
		var xorAddr stun.XORMappedAddress
		callbackErr = xorAddr.GetFrom(res.Message)
		if callbackErr != nil {
			log.Debugf("STUN error: %s", callbackErr)
			return
		}
		ip = xorAddr.IP
	}

	// Sending request to STUN server, waiting for response message.
	if err := client.Do(message, f); err != nil {
		log.Debugf("STUN error: %s", err)
		return nil, err
	}

	if callbackErr != nil {
		return nil, callbackErr
	}

	return ip, nil
}
