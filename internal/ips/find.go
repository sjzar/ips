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

	"github.com/sjzar/ips/domainlist"
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

// ParseText takes a text input, parses it into segments, and returns the serialized result
// based on the Manager configuration. It returns the combined result as a string.
func (m *Manager) ParseText(text string) (string, error) {

	tp := parser.NewTextParser(text).Parse()

	infoList := make([]interface{}, 0, len(tp.Segments))
	for _, segment := range tp.Segments {
		info, err := m.parseSegment(segment)
		if err != nil {
			log.Debug("m.parseSegment error: ", err)
			return "", err
		}
		infoList = append(infoList, info)
	}

	result, err := m.serialize(infoList)
	if err != nil {
		log.Debug("m.serialize error: ", err)
		return "", err
	}

	return result, nil
}

// parseSegment processes the provided segment and returns the corresponding data.
// This could be IP information, domain information, or raw text.
func (m *Manager) parseSegment(segment parser.Segment) (interface{}, error) {
	switch segment.Type {
	case parser.TextTypeIPv4, parser.TextTypeIPv6:
		return m.parseIP(segment.Content)
	case parser.TextTypeDomain:
		return m.parseDomain(segment.Content)
	case parser.TextTypeText:
		return segment.Content, nil
	}
	return nil, nil
}

// parseIP determines the type of IP (IPv4 or IPv6) and fetches the corresponding information.
func (m *Manager) parseIP(content string) (*model.IPInfo, error) {
	if ip := net.ParseIP(content); ip != nil {
		if ip.To4() != nil {
			return m.parseIPv4(ip)
		}
		return m.parseIPv6(ip)
	}

	return nil, errors.ErrInvalidIP
}

// parseIPv4 finds and returns the information associated with the provided IPv4 address.
func (m *Manager) parseIPv4(ip net.IP) (*model.IPInfo, error) {

	// lazyLoad initializes the IP readers if they haven't been initialized yet.
	if m.ipv4 == nil {
		var err error
		if m.ipv4, err = m.createReader(m.Conf.IPv4Format, m.Conf.IPv4File, false); err != nil {
			log.Debug("createReader error: ", err)
			return nil, err
		}
	}

	return m.ipv4.Find(ip)
}

// parseIPv6 finds and returns the information associated with the provided IPv6 address.
func (m *Manager) parseIPv6(ip net.IP) (*model.IPInfo, error) {

	// lazyLoad initializes the IP readers if they haven't been initialized yet.
	if m.ipv6 == nil {
		var err error
		if m.ipv6, err = m.createReader(m.Conf.IPv6Format, m.Conf.IPv6File, false); err != nil {
			log.Debug("createReader error: ", err)
			return nil, err
		}
	}

	return m.ipv6.Find(ip)
}

// parseDomain fetches the information for the given domain. Implementation is pending.
func (m *Manager) parseDomain(content string) (*model.DomainInfo, error) {
	if ret, ok := domainlist.GetDomainInfo(content); ok {
		return ret, nil
	}

	return &model.DomainInfo{
		Domain: content,
	}, nil
}

// serialize takes a segment and its associated data, then serializes the data
// based on the Manager configuration and returns the serialized string.
func (m *Manager) serialize(data []interface{}) (string, error) {
	switch m.Conf.OutputType {
	case OutputTypeJSON:
		list := &model.DataList{}
		for _, info := range data {
			switch v := info.(type) {
			case *model.IPInfo:
				list.AddItem(v.Output(m.Conf.UseDBFields))
			case *model.DomainInfo:
				list.AddDomain(v)
			case string:
				continue
			}
		}
		return m.serializeDataToJSON(list)
	case OutputTypeAlfred:
		list := &model.DataList{}
		for _, info := range data {
			switch v := info.(type) {
			case *model.IPInfo:
				list.AddAlfredItemByIPInfo(v)
			case *model.DomainInfo:
				list.AddAlfredItemByDomainInfo(v)
			case string:
				continue
			}
		}
		list.AddAlfredItemEmpty()
		return m.serializeDataToJSON(list)
	default:
		// default is OutputTypeText
		buf := &bytes.Buffer{}
		for _, info := range data {
			switch v := info.(type) {
			case *model.IPInfo:
				ret, err := m.serializeIPInfoToText(v)
				if err != nil {
					return "", err
				}
				buf.WriteString(ret)
			case *model.DomainInfo:
				ret, err := m.serializeDomainInfoToText(v)
				if err != nil {
					return "", err
				}
				buf.WriteString(ret)
			case string:
				buf.WriteString(v)
			}
		}
		return buf.String(), nil
	}
}

// serializeIPInfoToText takes an IPInfo, then serializes
// the IPInfo to a text format based on the Manager configuration.
func (m *Manager) serializeIPInfoToText(ipInfo *model.IPInfo) (string, error) {
	values := strings.Join(util.DeleteEmptyValue(ipInfo.Values()), m.Conf.TextValuesSep)
	if values != "" {
		ret := strings.Replace(m.Conf.TextFormat, "%origin", ipInfo.IP.String(), 1)
		ret = strings.Replace(ret, "%values", values, 1)
		return ret, nil
	}

	return ipInfo.IP.String(), nil
}

// serializeDomainInfoToText takes a DomainInfo, then serializes
// the DomainInfo to a text format based on the Manager configuration.
func (m *Manager) serializeDomainInfoToText(domainInfo *model.DomainInfo) (string, error) {
	values := strings.Join(util.DeleteEmptyValue(domainInfo.Values()), m.Conf.TextValuesSep)
	if values != "" {
		ret := strings.Replace(m.Conf.TextFormat, "%origin", domainInfo.Domain, 1)
		ret = strings.Replace(ret, "%values", values, 1)
		return ret, nil
	}

	return domainInfo.Domain, nil
}

// serializeDataToJSON serializes the provided DataList to a JSON format
// based on the Manager configuration. It returns the JSON string.
func (m *Manager) serializeDataToJSON(values *model.DataList) (string, error) {
	if len(values.Items) == 0 && len(values.Domains) == 0 {
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

// createReader decides whether to create a standard or hybrid reader based on the number of files provided.
// It creates a standard reader for a single file and a hybrid reader for multiple files.
func (m *Manager) createReader(_format, file []string, isPackMode bool) (format.Reader, error) {
	if len(file) == 1 {
		return m.createStandardReader(_format[0], file[0], isPackMode)
	}
	return m.createHybridReader(_format, file, isPackMode)
}

// createHybridReader constructs a hybrid reader using multiple IP database formats and files.
// It handles reader creation for each database file and aggregates them into a single hybrid reader.
func (m *Manager) createHybridReader(_format, file []string, isPackMode bool) (format.Reader, error) {
	readers := make([]format.Reader, 0, len(file))
	for i := range file {
		reader, err := m.createStandardReader(_format[i], file[i], isPackMode)
		if err != nil {
			log.Debug("createStandardReader error: ", err)
			return nil, err
		}
		readers = append(readers, reader)
	}

	reader, err := ipio.NewHybridReader(nil, readers...)
	if err != nil {
		log.Debug("ipio.NewHybridReader error: ", err)
		return nil, err
	}

	if len(m.Conf.HybridMode) != 0 {
		if err := reader.SetOption(ipio.HybridReaderOption{
			Mode: m.Conf.HybridMode,
		}); err != nil {
			log.Debug("reader.SetOption error: ", err)
			return nil, err
		}
	}

	if m.Conf.HybridMode != ipio.HybridComparisonMode {
		fs, err := m.newFieldSelector(reader.Meta(), isPackMode)
		if err != nil {
			return nil, err
		}
		reader.OperateChain.Use(fs.Do)
	}

	rw, err := m.newDataRewriter(isPackMode)
	if err != nil {
		return nil, err
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

// createStandardReader sets up and returns an IP reader based on the specified format and file.
// It includes additional processing for field selection and data rewriting based on the configuration.
func (m *Manager) createStandardReader(_format, file string, isPackMode bool) (format.Reader, error) {
	dbr, err := m.createDatabaseReader(_format, file)
	if err != nil {
		return nil, err
	}

	reader := ipio.NewStandardReader(dbr, nil)

	fs, err := m.newFieldSelector(reader.Meta(), isPackMode)
	if err != nil {
		return nil, err
	}
	reader.OperateChain.Use(fs.Do)

	rw, err := m.newDataRewriter(isPackMode)
	if err != nil {
		return nil, err
	}

	// special database process
	if !isPackMode {
		switch dbr.Meta().Format {
		case qqwry.DBFormat:
			rw.LoadString(data.QQwryCountry, data.QQwryArea)
		}
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

// createDatabaseReader initializes a database reader for the given format and file.
// It checks for file existence and downloads the database file if necessary.
func (m *Manager) createDatabaseReader(_format, file string) (format.Reader, error) {
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

	return dbr, nil
}

// newFieldSelector initializes a FieldSelector based on the provided metadata and the pack mode configuration.
// It selects different sets of fields based on whether the pack mode is enabled or not.
func (m *Manager) newFieldSelector(meta *model.Meta, isPackMode bool) (*operate.FieldSelector, error) {
	fields := m.Conf.Fields
	if isPackMode {
		fields = m.Conf.DPFields
	}
	fs, err := operate.NewFieldSelector(meta, fields)
	if err != nil {
		log.Debug("operate.NewFieldSelector error: ", err)
		return nil, err
	}
	return fs, nil
}

// newDataRewriter creates a DataRewriter based on the pack mode configuration.
// It loads different rewrite rules based on whether the pack mode is enabled or not.
func (m *Manager) newDataRewriter(isPackMode bool) (*operate.DataRewriter, error) {
	rw := operate.NewDataRewriter()

	rewriteFiles := m.Conf.RewriteFiles
	if isPackMode {
		rewriteFiles = m.Conf.DPRewriterFiles
	}
	if len(rewriteFiles) > 0 {
		if err := rw.LoadFiles(strings.Split(rewriteFiles, ",")); err != nil {
			log.Debug("rw.LoadFiles error: ", err)
			return nil, err
		}
	}

	// common process
	rw.LoadString(data.ASN2ISP, data.Province, data.City, data.ISP)
	return rw, nil
}
