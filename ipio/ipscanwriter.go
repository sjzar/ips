/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

// IPScanWriter IPScan 写入工具
type IPScanWriter struct {
	meta model.Meta
	w    io.Writer
	b    *bytes.Buffer
}

// NewIPScanWriter 初始化 IPScan 写入实例
func NewIPScanWriter(meta model.Meta, w io.Writer) (*IPScanWriter, error) {

	writer := &IPScanWriter{
		meta: meta,
		w:    w,
	}

	if err := writer.init(); err != nil {
		return nil, err
	}

	return writer, nil
}

// init 初始化 IPScan 写入实例
func (w *IPScanWriter) init() error {

	if w.w == nil {
		w.b = bytes.NewBuffer([]byte{})
		w.w = w.b
	}

	// write header
	str := fmt.Sprintf("# ScanTime: %s\n", time.Now().Local().Format("2006-01-02 15:04:05"))
	str += fmt.Sprintf("# Fields: %s\n", strings.Join(w.meta.Fields, FieldSep))
	str += fmt.Sprintf("# IPVersion: %d\n", w.meta.IPVersion)
	b, _ := json.Marshal(w.meta)
	str += fmt.Sprintf("%s%s\n", MetaPrefix, string(b))
	if _, err := fmt.Fprint(w.w, str); err != nil {
		return err
	}

	return nil
}

// Insert 插入数据
func (w *IPScanWriter) Insert(ipr *ipx.Range, values []string) error {
	if w.w == nil {
		if err := w.init(); err != nil {
			return err
		}
	}
	for _, ipNet := range ipr.IPNets() {
		if _, err := fmt.Fprintf(w.w, "%s\t%s\n", ipNet.String(), strings.Join(values, FieldSep)); err != nil {
			return err
		}
	}

	return nil
}

// Save 保存数据
func (w *IPScanWriter) Save(_w io.Writer) error {
	if w.b != nil {
		_, err := w.b.WriteTo(_w)
		if err != nil {
			return err
		}
	}
	return nil
}
