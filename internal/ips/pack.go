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
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/format/ipdb"
	"github.com/sjzar/ips/format/plain"
	"github.com/sjzar/ips/internal/data"
	"github.com/sjzar/ips/internal/ipio"
	"github.com/sjzar/ips/internal/operate"
)

// Pack reads data from a database file, processes it, and writes it to an output.
func (m *Manager) Pack(_format, file, _outputFormat, outputFile string) error {

	// Create a new DB Reader
	dbr, err := format.NewReader(_format, file)
	if err != nil {
		log.Debug("format.NewReader error: ", err)
		return err
	}
	defer func() {
		_ = dbr.Close()
	}()

	// Initialize the reader
	reader := ipio.NewStandardReader(dbr, nil)

	// Setup field selector
	fs, err := operate.NewFieldSelector(reader.Meta(), m.Conf.DPFields)
	if err != nil {
		log.Debug("operate.NewFieldSelector error: ", err)
		return err
	}
	reader.OperateChain.Use(fs.Do)

	// Setup data rewriter
	rw := operate.NewDataRewriter()
	if len(m.Conf.DPRewriterFiles) > 0 {
		if err := rw.LoadFiles(strings.Split(m.Conf.DPRewriterFiles, ",")); err != nil {
			log.Debug("rw.LoadFiles error: ", err)
			return err
		}
	}

	// common database process
	rw.LoadString(data.ASN2ISP, data.Province, data.City, data.ISP)

	reader.OperateChain.Use(rw.Do)

	// Add specific logic based on the db reader type
	switch dbr.(type) {
	case *ipdb.Reader:
	}

	// Setup the writer
	writer, err := format.NewWriter(_outputFormat, outputFile, reader.Meta())
	if err != nil {
		log.Debug("format.NewWriter error: ", err)
		return err
	}

	// Setup output destination
	output := os.Stdout
	if len(outputFile) != 0 {
		var err error
		output, err = os.Create(outputFile)
		if err != nil {
			log.Debug("os.Create error: ", err)
			return err
		}
		defer func() {
			_ = output.Close()
		}()
	}

	// Add specific logic based on the writer type
	switch writer.(type) {
	case *plain.Writer:
		if err := writer.SetOption(plain.WriterOption{IW: output}); err != nil {
			log.Debug("writer.SetOption error: ", err)
			return err
		}
	}

	// Dump data using the dumper
	dumper := ipio.NewStandardDumper(reader, writer)
	if err := dumper.Dump(); err != nil {
		log.Debug("dumper.Dump error: ", err)
		return err
	}

	// Write to the output destination
	if _, err := dumper.WriteTo(output); err != nil {
		log.Debug("dumper.WriteTo error: ", err)
		return err
	}

	return nil
}
