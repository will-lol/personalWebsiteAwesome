package merge

import (
	"strings"

	"github.com/a-h/templ"
)

func MergeAttrs(attrs1 templ.Attributes, attrs2 templ.Attributes) templ.Attributes {
	for k, v := range attrs2 {
		if attrs1[k] != nil {
			switch v.(type) {
			case string:
				attrs1[k] = strings.Join([]string{attrs1[k].(string), attrs2[k].(string)}, " ")
			}
		} else {
			attrs1[k] = v
		}
	}
	return attrs1
}
