package validate

import (
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,79}$`)
var reserved = map[string]bool{"api": true, "admin": true, "d": true, "themes": true, "uploads": true, "static": true, "feed": true, "sitemap": true, "robots": true, "healthz": true, "install": true, "login": true}

func Slug(s string) bool     { return slugPattern.MatchString(s) && !reserved[strings.ToLower(s)] }
func Username(s string) bool { ok, _ := regexp.MatchString(`^[a-z0-9_]{3,50}$`, s); return ok }
