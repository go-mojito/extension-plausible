package plausible

import (
	"regexp"
)

var uaRegex = regexp.MustCompile(`(?m)Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini`)

func guessWidthFromUA(userAgent string) int {
	if uaRegex.MatchString(userAgent) {
		return 360
	}
	return 1920
}
