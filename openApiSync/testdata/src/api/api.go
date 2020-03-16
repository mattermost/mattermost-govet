// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api/v4'

	Users *mux.Router // 'api/v4/userzs'
}

type API struct {
	BaseRoutes *Routes
}

func (*API) ApiSessionRequired(h func(*context.Context, http.ResponseWriter, *http.Request)) http.Handler {
	return nil
}
func Init(root *mux.Router) *API {
	api := &API{
		BaseRoutes: &Routes{},
	}
	api.BaseRoutes.Root = root
	api.BaseRoutes.ApiRoot = root.PathPrefix("api/v4").Subrouter()

	api.BaseRoutes.Users = api.BaseRoutes.ApiRoot.PathPrefix("/users").Subrouter() // want "PathPrefix doesn't match field comment for field 'Users': 'api/v4/users' vs 'api/v4/userzs'"

	api.InitUsers()

	fmt.Errorf("") // wzant "result of fmt.Errorf call not used"
	return api
}

func _() {
	_ = Init(nil)
}
