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
	"bytes"
	"encoding/json"
	"net"
	"net/url"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/format/mmdb"
	"github.com/sjzar/ips/format/qqwry"
	"github.com/sjzar/ips/internal/data"
	"github.com/sjzar/ips/internal/ipio"
	"github.com/sjzar/ips/internal/operate"
	"github.com/sjzar/ips/internal/parser"
	"github.com/sjzar/ips/internal/util"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// ParseText parses the provided text and returns the result based on the Manager configuration.
func (m *Manager) ParseText(text string) (string, error) {

	buf := &bytes.Buffer{}
	tp := parser.NewTextParser(text).Parse()

	for _, segment := range tp.Segments {
		var result string
		var err error
		switch m.Conf.OutputType {
		case OutputTypeJSON:
			result, err = m.GetJsonResult(segment.Type, segment.Content)
		default:
			// OutputTypeText
			result = m.GetTextResult(segment.Type, segment.Content)
		}
		if err != nil {
			log.Debug("GetJsonResult error: ", err)
			return "", err
		}
		buf.WriteString(result)
	}

	return buf.String(), nil
}

// GetJsonResult returns the JSON representation for the given type and content.
func (m *Manager) GetJsonResult(_type, content string) (string, error) {
	var values interface{}
	switch _type {
	case parser.TextTypeIPv4:
		ipInfo, err := m.IPv4Find(net.ParseIP(content))
		if err != nil {
			log.Debug("m.IPv4.Find error: ", err)
			return "", err
		}
		values = ipInfo.Output(m.Conf.UseDBFields)
	case parser.TextTypeIPv6:
		ipInfo, err := m.IPv6Find(net.ParseIP(content))
		if err != nil {
			log.Debug("m.IPv6.Find error: ", err)
			return "", err
		}
		values = ipInfo.Output(m.Conf.UseDBFields)
	case parser.TextTypeDomain:
	}

	if values == nil {
		return "", nil
	}

	var ret []byte
	var err error
	if m.Conf.JsonIndent {
		ret, err = json.MarshalIndent(values, "", "  ")
	} else {
		ret, err = json.Marshal(values)
	}
	if err != nil {
		log.Debug("json.Marshal error: ", err)
		return "", err
	}

	return string(ret) + "\n", nil
}

// GetTextResult returns the text representation for the given type and content.
func (m *Manager) GetTextResult(_type, content string) string {
	values := ""
	switch _type {
	case parser.TextTypeIPv4:
		if ipInfo, err := m.IPv4Find(net.ParseIP(content)); err == nil {
			values = strings.Join(util.DeleteEmptyValue(ipInfo.Values()), m.Conf.TextValuesSep)
		}
	case parser.TextTypeIPv6:
		if ipInfo, err := m.IPv6Find(net.ParseIP(content)); err == nil {
			values = strings.Join(util.DeleteEmptyValue(ipInfo.Values()), m.Conf.TextValuesSep)
		}
	case parser.TextTypeDomain:
	}

	if values != "" {
		ret := strings.Replace(m.Conf.TextFormat, "%origin", content, 1)
		ret = strings.Replace(ret, "%values", values, 1)
		return ret
	}

	return content
}

// IPv4Find finds the IP information for the given IPv4 address.
func (m *Manager) IPv4Find(ip net.IP) (*model.IPInfo, error) {

	// lazyLoad initializes the IP readers if they haven't been initialized yet.
	if m.ipv4 == nil {
		var err error
		if m.ipv4, err = m.createReader(m.Conf.IPv4Format, m.Conf.IPv4File); err != nil {
			log.Debug("createReader error: ", err)
			return nil, err
		}
	}

	return m.ipv4.Find(ip)
}

// IPv6Find finds the IP information for the given IPv6 address.
func (m *Manager) IPv6Find(ip net.IP) (*model.IPInfo, error) {

	// lazyLoad initializes the IP readers if they haven't been initialized yet.
	if m.ipv6 == nil {
		var err error
		if m.ipv6, err = m.createReader(m.Conf.IPv6Format, m.Conf.IPv6File); err != nil {
			log.Debug("createReader error: ", err)
			return nil, err
		}
	}

	return m.ipv6.Find(ip)
}

// createReader creates an IP reader based on the provided format and file.
func (m *Manager) createReader(_format, file string) (format.Reader, error) {
	if !util.IsFileExist(file) {
		fullpath := filepath.Join(m.Conf.IPSDir, file)
		if !util.IsFileExist(fullpath) {
			// init database file
			_, ok := DownloadMap[file]
			if !ok {
				log.Debugf("file not found %s", file)
				return nil, errors.ErrFileNotFound
			}
			if err := m.Download(file, ""); err != nil {
				return nil, err
			}
		}
		file = fullpath
	}

	dbr, err := format.NewReader(_format, file)
	if err != nil {
		log.Debug("format.NewReader error: ", _format, file, err)
		return nil, err
	}

	switch dbr.(type) {
	case *mmdb.Reader:
		readerOptionArg, err := url.ParseQuery(m.Conf.ReaderOption)
		if err != nil {
			log.Debug("url.ParseQuery error: ", err)
			return nil, err
		}
		option := mmdb.ReaderOption{
			DisableExtraData: readerOptionArg.Get("disable_extra_data") == "true",
			UseFullField:     readerOptionArg.Get("use_full_field") == "true",
		}
		if err := dbr.SetOption(option); err != nil {
			log.Debug("reader.SetOption error: ", err)
			return nil, err
		}
	}

	reader := ipio.NewStandardReader(dbr, nil)

	fs, err := operate.NewFieldSelector(reader.Meta(), m.Conf.Fields)
	if err != nil {
		log.Debug("operate.NewFieldSelector error: ", err)
		return nil, err
	}
	reader.OperateChain.Use(fs.Do)

	rw := operate.NewDataRewriter()
	if len(m.Conf.RewriteFiles) > 0 {
		if err := rw.LoadFiles(strings.Split(m.Conf.RewriteFiles, ",")); err != nil {
			log.Debug("rw.LoadFiles error: ", err)
			return nil, err
		}
	}

	// common process
	rw.LoadString(data.ASN2ISP, data.Province, data.City, data.ISP)

	// special database process
	switch dbr.Meta().Format {
	case qqwry.DBFormat:
		rw.LoadString(data.QQwryCountry, data.QQwryArea)
	}

	reader.OperateChain.Use(rw.Do)

	if len(m.Conf.Lang) != 0 {
		tl, err := operate.NewTranslator(m.Conf.Lang)
		if err != nil {
			log.Debug("operate.NewTranslator error: ", err)
			return nil, err
		}
		reader.OperateChain.Use(tl.Do)
	}

	return reader, nil
}
