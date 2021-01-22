// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package telemetry

import (
	"model"
)

const (
	TrackConfigService       = "config_service"
	TrackConfigMessageExport = "config_message_export"
)

type ServerIface interface {
	Config() *model.Config
}

type TelemetryService struct {
	srv ServerIface
}

func (ts *TelemetryService) sendTelemetry(event string, properties map[string]interface{}) {
}

func isDefault(setting interface{}, defaultValue interface{}) bool {
	return setting == defaultValue
}

func (ts *TelemetryService) trackConfig() {
	cfg := ts.srv.Config()
	ts.sendTelemetry(TrackConfigService, map[string]interface{}{
		"web_server_mode":                                         *cfg.ServiceSettings.WebserverMode,
		"enable_security_fix_alert":                               *cfg.ServiceSettings.EnableSecurityFixAlert,
		"enable_insecure_outgoing_connections":                    *cfg.ServiceSettings.EnableInsecureOutgoingConnections,
		"enable_incoming_webhooks":                                cfg.ServiceSettings.EnableIncomingWebhooks,
		"enable_outgoing_webhooks":                                cfg.ServiceSettings.EnableOutgoingWebhooks,
		"enable_commands":                                         *cfg.ServiceSettings.EnableCommands,
		"enable_only_admin_integrations":                          *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations,
		"enable_post_username_override":                           cfg.ServiceSettings.EnablePostUsernameOverride,
		"enable_post_icon_override":                               cfg.ServiceSettings.EnablePostIconOverride,
		"enable_user_access_tokens":                               *cfg.ServiceSettings.EnableUserAccessTokens,
		"enable_custom_emoji":                                     *cfg.ServiceSettings.EnableCustomEmoji,
		"enable_emoji_picker":                                     *cfg.ServiceSettings.EnableEmojiPicker,
		"enable_gif_picker":                                       *cfg.ServiceSettings.EnableGifPicker,
		"gfycat_api_key":                                          isDefault(*cfg.ServiceSettings.GfycatApiKey, model.SERVICE_SETTINGS_DEFAULT_GFYCAT_API_KEY),
		"gfycat_api_secret":                                       isDefault(*cfg.ServiceSettings.GfycatApiSecret, model.SERVICE_SETTINGS_DEFAULT_GFYCAT_API_SECRET),
		"experimental_enable_authentication_transfer":             *cfg.ServiceSettings.ExperimentalEnableAuthenticationTransfer,
		"restrict_custom_emoji_creation":                          *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_RestrictCustomEmojiCreation,
		"enable_testing":                                          cfg.ServiceSettings.EnableTesting,
		"enable_developer":                                        *cfg.ServiceSettings.EnableDeveloper,
		"enable_multifactor_authentication":                       *cfg.ServiceSettings.EnableMultifactorAuthentication,
		"enforce_multifactor_authentication":                      *cfg.ServiceSettings.EnforceMultifactorAuthentication,
		"enable_oauth_service_provider":                           cfg.ServiceSettings.EnableOAuthServiceProvider,
		"connection_security":                                     *cfg.ServiceSettings.ConnectionSecurity,
		"tls_strict_transport":                                    *cfg.ServiceSettings.TLSStrictTransport,
		"uses_letsencrypt":                                        *cfg.ServiceSettings.UseLetsEncrypt,
		"forward_80_to_443":                                       *cfg.ServiceSettings.Forward80To443,
		"maximum_login_attempts":                                  *cfg.ServiceSettings.MaximumLoginAttempts,
		"extend_session_length_with_activity":                     *cfg.ServiceSettings.ExtendSessionLengthWithActivity,
		"session_length_web_in_days":                              *cfg.ServiceSettings.SessionLengthWebInDays,
		"session_length_mobile_in_days":                           *cfg.ServiceSettings.SessionLengthMobileInDays,
		"session_length_sso_in_days":                              *cfg.ServiceSettings.SessionLengthSSOInDays,
		"session_cache_in_minutes":                                *cfg.ServiceSettings.SessionCacheInMinutes,
		"session_idle_timeout_in_minutes":                         *cfg.ServiceSettings.SessionIdleTimeoutInMinutes,
		"isdefault_site_url":                                      isDefault(*cfg.ServiceSettings.SiteURL, model.SERVICE_SETTINGS_DEFAULT_SITE_URL),
		"isdefault_tls_cert_file":                                 isDefault(*cfg.ServiceSettings.TLSCertFile, model.SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE),
		"isdefault_tls_key_file":                                  isDefault(*cfg.ServiceSettings.TLSKeyFile, model.SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE),
		"isdefault_read_timeout":                                  isDefault(*cfg.ServiceSettings.ReadTimeout, model.SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT),
		"isdefault_write_timeout":                                 isDefault(*cfg.ServiceSettings.WriteTimeout, model.SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT),
		"isdefault_idle_timeout":                                  isDefault(*cfg.ServiceSettings.IdleTimeout, model.SERVICE_SETTINGS_DEFAULT_IDLE_TIMEOUT),
		"isdefault_google_developer_key":                          isDefault(cfg.ServiceSettings.GoogleDeveloperKey, ""),
		"isdefault_allow_cors_from":                               isDefault(*cfg.ServiceSettings.AllowCorsFrom, model.SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM),
		"isdefault_cors_exposed_headers":                          isDefault(cfg.ServiceSettings.CorsExposedHeaders, ""),
		"cors_allow_credentials":                                  *cfg.ServiceSettings.CorsAllowCredentials,
		"cors_debug":                                              *cfg.ServiceSettings.CorsDebug,
		"isdefault_allowed_untrusted_internal_connections":        isDefault(*cfg.ServiceSettings.AllowedUntrustedInternalConnections, ""),
		"restrict_post_delete":                                    *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_RestrictPostDelete,
		"allow_edit_post":                                         *cfg.ServiceSettings.DEPRECATED_DO_NOT_USE_AllowEditPost,
		"post_edit_time_limit":                                    *cfg.ServiceSettings.PostEditTimeLimit,
		"enable_user_typing_messages":                             *cfg.ServiceSettings.EnableUserTypingMessages,
		"enable_channel_viewed_messages":                          *cfg.ServiceSettings.EnableChannelViewedMessages,
		"time_between_user_typing_updates_milliseconds":           *cfg.ServiceSettings.TimeBetweenUserTypingUpdatesMilliseconds,
		"cluster_log_timeout_milliseconds":                        *cfg.ServiceSettings.ClusterLogTimeoutMilliseconds,
		"enable_post_search":                                      *cfg.ServiceSettings.EnablePostSearch,
		"minimum_hashtag_length":                                  *cfg.ServiceSettings.MinimumHashtagLength,
		"enable_user_statuses":                                    *cfg.ServiceSettings.EnableUserStatuses,
		"close_unused_direct_messages":                            *cfg.ServiceSettings.CloseUnusedDirectMessages,
		"enable_preview_features":                                 *cfg.ServiceSettings.EnablePreviewFeatures,
		"enable_tutorial":                                         *cfg.ServiceSettings.EnableTutorial,
		"experimental_enable_default_channel_leave_join_messages": *cfg.ServiceSettings.ExperimentalEnableDefaultChannelLeaveJoinMessages,
		"experimental_group_unread_channels":                      *cfg.ServiceSettings.ExperimentalGroupUnreadChannels,
		"collapsed_threads":                                       *cfg.ServiceSettings.CollapsedThreads,
		"websocket_url":                                           isDefault(*cfg.ServiceSettings.WebsocketURL, ""),
		"allow_cookies_for_subdomains":                            *cfg.ServiceSettings.AllowCookiesForSubdomains,
		"enable_api_team_deletion":                                *cfg.ServiceSettings.EnableAPITeamDeletion,
		"enable_api_user_deletion":                                *cfg.ServiceSettings.EnableAPIUserDeletion,
		"enable_api_channel_deletion":                             *cfg.ServiceSettings.EnableAPIChannelDeletion,
		"experimental_enable_hardened_mode":                       *cfg.ServiceSettings.ExperimentalEnableHardenedMode,
		"disable_legacy_mfa":                                      *cfg.ServiceSettings.DisableLegacyMFA,
		"experimental_strict_csrf_enforcement":                    *cfg.ServiceSettings.ExperimentalStrictCSRFEnforcement,
		"enable_email_invitations":                                *cfg.ServiceSettings.EnableEmailInvitations,
		"experimental_channel_organization":                       *cfg.ServiceSettings.ExperimentalChannelOrganization,
		"disable_bots_when_owner_is_deactivated":                  *cfg.ServiceSettings.DisableBotsWhenOwnerIsDeactivated,
		"enable_bot_account_creation":                             *cfg.ServiceSettings.EnableBotAccountCreation,
		"enable_svgs":                                             *cfg.ServiceSettings.EnableSVGs,
		"enable_latex":                                            *cfg.ServiceSettings.EnableLatex,
		"enable_opentracing":                                      *cfg.ServiceSettings.EnableOpenTracing,
		"enable_local_mode":                                       *cfg.ServiceSettings.EnableLocalMode,
		"managed_resource_paths":                                  isDefault(*cfg.ServiceSettings.ManagedResourcePaths, ""),
		"enable_legacy_sidebar":                                   *cfg.ServiceSettings.EnableLegacySidebar,
	})

	ts.sendTelemetry(TrackConfigMessageExport, map[string]interface{}{
		"enable_message_export":                 *cfg.MessageExportSettings.EnableExport,
		"export_format":                         *cfg.MessageExportSettings.ExportFormat,
		"daily_run_time":                        *cfg.MessageExportSettings.DailyRunTime,
		"default_export_from_timestamp":         *cfg.MessageExportSettings.ExportFromTimestamp,
		"batch_size":                            *cfg.MessageExportSettings.BatchSize,
		"global_relay_customer_type":            *cfg.MessageExportSettings.GlobalRelaySettings.CustomerType,
		"is_default_global_relay_smtp_username": isDefault(*cfg.MessageExportSettings.GlobalRelaySettings.SmtpUsername, ""),
		"is_default_global_relay_smtp_password": isDefault(*cfg.MessageExportSettings.GlobalRelaySettings.SmtpPassword, ""),
		"is_default_global_relay_email_address": isDefault(*cfg.MessageExportSettings.GlobalRelaySettings.EmailAddress, ""),
		"global_relay_smtp_server_timeout":      *cfg.MessageExportSettings.GlobalRelaySettings.SMTPServerTimeout,
		"download_export_results":               *cfg.MessageExportSettings.DownloadExportResults,
	})
}
