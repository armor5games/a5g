package a5glogs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func DirectoryAndLogName(logPath string, fileExcludeRegexp *regexp.Regexp) (
	string, string, error) {
	var (
		logDirName  string
		logFileName string
		err         error
	)
	fileInfo, err := os.Stat(logPath)
	if err != nil && os.IsNotExist(err) {
		var (
			x = filepath.Dir(logPath)
			y = filepath.Base(logPath)
		)
		if strings.HasSuffix(x, y) {
			// This mean "path" contains an directory (without file name).
			logDirName = logPath
		} else {
			// This mean "path" contains an directory and file name.
			logDirName = x
			logFileName = y
		}
	} else if err != nil {
		return "", "", errors.WithStack(err)
	} else {
		if fileMode := fileInfo.Mode(); fileMode.IsDir() {
			logDirName = logPath
		} else {
			var (
				x = filepath.Dir(logPath)
				y = filepath.Base(logPath)
			)
			logDirName = x
			logFileName = y
		}
	}
	logFileName = path.Clean(logFileName)
	if fileExcludeRegexp != nil {
		if fileExcludeRegexp.MatchString(logFileName) {
			logFileName = ""
		}
	}
	return path.Clean(logDirName), logFileName, nil
}

// PartitionedPathByUserID produces path "000/000/001" for id 1.
// You may use this to generate path like:
// "/var/log/appname/123/456/789/123456789/your.log".
// <https://gist.github.com/xlab/6e204ef96b4433a697b3>
func PartitionedPathByUserID(i int64) (string, error) {
	const (
		chunksFormat = "%09d"
		chunkSize    = 3
	)
	var (
		a            = []byte(fmt.Sprintf(chunksFormat, i))
		allChunks    = make([][]byte, 0, len(a)/chunkSize+1)
		currentChunk []byte
	)
	for len(a) >= chunkSize {
		currentChunk, a = a[:chunkSize], a[chunkSize:]
		allChunks = append(allChunks, currentChunk)
	}
	if len(a) > 0 {
		allChunks = append(allChunks, a[:])
	}
	var stringChunks []string
	for _, a := range allChunks {
		stringChunks = append(stringChunks, string(a))
	}
	return strings.Join(stringChunks, "/"), nil
}
