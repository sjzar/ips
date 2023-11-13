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
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/internal/operate"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

const (
	// HybridComparisonMode is used in scenarios where the goal is to compare data
	// across multiple IP database readers. In this mode, the output includes data
	// from all readers, facilitating a comprehensive comparison of the differences
	// between each data source. This mode is particularly useful for analyzing
	// discrepancies and inconsistencies across various IP databases.
	HybridComparisonMode = "comparison"

	// HybridAggregationMode is designed for situations where a unified view of data
	// is needed. In this mode, the output is a single, aggregated set of data that
	// combines information from all the IP database readers. If a particular field is
	// missing from one reader, it is supplemented by another, ensuring a more
	// complete and cohesive dataset. This mode is ideal for creating a comprehensive
	// and enriched view of IP information by leveraging the strengths of multiple databases.
	HybridAggregationMode = "aggregation"
)

// HybridReader integrates multiple IP database readers into a single entity.
// It supports various operational modes (like comparison and aggregation) for querying
// and combining results from different IP databases.
type HybridReader struct {
	OperateChain *operate.IPOperateChain // Chain of operations to be applied to IP information.
	dbReaders    []format.Reader         // Collection of database readers.
	meta         *model.Meta             // Combined metadata from all readers.
	hybridMode   string                  // Operational mode of the HybridReader.
}

// NewHybridReader constructs a new HybridReader with the provided IP operation chain and database readers.
// It initializes the reader's metadata by merging metadata from all provided readers.
func NewHybridReader(operateChain *operate.IPOperateChain, dbReaders ...format.Reader) (*HybridReader, error) {
	// Check if there are any readers provided
	if len(dbReaders) == 0 {
		return nil, errors.ErrNoDatabaseReaders
	}

	if operateChain == nil {
		operateChain = operate.NewIPOperateChain()
	}

	// Initialize meta with the first reader's meta information
	hybridMeta := &model.Meta{
		Fields:     make([]string, 0),
		FieldAlias: make(map[string]string),
	}

	// Merge Fields and FieldAlias from all readers
	for i, reader := range dbReaders {
		readerMeta := reader.Meta()
		if i == 0 {
			hybridMeta.MetaVersion = readerMeta.MetaVersion
			hybridMeta.IPVersion = readerMeta.IPVersion
			hybridMeta.Format = readerMeta.Format
		}

		for _, field := range readerMeta.Fields {
			if strings.HasPrefix(field, fmt.Sprintf("%d_", i)) {
				hybridMeta.Fields = append(hybridMeta.Fields, field)
				continue
			}
			if split := strings.SplitN(field, "_", 2); len(split) == 2 {
				if _, err := strconv.ParseInt(split[0], 10, 64); err == nil {
					continue
				}
			}
			prefixedField := fmt.Sprintf("%d_%s", i, field)
			hybridMeta.Fields = append(hybridMeta.Fields, prefixedField)
		}

		for key, value := range readerMeta.FieldAlias {
			prefixedKey := fmt.Sprintf("%d_%s", i, key)
			prefixedVal := fmt.Sprintf("%d_%s", i, value)
			hybridMeta.FieldAlias[prefixedKey] = prefixedVal
		}
	}

	return &HybridReader{
		OperateChain: operateChain,
		dbReaders:    dbReaders,
		meta:         hybridMeta,
	}, nil
}

// Meta returns the combined metadata of the IP databases attached to the HybridReader.
func (h *HybridReader) Meta() *model.Meta {
	return h.meta
}

// Find performs parallel queries across all underlying database readers using the specified IP address.
// It combines the results based on the operational mode of the HybridReader: either comparing data or aggregating it.
func (h *HybridReader) Find(ip net.IP) (*model.IPInfo, error) {
	var wg sync.WaitGroup
	results := make([]*model.IPInfo, len(h.dbReaders))
	errs := make([]error, len(h.dbReaders))

	// Parallel query from each Reader
	for i, reader := range h.dbReaders {
		wg.Add(1)
		go func(index int, rd format.Reader) {
			defer wg.Done()
			result, err := rd.Find(ip)
			results[index] = result
			if err != nil {
				errs[index] = fmt.Errorf("error in Reader %d: %w", index, err)
			}
		}(i, reader)
	}

	wg.Wait()

	// Check for errors and combine results
	hybridIPInfo := &model.IPInfo{
		IP:            ip,
		Data:          make(map[string]string),
		FieldAlias:    h.Meta().FieldAlias,
		Fields:        h.Meta().Fields,
		ReplaceFields: make(map[string]string),
	}

	for _, err := range errs {
		if err != nil {
			// Return the first non-nil error with reader index
			return nil, err
		}
	}

	for i, result := range results {
		if i == 0 {
			hybridIPInfo.IPNet = result.IPNet
		} else if ok := hybridIPInfo.IPNet.CommonRange(ip, result.IPNet); !ok {
			log.Debug("IPNet.CommonRange() failed ", hybridIPInfo.IPNet, ip, result.IPNet)
		}

		for key, value := range result.Data {
			prefixedKey := fmt.Sprintf("%d_%s", i, key)
			hybridIPInfo.Data[prefixedKey] = value
		}

		for key, value := range result.ReplaceFields {
			prefixedKey := fmt.Sprintf("%d_%s", i, key)
			hybridIPInfo.ReplaceFields[prefixedKey] = value
		}

		switch h.hybridMode {
		case HybridComparisonMode:
		default:
			// HybridAggregationMode
			for _, field := range h.Meta().Fields {
				if _, ok := hybridIPInfo.Data[field]; ok {
					continue
				}
				val, ok := result.GetData(field)
				if !ok {
					continue
				}
				hybridIPInfo.Data[field] = val
				if replaceVal, ok := result.ReplaceFields[field]; ok {
					hybridIPInfo.Data[field] = replaceVal
				}
			}
		}
	}

	if h.OperateChain != nil {
		if err := h.OperateChain.Do(hybridIPInfo); err != nil {
			return nil, err
		}
	}

	return hybridIPInfo, nil
}

type HybridReaderOption struct {
	Mode string
}

// SetOption configures the HybridReader with the provided option, particularly the operational mode.
func (h *HybridReader) SetOption(option interface{}) error {
	if opt, ok := option.(HybridReaderOption); ok {
		h.hybridMode = opt.Mode
	}
	return nil
}

// Close method ensures that all underlying database readers are properly closed.
func (h *HybridReader) Close() error {
	for _, reader := range h.dbReaders {
		if err := reader.Close(); err != nil {
			return err
		}
	}
	return nil
}
