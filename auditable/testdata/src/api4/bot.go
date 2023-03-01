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

	var foo interface {
		Auditable() map[string]interface{}
	}

	bot := &model.Bot{}
	bot.Patch(botPatch)

	botValue := model.Bot{}
	botPatchValue := model.BotPatch{}

	botSlice := []*model.Bot{}
	botPatchSlice := []*model.BotPatch{}

	slice := []string{"a", "b", "c"}

	botMap := map[string]*model.Bot{}
	botPatchMap := map[string]*model.BotPatch{}

	props := map[string]any{}

	auditRec := c.MakeAuditRecord("createBot", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddEventParameter("bot", bot)
	auditRec.AddEventParameter("patch", botPatch) // want `model.BotPatch is not auditable, but it is added to the audit record`
	auditRec.AddEventParameter("bot value", botValue)
	auditRec.AddEventParameter("patch value", botPatchValue) // want `model.BotPatch is not auditable, but it is added to the audit record`
	auditRec.AddEventParameter("smth", "something else")
	auditRec.AddEventParameter("bot slice", botSlice)
	auditRec.AddEventParameter("bot patch slice", botPatchSlice) // want `model.BotPatch is not auditable, but it is added to the audit record`
	auditRec.AddEventParameter("bot map", botMap)
	auditRec.AddEventParameter("bot patch map", botPatchMap)      // want `model.BotPatch is not auditable, but it is added to the audit record`
	auditRec.AddEventParameter("string interface", props)         // want `is not auditable, but it is added to the audit record`
	auditRec.AddEventParameter("string interface ignore ", props) // auditable:ignore
	auditRec.AddEventParameter("anonymous interface ", foo)
	auditRec.AddEventParameter("slice", slice)
	other := "other"
	auditRec.AddEventParameter("other", other)

	auditRec.Success()
	auditRec.AddEventObjectType("bot")
	auditRec.AddEventResultState(bot) // overwrite meta

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bot); err != nil {
		return
	}
}
