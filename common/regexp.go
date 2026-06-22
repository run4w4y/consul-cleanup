package common

import "regexp"

var uuidRegexString = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

var uuidRegexp = regexp.MustCompile(uuidRegexString)

// extracts uuid (nomad alloc id) from a full service id in the consul service entry
func extractUUID(serviceId string) string {
	return uuidRegexp.FindString(serviceId)
}
