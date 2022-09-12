package mode

import (
	"regexp"
	"strings"
)

type CompilationMode int

const (
	// Independent is the default mode
	Independent CompilationMode = iota
	Extend
)

func CheckMode(comments []string) (CompilationMode, []string) {
	// 1. check Extend mode
	// compile reg expression
	reg := regexp.MustCompile(`@exclude extends ([^\s]+)`)
	if reg == nil {
		panic("regexp comile error")
	}
	// extract information
	for _, comment := range comments {
		result := reg.FindStringSubmatch(comment)

		if result != nil {
			return Extend, strings.Split(result[1], ",")
		}
	}

	// default mode
	return Independent, nil
}
