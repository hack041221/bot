package bot

import (
	"regexp"
)

var r = regexp.MustCompile("(?:\\/v\\/|watch\\/|\\?v=|\\&v=|youtu\\.be\\/|\\/v=|^youtu\\.be\\/|watch\\%3Fv\\%3D)([a-zA-Z0-9_-]{11})+")

func hasYoutubeLink(text string) (links []string) {
	for _, p := range r.FindAllStringSubmatch(text, -1) {
		links = append(links, p[1])
	}
	return
}
