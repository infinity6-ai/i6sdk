package main_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/brwjsz/brwjsz"
	"github.com/infinity6-ai/gox/brwjsz/brwqunit"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"github.com/stretchr/testify/require"
)

var defaultImportMap = map[string]string{
	"sql.js": "https://cdn.jsdelivr.net/npm/sql.js@1.14.2/+esm",
}

func findSdkWebDir(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return startDir
		}
		dir = parent
	}
}

func TestWebSdk(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 120*time.Second)
	defer cancel()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	sdkDir := findSdkWebDir(cwd)

	opts := brwjsz.Options{
		Site: &brwjsz.Site{
			Disk: pathz.MustParse(sdkDir),
		},
	}
	opts.Chrome("", true, "")

	b := brwjsz.New(ctx, opts)
	defer b.Close()

	err = b.Start(ctx)
	require.NoError(t, err)

	client := b.Client()

	files, err := filepath.Glob(filepath.Join(sdkDir, "src", "test", "*-test.js"))
	require.NoError(t, err)
	require.NotEmpty(t, files, "no test files found matching src/test/*-test.js")

	sort.Strings(files)

	for _, file := range files {
		fileName := filepath.Base(file)
		testUrl := "/src/test/" + fileName
		t.Run(fileName, func(t *testing.T) {
			brwqunit.TestWithOptions(t, client, brwqunit.Options{
				TestUrls:  []string{testUrl},
				ImportMap: defaultImportMap,
			})
		})
	}
}
