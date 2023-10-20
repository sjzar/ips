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

package ipio

import (
	"net"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/model"
)

// StandardDumper serves as a standard mechanism to transfer IP database from one format to another.
type StandardDumper struct {
	format.Reader
	format.Writer

	marker net.IP        // keeps track of the current IP address being processed
	info   *model.IPInfo // holds the current IP information
	done   bool          // flag to indicate if processing is complete
	err    error         // holds any error that occurs during processing
}

// NewStandardDumper initializes and returns a new StandardDumper.
func NewStandardDumper(r format.Reader, w format.Writer) *StandardDumper {
	return &StandardDumper{
		Reader: r,
		Writer: w,
	}
}

// Dump is a convenience method to transfer IP data from a reader to a writer.
// It is equivalent to calling NewStandardDumper(r, w).Dump().
func Dump(r format.Reader, w format.Writer) error {
	return NewStandardDumper(r, w).Dump()
}

// Dump transfers IP data from the Reader to the Writer.
func (d *StandardDumper) Dump() error {
	for d.Next() {
		if err := d.Insert(d.info); err != nil {
			return err
		}
	}
	if d.err != nil {
		return d.err
	}
	return nil
}

// Next fetches the next IP information from the Reader.
func (d *StandardDumper) Next() bool {
	if d.done {
		return false
	}
	if d.marker == nil {
		if d.Meta().IsIPv6Support() {
			d.marker = make(net.IP, net.IPv6len)
		} else {
			d.marker = net.IPv4(0, 0, 0, 0)
		}
	}
	d.info = nil

	// Continuously fetch IP information until either a change in the data is found or the end of the IP range is reached.
	for {
		info, err := d.Find(d.marker)
		if err != nil {
			log.Debug("StandardDumper Find() failed ", d.marker, info, err)
			d.err = err
			return false
		}

		if d.info == nil {
			d.info = info
		} else {
			if !reflect.DeepEqual(d.info.Values(), info.Values()) {
				break
			}
			if ok := d.info.IPNet.Join(info.IPNet); !ok {
				break
			}
		}

		// Move to the next IP address in the range.
		d.marker = ipnet.NextIP(info.IPNet.End)
		if ipnet.IsLastIP(info.IPNet.End, d.Meta().IsIPv6Support()) {
			d.done = true
			break
		}
	}
	return true
}
