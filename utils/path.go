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
