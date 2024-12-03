// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package helpers

import "strings"

func IsIgnore(file string, ignoreFiles []string) bool {
	for _, f := range ignoreFiles {
		if strings.HasSuffix(file, f) {
			return true
		}
	}
	return false
}
