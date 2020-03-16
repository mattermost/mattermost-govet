// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"context"
	"net/http"
)

func (a *API) InitUsers() {
	a.BaseRoutes.Users.Handle("/ids", a.ApiSessionRequired(getUsersByIds)).Methods("POST") // want `Cannot find /userzs/ids method: POST in OpenAPI 3 spec. \(maybe you meant: \[/users/ids\]\)`
}
func getUsersByIds(c *context.Context, w http.ResponseWriter, r *http.Request) {
}
