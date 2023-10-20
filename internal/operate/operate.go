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

package operate

import (
	"github.com/sjzar/ips/pkg/model"
)

// IPOperateFunc defines a function type for operations on IPInfo,
// such as filtering fields, rewriting data, supplementing data, etc.
type IPOperateFunc func(info *model.IPInfo) error

// IPOperateChain represents a chain of IPOperateFuncs,
// organized to be executed in the order they are added.
type IPOperateChain struct {
	operates []IPOperateFunc
}

// NewIPOperateChain initializes and returns a new IPOperateChain.
func NewIPOperateChain() *IPOperateChain {
	return &IPOperateChain{
		operates: make([]IPOperateFunc, 0),
	}
}

// Do executes all the operations in the chain on the given IPInfo.
func (c *IPOperateChain) Do(info *model.IPInfo) error {
	if c == nil || c.operates == nil || len(c.operates) == 0 || info == nil {
		return nil
	}

	for _, operate := range c.operates {
		if err := operate(info); err != nil {
			return err
		}
	}

	return nil
}

// Use adds an IPOperateFunc to the end of the operation chain.
func (c *IPOperateChain) Use(f IPOperateFunc) {
	if c == nil {
		panic("IPOperateChain is nil")
	}

	if c.operates == nil {
		c.operates = make([]IPOperateFunc, 0)
	}

	c.operates = append(c.operates, f)
}
