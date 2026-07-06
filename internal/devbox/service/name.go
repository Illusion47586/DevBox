package service

import (
	"path/filepath"
	"strings"
)

func deriveProjectName(gitURL string) string {
	base := filepath.Base(strings.TrimSuffix(gitURL, "/"))
	base = strings.TrimSuffix(base, ".git")
	base = strings.ToLower(base)
	base = replaceInvalidNameRunes(base)
	if base == "" {
		return "project"
	}
	return base
}

func replaceInvalidNameRunes(value string) string {
	var builder strings.Builder
	lastHyphen := false
	for _, r := range value {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			builder.WriteRune(r)
			lastHyphen = false
			continue
		}
		if !lastHyphen {
			builder.WriteByte('-')
			lastHyphen = true
		}
	}
	return strings.Trim(builder.String(), "-")
}
