// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"audit"
	"model"
)

func createBot(c *Context, w http.ResponseWriter, r *http.Request) {
	var botPatch *model.BotPatch
	err := json.NewDecoder(r.Body).Decode(&botPatch)
	if err != nil {
		return
	}

	bot := &model.Bot{}
	bot.Patch(botPatch)

	auditRec := c.MakeAuditRecord("createBot", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddEventParameter("bot", bot)

	auditRec.Success()
	auditRec.AddEventObjectType("bot")
	auditRec.AddEventResultState(bot) // overwrite meta

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bot); err != nil {
		return
	}
}
