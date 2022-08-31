package utils

import (
	"strings"
)

func SplitDirectoryAndFilename(generatedFilenamePrefix string) (path string, fileName string) {
	idx := strings.LastIndex(generatedFilenamePrefix, "/")
	path = generatedFilenamePrefix[:idx]
	fileName = generatedFilenamePrefix[idx+1:]
	return path, fileName
}

func LowerCammel(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
