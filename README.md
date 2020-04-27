## Mattermost Go Vet

This repository contains mattermost-specific go-vet rules that are used to maintain code consistency in `mattermost-server`.

## Included analyzers

1. **equalLenAsserts** - check for (require/assert).Equal(t, X, len(Y))
1. **inconsistentReceiverName** - check for inconsistent receiver names in the methods of a struct
1. **license** - check the license header
1. **openApiSync** - check for inconsistencies between OpenAPI spec and the source code
1. **structuredLogging** - check invalid usage of logging (must use structured logging)
1. **tFatal** - check invalid usage of t.Fatal assertions (instead of testify methods)
1. **apiAuditLogs** - check that audit records are properly created in the API layer
1. **rawSql** - check invalid usage of raw SQL queries instead of using the squirrel lib
