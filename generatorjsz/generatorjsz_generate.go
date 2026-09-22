package generatorjsz

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/infinity6-ai/gox/routez/apiz"
)

type generatedFile struct {
	name string
	code string
}

func Generate(outDir string, apis []apiz.Api) ([]string, error) {
	files, err := buildFiles(apis)
	if err != nil {
		return nil, err
	}

	if err := removeStale(outDir, files); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))
	for _, f := range files {
		if err := filez.WriteFile(filepath.Join(outDir, f.name), []byte(f.code)); err != nil {
			return nil, fmt.Errorf("generatorjsz: failed to write %s: %w", f.name, err)
		}
		names = append(names, f.name)
	}
	return names, nil
}

func buildFiles(apis []apiz.Api) ([]generatedFile, error) {
	files := make([]generatedFile, 0, len(apis))

	seenFn := map[string]string{}
	seenFile := map[string]string{}

	for _, api := range apis {
		spec := api.ApiSpec()
		fnName, code, err := BuildJsFetch(spec, api.GetDataRefs())

		if err != nil {
			return nil, fmt.Errorf("generatorjsz: api %q: %w", spec.Id, err)
		}

		name := FileName(spec.Id)

		if other, ok := seenFn[fnName]; ok {
			return nil, fmt.Errorf("generatorjsz: apis %q e %q geram a mesma função %q", other, spec.Id, fnName)
		}

		if other, ok := seenFile[strings.ToLower(name)]; ok {
			return nil, fmt.Errorf("generatorjsz: apis %q e %q geram o mesmo arquivo %q", other, spec.Id, name)
		}

		seenFn[fnName] = spec.Id
		seenFile[strings.ToLower(name)] = spec.Id
		files = append(files, generatedFile{name: name, code: code})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	return files, nil
}

func removeStale(outDir string, keep []generatedFile) error {
	entries, err := os.ReadDir(outDir)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("generatorjsz: failed to list %s: %w", outDir, err)
	}

	wanted := make(map[string]bool, len(keep))
	for _, f := range keep {
		wanted[f.name] = true
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".js") || wanted[e.Name()] {
			continue
		}

		path := filepath.Join(outDir, e.Name())
		content, err := os.ReadFile(path)

		if err != nil {
			return fmt.Errorf("generatorjsz: failed to read %s: %w", path, err)
		}

		if !strings.HasPrefix(string(content), generatedHeaderPrefix) {
			continue
		}

		if err := os.Remove(path); err != nil {
			return fmt.Errorf("generatorjsz: failed to remove stale %s: %w", path, err)
		}
	}

	return nil
}
