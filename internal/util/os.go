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

package util

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/pkg/errors"
)

// PrepareDir ensures that the specified directory path exists.
// If the directory does not exist, it attempts to create it.
func PrepareDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !stat.IsDir() {
		log.Debugf("%s is not a directory", path)
		return errors.ErrInvalidDirectory
	}
	return nil
}

// IsFileExist checks if a file exists at the given path.
func IsFileExist(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}
