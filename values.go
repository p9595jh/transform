package transform

import "strings"

const (
	_TAG = "transform"

	// tags must be able to be seperated by comma(,)
	// ex. transform:"a,b,c:10"
	_SEPERATOR = ","
)

// 'lower' and 'upper' are provided as default
// these can be overwritten
var defaults = []I{
	{
		Name: "lower",
		F: F1(func(s string, _ string) string {
			return strings.ToLower(s)
		}),
	}, {
		Name: "upper",
		F: F1(func(s string, _ string) string {
			return strings.ToUpper(s)
		}),
	},
}
