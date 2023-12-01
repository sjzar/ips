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
	"context"
	"net"
	"runtime"
	"slices"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/format/plain"
	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// ChannelBufferSize why 1000? I don't know :)
const ChannelBufferSize = 1000

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
	return NewStandardDumper(r, w).Dump(0)
}

// Dump transfers IP data from the Reader to the Writer.
func (d *StandardDumper) Dump(readerJobs int) error {
	if readerJobs <= 0 {
		switch d.WriterFormat() {
		case plain.DBFormat:
			readerJobs = 1
		default:
			readerJobs = runtime.NumCPU()
		}
	}

	ipStart, ipEnd := net.IPv4(0, 0, 0, 0), ipnet.LastIPv4
	if d.Meta().IsIPv6Support() {
		ipStart, ipEnd = make(net.IP, net.IPv6len), ipnet.LastIPv6
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	retChan := make(chan *model.IPInfo, ChannelBufferSize)
	errChan := make(chan error, readerJobs)
	defer close(errChan)
	wg := sync.WaitGroup{}

	split := ipnet.SplitIPNet(ipStart, ipEnd, readerJobs)
	for i := 0; i < len(split)-1; i++ {
		wg.Add(1)
		go func(ctx context.Context, start, end net.IP) {
			defer wg.Done()
			sd := SimpleDumper{
				Reader:  d.Reader,
				ipStart: start,
				ipEnd:   end,
			}
			if err := sd.Dump(ctx, retChan); err != nil {
				errChan <- err
				return
			}
		}(ctx, split[i], split[i+1])
	}

	go func() {
		wg.Wait()
		close(retChan)
	}()

	for {
		select {
		case ipInfo, ok := <-retChan:
			if !ok {
				return nil
			}
			if err := d.Insert(ipInfo); err != nil {
				log.Debug("StandardDumper Insert() failed ", ipInfo, err)
				return err
			}
		case err := <-errChan:
			if err != nil {
				log.Debug("StandardDumper Dump2() failed ", err)
				return err
			}
		}
	}
}

// SimpleDumper is a structure that facilitates the extraction of IP information within a specified range.
type SimpleDumper struct {
	format.Reader
	ipStart net.IP
	ipEnd   net.IP
}

// Dump iterates over IP addresses in the specified range, sending IP information to retChan.
// The operation is context-aware and will stop if the context is cancelled.
func (d *SimpleDumper) Dump(ctx context.Context, retChan chan<- *model.IPInfo) error {
	// Validate the IP range before proceeding.
	if d.ipStart == nil || d.ipEnd == nil || ipnet.IPLess(d.ipEnd, d.ipStart) {
		return errors.ErrInvalidIPRange
	}

	marker := d.ipStart
	info, err := d.Find(marker)
	if err != nil {
		return err
	}

	// Adjust marker if it doesn't match the start of the IP range.
	if !marker.Equal(info.IPNet.Start) {
		marker = ipnet.NextIP(info.IPNet.End.To16())
	}

	// Iterate over IP addresses until the end of the range is reached.
	for done := d.done(marker, false); !done; done = d.done(marker, true) {
		select {
		case <-ctx.Done(): // Check if the operation was cancelled.
			return nil
		default:
		}
		info, err = d.next(marker)
		if err != nil {
			return err
		}
		if info == nil {
			break
		}
		retChan <- info // Send IP information to the channel.
		marker = ipnet.NextIP(info.IPNet.End.To16())
	}

	return nil
}

// done checks whether the end of the IP range has been reached.
func (d *SimpleDumper) done(marker net.IP, started bool) bool {
	if marker == nil {
		return false
	}

	// Check if the current marker is outside the IP range.
	if !ipnet.Contains(d.ipStart, d.ipEnd, marker) {
		return true
	}

	// Check if the end of the range is reached.
	return started && d.ipEnd.Equal(ipnet.PrevIP(marker))
}

// next retrieves the next IP information from the Reader based on the marker.
func (d *SimpleDumper) next(marker net.IP) (*model.IPInfo, error) {
	if marker == nil {
		marker = d.ipStart
	}

	var currentInfo *model.IPInfo
	for {
		info, err := d.Find(marker)
		if err != nil {
			return nil, err
		}

		// Determine if a new IP range is encountered.
		if currentInfo == nil {
			currentInfo = info
		} else {
			if !slices.Equal(currentInfo.Values(), info.Values()) {
				break
			}
			if ok := currentInfo.IPNet.Join(info.IPNet); !ok {
				break
			}
		}

		marker = ipnet.NextIP(info.IPNet.End)
		if d.done(marker, true) {
			break
		}
	}
	return currentInfo, nil
}
