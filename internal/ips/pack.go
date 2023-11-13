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
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/format/mmdb"
	"github.com/sjzar/ips/format/plain"
	"github.com/sjzar/ips/internal/ipio"
	"github.com/sjzar/ips/pkg/errors"
)

// Pack reads data from a database file, processes it, and writes it to an output.
func (m *Manager) Pack(_format, file []string, _outputFormat, outputFile string) error {

	if len(_format) == 0 {
		_format = make([]string, len(file))
	} else if len(file) != len(_format) {
		return errors.ErrInvalidFormat
	}

	reader, err := m.createReader(_format, file, true)
	if err != nil {
		log.Debug("m.createReader error: ", err)
		return err
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
	case *mmdb.Writer:
		writerOptionArg, err := url.ParseQuery(m.Conf.WriterOption)
		if err != nil {
			log.Debug("url.ParseQuery error: ", err)
			return err
		}
		option := mmdb.WriterOption{
			SelectLanguages: writerOptionArg.Get("select_languages"),
		}
		if err := writer.SetOption(option); err != nil {
			log.Debug("writer.SetOption error: ", err)
			return err
		}
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
