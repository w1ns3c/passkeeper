package entities

import "strings"

// HideDBpass return dbURL without password
func HideDBpass(dbURL string) string {
	var (
		prefix     = "://"
		hiddenPass = "******"
	)

	userInd := strings.Index(dbURL, prefix) + len(prefix)
	passInd1 := strings.Index(dbURL[userInd:], ":")
	passInd2 := strings.LastIndex(dbURL, "@")

	if userInd+passInd1 < passInd2 && passInd2 < len(dbURL) {
		return dbURL[:userInd+passInd1+1] + hiddenPass + dbURL[passInd2:]
	}

	return "can't correctly hide db password"
}
