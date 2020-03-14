// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package apiAuditLogs

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "apiAuditLogs",
	Doc:  "check the audit logs usage in the API",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	whiteList := map[string]bool{
		"autocompleteChannelsForTeam":          true,
		"autocompleteChannelsForTeamForSearch": true,
		"autocompleteEmojis":                   true,
		"autocompleteUsers":                    true,
		"channelMembersMinusGroupMembers":      true,
		"checkUserMfa":                         true,
		"connectWebSocket":                     true,
		"createEphemeralPost":                  true,
		"deleteReaction":                       true,
		"disableBot":                           true,
		"doPostAction":                         true,
		"enableBot":                            true,
		"executeCommand":                       true,
		"generateMfaSecret":                    true,
		"getAllChannels":                       true,
		"getAllTeams":                          true,
		"getAnalytics":                         true,
		"getAudits":                            true,
		"getAuthorizedOAuthApps":               true,
		"getBot":                               true,
		"getBotIconImage":                      true,
		"getBots":                              true,
		"getBrandImage":                        true,
		"getBulkReactions":                     true,
		"getChannel":                           true,
		"getChannelByName":                     true,
		"getChannelByNameForTeamName":          true,
		"getChannelMember":                     true,
		"getChannelMembers":                    true,
		"getChannelMembersByIds":               true,
		"getChannelMembersForUser":             true,
		"getChannelMembersTimezones":           true,
		"getChannelModerations":                true,
		"getChannelsForScheme":                 true,
		"getChannelsForTeamForUser":            true,
		"getChannelStats":                      true,
		"getChannelUnread":                     true,
		"getClientConfig":                      true,
		"getClientLicense":                     true,
		"getClusterStatus":                     true,
		"getCommand":                           true,
		"getComplianceReport":                  true,
		"getComplianceReports":                 true,
		"getConfig":                            true,
		"getDefaultProfileImage":               true,
		"getDeletedChannelsForTeam":            true,
		"getEmoji":                             true,
		"getEmojiByName":                       true,
		"getEmojiImage":                        true,
		"getEmojiList":                         true,
		"getEnvironmentConfig":                 true,
		"getFile":                              true,
		"getFileInfo":                          true,
		"getFileInfosForPost":                  true,
		"getFileLink":                          true,
		"getFilePreview":                       true,
		"getFileThumbnail":                     true,
		"getFlaggedPostsForUser":               true,
		"getGroup":                             true,
		"getGroupMembers":                      true,
		"getGroups":                            true,
		"getGroupsByChannel":                   true,
		"getGroupsByTeam":                      true,
		"getGroupSyncable":                     true,
		"getGroupSyncables":                    true,
		"getImage":                             true,
		"getIncomingHooks":                     true,
		"getInviteInfo":                        true,
		"getJob":                               true,
		"getJobs":                              true,
		"getJobsByType":                        true,
		"getLatestTermsOfService":              true,
		"getLdapGroups":                        true,
		"getLogs":                              true,
		"getMarketplacePlugins":                true,
		"getOAuthApp":                          true,
		"getOAuthAppInfo":                      true,
		"getOAuthApps":                         true,
		"getOpenGraphMetadata":                 true,
		"getOutgoingHooks":                     true,
		"getPinnedPosts":                       true,
		"getPlugins":                           true,
		"getPluginStatuses":                    true,
		"getPolicy":                            true,
		"getPost":                              true,
		"getPostsForChannel":                   true,
		"getPostsForChannelAroundLastUnread":   true,
		"getPostThread":                        true,
		"getPreferenceByCategoryAndName":       true,
		"getPreferences":                       true,
		"getPreferencesByCategory":             true,
		"getProfileImage":                      true,
		"getPublicChannelsByIdsForTeam":        true,
		"getPublicChannelsForTeam":             true,
		"getPublicFile":                        true,
		"getReactions":                         true,
		"getRedirectLocation":                  true,
		"getRole":                              true,
		"getRoleByName":                        true,
		"getRolesByNames":                      true,
		"getSamlCertificateStatus":             true,
		"getSamlMetadata":                      true,
		"getSamlMetadataFromIdp":               true,
		"getScheme":                            true,
		"getSchemes":                           true,
		"getServerBusyExpires":                 true,
		"getSessions":                          true,
		"getSupportedTimezones":                true,
		"getSystemPing":                        true,
		"getTeam":                              true,
		"getTeamByName":                        true,
		"getTeamIcon":                          true,
		"getTeamMember":                        true,
		"getTeamMembers":                       true,
		"getTeamMembersByIds":                  true,
		"getTeamMembersForUser":                true,
		"getTeamsForScheme":                    true,
		"getTeamsForUser":                      true,
		"getTeamStats":                         true,
		"getTeamsUnreadForUser":                true,
		"getTeamUnread":                        true,
		"getTotalUsersStats":                   true,
		"getUser":                              true,
		"getUserAccessToken":                   true,
		"getUserAccessTokens":                  true,
		"getUserAccessTokensForUser":           true,
		"getUserAudits":                        true,
		"getUserByEmail":                       true,
		"getUserByUsername":                    true,
		"getUsers":                             true,
		"getUsersByGroupChannelIds":            true,
		"getUsersByIds":                        true,
		"getUsersByNames":                      true,
		"getUserStatus":                        true,
		"getUserStatusesByIds":                 true,
		"getUserTermsOfService":                true,
		"getWebappPlugins":                     true,
		"listAutocompleteCommands":             true,
		"listCommands":                         true,
		"logout":                               true,
		"openDialog":                           true,
		"patchChannelModerations":              true,
		"pinPost":                              true,
		"postLog":                              true,
		"pushNotificationAck":                  true,
		"saveReaction":                         true,
		"searchAllChannels":                    true,
		"searchArchivedChannelsForTeam":        true,
		"searchChannelsForTeam":                true,
		"searchEmojis":                         true,
		"searchGroupChannels":                  true,
		"searchPosts":                          true,
		"searchTeams":                          true,
		"searchUserAccessTokens":               true,
		"searchUsers":                          true,
		"setPostUnread":                        true,
		"submitDialog":                         true,
		"teamExists":                           true,
		"teamMembersMinusGroupMembers":         true,
		"testElasticsearch":                    true,
		"testEmail":                            true,
		"testLdap":                             true,
		"testS3":                               true,
		"testSiteURL":                          true,
		"unpinPost":                            true,
		"updateUserStatus":                     true,
		"uploadFileStream":                     true,
		"viewChannel":                          true,
	}

	if pass.Pkg.Path() != "github.com/mattermost/mattermost-server/v5/api4" {
		return nil, nil
	}

	for _, file := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		if strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "apitestlib.go") {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch funDecl := node.(type) {
			case *ast.FuncDecl:
				if whiteList[funDecl.Name.Name] {
					return false
				}
				if !isApiHandler(funDecl) {
					return false
				}
				initializationFound := false
				logCallFound := false
				successCallFound := false
				ast.Inspect(funDecl, func(node ast.Node) bool {
					switch n := node.(type) {
					case *ast.CallExpr:
						fun, ok := n.Fun.(*ast.SelectorExpr)
						if !ok {
							return true
						}

						ident, ok := fun.X.(*ast.Ident)
						if !ok {
							return true
						}

						// must have a auditRec.Success()
						if ident.Name == "auditRec" && fun.Sel.Name == "Success" {
							successCallFound = true
						}
						if ident.Name == "c" && (fun.Sel.Name == "LogAuditRec" || fun.Sel.Name == "LogAuditRecWithLevel") {
							logCallFound = true
						}

						if ident.Name == "c" && fun.Sel.Name == "MakeAuditRecord" {
							initializationFound = true
							firstArg, ok := n.Args[0].(*ast.BasicLit)
							if !ok {
								pass.Reportf(n.Args[0].Pos(), "Invalid record name, expected \"%s\", found \"%v\"", funDecl.Name.Name, n.Args[0])
								return true
							}
							secondArg, ok := n.Args[1].(*ast.SelectorExpr)
							if !ok {
								pass.Reportf(n.Args[1].Pos(), "Invalid initial state for record, expected \"audit.Fail\", found \"%v\"", n.Args[1])
								return true
							}
							if firstArg.Kind != token.STRING || firstArg.Value != fmt.Sprintf("\"%s\"", funDecl.Name.Name) {
								pass.Reportf(n.Args[0].Pos(), "Invalid record name, expected \"%s\", found %s", funDecl.Name.Name, firstArg.Value)
								return true
							}
							secondArgX, ok := secondArg.X.(*ast.Ident)
							if !ok {
								pass.Reportf(n.Args[1].Pos(), "Invalid initial state for record, expected \"audit.Fail\", found \"%v\"", secondArg.X)
								return true
							}
							if secondArgX.Name != "audit" || secondArg.Sel.Name != "Fail" {
								pass.Reportf(n.Args[1].Pos(), "Invalid initial state for record, expected \"audit.Fail\", found \"%s.%s\"", secondArgX.Name, secondArg.Sel.Name)
								return true
							}
						}
					}
					return true
				})
				if !initializationFound {
					pass.Reportf(funDecl.Pos(), "Expected audit log in this function, but not found, please add the audit logs or add the \"%s\" function to the white list", funDecl.Name.Name)
				} else {
					if !logCallFound {
						pass.Reportf(funDecl.Pos(), "Expected LogAuditRec or LogAuditRecWithLevel call, but not found")
					}
					if !successCallFound {
						pass.Reportf(funDecl.Pos(), "Expected Success call, but not found")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func isApiHandler(funDecl *ast.FuncDecl) bool {
	funcType := funDecl.Type
	if len(funcType.Params.List) < 3 {
		return false
	}
	arg0Type, ok := funcType.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	arg0X, ok := arg0Type.X.(*ast.Ident)
	if !ok {
		return false
	}
	if arg0X.Name != "Context" {
		return false
	}

	arg1Type, ok := funcType.Params.List[1].Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	arg1X, ok := arg1Type.X.(*ast.Ident)
	if !ok {
		return false
	}
	if arg1X.Name != "http" || arg1Type.Sel.Name != "ResponseWriter" {
		return false
	}

	arg2Type, ok := funcType.Params.List[2].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	arg2X, ok := arg2Type.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	arg2XX, ok := arg2X.X.(*ast.Ident)
	if !ok {
		return false
	}

	if arg2XX.Name != "http" || arg2X.Sel.Name != "Request" {
		return false
	}
	return true
}
