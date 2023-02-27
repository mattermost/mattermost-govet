// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

// BotPatch is a description of what fields to update on an existing bot.
type BotPatch struct {
	Username    *string `json:"username"`
	DisplayName *string `json:"display_name"`
	Description *string `json:"description"`
}

// Bot is a special type of User meant for programmatic interactions.
// Note that the primary key of a bot is the UserId, and matches the primary key of the
// corresponding user.
type Bot struct {
	UserId         string `json:"user_id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name,omitempty"`
	Description    string `json:"description,omitempty"`
	OwnerId        string `json:"owner_id"`
	LastIconUpdate int64  `json:"last_icon_update,omitempty"`
	CreateAt       int64  `json:"create_at"`
	UpdateAt       int64  `json:"update_at"`
	DeleteAt       int64  `json:"delete_at"`
}

func (b *Bot) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"user_id":          b.UserId,
		"username":         b.Username,
		"display_name":     b.DisplayName,
		"description":      b.Description,
		"owner_id":         b.OwnerId,
		"last_icon_update": b.LastIconUpdate,
		"create_at":        b.CreateAt,
		"update_at":        b.UpdateAt,
		"delete_at":        b.DeleteAt,
	}
}

// Patch modifies an existing bot with optional fields from the given patch.
// TODO 6.0: consider returning a boolean to indicate whether or not the patch
// applied any changes.
func (b *Bot) Patch(patch *BotPatch) {
	if patch.Username != nil {
		b.Username = *patch.Username
	}

	if patch.DisplayName != nil {
		b.DisplayName = *patch.DisplayName
	}

	if patch.Description != nil {
		b.Description = *patch.Description
	}
}
