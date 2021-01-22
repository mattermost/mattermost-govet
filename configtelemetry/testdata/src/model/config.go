// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

const (
	SERVICE_SETTINGS_DEFAULT_SITE_URL           = "http://localhost:8065"
	SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE      = ""
	SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE       = ""
	SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT       = 300
	SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT      = 300
	SERVICE_SETTINGS_DEFAULT_IDLE_TIMEOUT       = 60
	SERVICE_SETTINGS_DEFAULT_MAX_LOGIN_ATTEMPTS = 10
	SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM    = ""
	SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS = ":8065"
	SERVICE_SETTINGS_DEFAULT_GFYCAT_API_KEY     = "2_KtH_W5"
	SERVICE_SETTINGS_DEFAULT_GFYCAT_API_SECRET  = "3wLVZPiswc3DnaiaFoLkDvB4X0IV6CpMkj4tf2inJRsBY6-FnkT08zGmppWFgeof"

	LOCAL_MODE_SOCKET_PATH = "/var/tmp/mattermost_local.socket"
)

type ServiceSettings struct {
	SiteURL                                           *string  `access:"environment,authentication,write_restrictable"`
	WebsocketURL                                      *string  `access:"write_restrictable,cloud_restrictable"`
	LicenseFileLocation                               *string  `access:"write_restrictable,cloud_restrictable"`             // telemetry: none
	ListenAddress                                     *string  `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
	ConnectionSecurity                                *string  `access:"environment,write_restrictable,cloud_restrictable"`
	TLSCertFile                                       *string  `access:"environment,write_restrictable,cloud_restrictable"`
	TLSKeyFile                                        *string  `access:"environment,write_restrictable,cloud_restrictable"`
	TLSMinVer                                         *string  `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	TLSStrictTransport                                *bool    `access:"write_restrictable,cloud_restrictable"`
	TLSStrictTransportMaxAge                          *int64   `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	TLSOverwriteCiphers                               []string `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	UseLetsEncrypt                                    *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	LetsEncryptCertificateCacheFile                   *string  `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
	Forward80To443                                    *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	TrustedProxyIPHeader                              []string `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	ReadTimeout                                       *int     `access:"environment,write_restrictable,cloud_restrictable"`
	WriteTimeout                                      *int     `access:"environment,write_restrictable,cloud_restrictable"`
	IdleTimeout                                       *int     `access:"write_restrictable,cloud_restrictable"`
	MaximumLoginAttempts                              *int     `access:"authentication,write_restrictable,cloud_restrictable"`
	GoroutineHealthThreshold                          *int     `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	GoogleDeveloperKey                                *string  `access:"site,write_restrictable,cloud_restrictable"`
	EnableOAuthServiceProvider                        *bool    `access:"integrations"`
	EnableIncomingWebhooks                            *bool    `access:"integrations"`
	EnableOutgoingWebhooks                            *bool    `access:"integrations"`
	EnableCommands                                    *bool    `access:"integrations"`
	DEPRECATED_DO_NOT_USE_EnableOnlyAdminIntegrations *bool    `json:"EnableOnlyAdminIntegrations" mapstructure:"EnableOnlyAdminIntegrations"` // This field is deprecated and must not be used.
	EnablePostUsernameOverride                        *bool    `access:"integrations"`
	EnablePostIconOverride                            *bool    `access:"integrations"`
	EnableLinkPreviews                                *bool    `access:"site"` // telemetry: none
	EnableTesting                                     *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	EnableDeveloper                                   *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	EnableOpenTracing                                 *bool    `access:"write_restrictable,cloud_restrictable"`
	EnableSecurityFixAlert                            *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	EnableInsecureOutgoingConnections                 *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	AllowedUntrustedInternalConnections               *string  `access:"environment,write_restrictable,cloud_restrictable"`
	EnableMultifactorAuthentication                   *bool    `access:"authentication"`
	EnforceMultifactorAuthentication                  *bool    `access:"authentication"`
	EnableUserAccessTokens                            *bool    `access:"integrations"`
	AllowCorsFrom                                     *string  `access:"integrations,write_restrictable,cloud_restrictable"`
	CorsExposedHeaders                                *string  `access:"integrations,write_restrictable,cloud_restrictable"`
	CorsAllowCredentials                              *bool    `access:"integrations,write_restrictable,cloud_restrictable"`
	CorsDebug                                         *bool    `access:"integrations,write_restrictable,cloud_restrictable"`
	AllowCookiesForSubdomains                         *bool    `access:"write_restrictable,cloud_restrictable"`
	ExtendSessionLengthWithActivity                   *bool    `access:"environment,write_restrictable,cloud_restrictable"`
	SessionLengthWebInDays                            *int     `access:"environment,write_restrictable,cloud_restrictable"`
	SessionLengthMobileInDays                         *int     `access:"environment,write_restrictable,cloud_restrictable"`
	SessionLengthSSOInDays                            *int     `access:"environment,write_restrictable,cloud_restrictable"`
	SessionCacheInMinutes                             *int     `access:"environment,write_restrictable,cloud_restrictable"`
	SessionIdleTimeoutInMinutes                       *int     `access:"environment,write_restrictable,cloud_restrictable"`
	WebsocketSecurePort                               *int     `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	WebsocketPort                                     *int     `access:"write_restrictable,cloud_restrictable"` // telemetry: none
	WebserverMode                                     *string  `access:"environment,write_restrictable,cloud_restrictable"`
	EnableCustomEmoji                                 *bool    `access:"site"`
	EnableEmojiPicker                                 *bool    `access:"site"`
	EnableGifPicker                                   *bool    `access:"integrations"`
	GfycatApiKey                                      *string  `access:"integrations"`
	GfycatApiSecret                                   *string  `access:"integrations"`
	DEPRECATED_DO_NOT_USE_RestrictCustomEmojiCreation *string  `json:"RestrictCustomEmojiCreation" mapstructure:"RestrictCustomEmojiCreation"` // This field is deprecated and must not be used.
	DEPRECATED_DO_NOT_USE_RestrictPostDelete          *string  `json:"RestrictPostDelete" mapstructure:"RestrictPostDelete"`                   // This field is deprecated and must not be used.
	DEPRECATED_DO_NOT_USE_AllowEditPost               *string  `json:"AllowEditPost" mapstructure:"AllowEditPost"`                             // This field is deprecated and must not be used.
	PostEditTimeLimit                                 *int     `access:"user_management_permissions"`
	TimeBetweenUserTypingUpdatesMilliseconds          *int64   `access:"experimental,write_restrictable,cloud_restrictable"`
	EnablePostSearch                                  *bool    `access:"write_restrictable,cloud_restrictable"`
	MinimumHashtagLength                              *int     `access:"environment,write_restrictable,cloud_restrictable"`
	EnableUserTypingMessages                          *bool    `access:"experimental,write_restrictable,cloud_restrictable"`
	EnableChannelViewedMessages                       *bool    `access:"experimental,write_restrictable,cloud_restrictable"`
	EnableUserStatuses                                *bool    `access:"write_restrictable,cloud_restrictable"`
	ExperimentalEnableAuthenticationTransfer          *bool    `access:"experimental,write_restrictable,cloud_restrictable"`
	ClusterLogTimeoutMilliseconds                     *int     `access:"write_restrictable,cloud_restrictable"`
	CloseUnusedDirectMessages                         *bool    `access:"experimental"`
	EnablePreviewFeatures                             *bool    `access:"experimental"`
	EnableTutorial                                    *bool    `access:"experimental"`
	ExperimentalEnableDefaultChannelLeaveJoinMessages *bool    `access:"experimental"`
	ExperimentalGroupUnreadChannels                   *string  `access:"experimental"`
	ExperimentalChannelOrganization                   *bool    `access:"experimental"`
	DEPRECATED_DO_NOT_USE_ImageProxyType              *string  `json:"ImageProxyType" mapstructure:"ImageProxyType"`       // This field is deprecated and must not be used.
	DEPRECATED_DO_NOT_USE_ImageProxyURL               *string  `json:"ImageProxyURL" mapstructure:"ImageProxyURL"`         // This field is deprecated and must not be used.
	DEPRECATED_DO_NOT_USE_ImageProxyOptions           *string  `json:"ImageProxyOptions" mapstructure:"ImageProxyOptions"` // This field is deprecated and must not be used.
	EnableAPITeamDeletion                             *bool
	EnableAPIUserDeletion                             *bool
	ExperimentalEnableHardenedMode                    *bool `access:"experimental"`
	DisableLegacyMFA                                  *bool `access:"write_restrictable,cloud_restrictable"`
	ExperimentalStrictCSRFEnforcement                 *bool `access:"experimental,write_restrictable,cloud_restrictable"`
	EnableEmailInvitations                            *bool `access:"authentication"`
	DisableBotsWhenOwnerIsDeactivated                 *bool `access:"integrations,write_restrictable,cloud_restrictable"`
	EnableBotAccountCreation                          *bool `access:"integrations"`
	EnableSVGs                                        *bool `access:"site"`
	EnableLatex                                       *bool `access:"site"`
	EnableAPIChannelDeletion                          *bool
	EnableLocalMode                                   *bool
	LocalModeSocketLocation                           *string // telemetry: none
	EnableAWSMetering                                 *bool   // telemetry: none
	SplitKey                                          *string `access:"environment,write_restrictable"` // telemetry: none
	FeatureFlagSyncIntervalSeconds                    *int    `access:"environment,write_restrictable"` // telemetry: none
	DebugSplit                                        *bool   `access:"environment,write_restrictable"` // telemetry: none
	ThreadAutoFollow                                  *bool   `access:"experimental"`                   // telemetry: none
	CollapsedThreads                                  *string `access:"experimental"`
	ManagedResourcePaths                              *string `access:"environment,write_restrictable,cloud_restrictable"`
	EnableLegacySidebar                               *bool   `access:"experimental"`
}

type MessageExportSettings struct {
	EnableExport          *bool   `access:"compliance"`
	ExportFormat          *string `access:"compliance"`
	DailyRunTime          *string `access:"compliance"`
	ExportFromTimestamp   *int64  `access:"compliance"`
	BatchSize             *int    `access:"compliance"`
	DownloadExportResults *bool   `access:"compliance"`

	// formatter-specific settings - these are only expected to be non-nil if ExportFormat is set to the associated format
	GlobalRelaySettings *GlobalRelayMessageExportSettings `access:"compliance"`
}

type GlobalRelayMessageExportSettings struct {
	CustomerType      *string `access:"compliance"` // must be either A9 or A10, dictates SMTP server url
	SmtpUsername      *string `access:"compliance"`
	SmtpPassword      *string `access:"compliance"`
	EmailAddress      *string `access:"compliance"` // the address to send messages to
	SMTPServerTimeout *int    `access:"compliance"`
}

type CloudSettings struct {
	CWSUrl *string `access:"environment,write_restrictable"`
}

type Config struct {
	ServiceSettings       ServiceSettings
	MessageExportSettings MessageExportSettings
	CloudSettings         CloudSettings // telemetry: none
}
