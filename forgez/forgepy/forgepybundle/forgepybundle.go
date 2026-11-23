package forgepybundle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyapi"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyschema"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepysdkfunc"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

type Options struct {
	Apis     []*apispecz.Api
	SdkFuncs []*moldsdkfunc.SdkFunc
}

type pygenerator struct {
	Apis     []*apispecz.Api
	SdkFuncs []*moldsdkfunc.SdkFunc
	Typedefs typeastz.Typedefs
}

func (py *pygenerator) Prepare(ctx context.Context, opts Options) {
	py.Apis = opts.Apis
	py.SdkFuncs = opts.SdkFuncs
	for _, element := range py.Apis {
		logger.Info(ctx, "parsing py api", map[string]any{"spec": element.Spec.Id})
		py.Typedefs.Merge(element.Typedefs)
	}
	py.Typedefs.Merge(apispecz.ExtractTypedefs(py.Apis))

	for _, element := range py.SdkFuncs {
		logger.Info(ctx, "parsing py sdk func", map[string]any{"spec": element.Spec.Id})
		py.Typedefs.Merge(element.Typedefs)
	}
	py.Typedefs.Merge(forgepysdkfunc.ExtractTypedefs(py.SdkFuncs))
}

// Write generates a single consolidated Python module string containing typedefs, APIs, and SDK functions.
func Write(ctx context.Context, opts Options) string {
	var py pygenerator
	py.Prepare(ctx, opts)
	var code strings.Builder
	code.WriteString("\"\"\"Generated SDK APIs, types, and workflows for i6 Python SDK.\"\"\"\n\n")
	code.WriteString("from __future__ import annotations\n\n")
	code.WriteString("from typing import Any, Callable, Dict, Generic, List, Optional, TypeVar, TypedDict, Union\n")
	code.WriteString("from urllib.parse import quote\n\n")
	code.WriteString("T = TypeVar(\"T\")\n\n\n")
	code.WriteString("class Resp(Generic[T]):\n")
	code.WriteString("    \"\"\"HTTP response container with lazy parsing.\"\"\"\n\n")
	code.WriteString("    status: int\n")
	code.WriteString("    headers: Dict[str, str]\n")
	code.WriteString("    ok: bool\n")
	code.WriteString("    text: str\n\n")
	code.WriteString("    def parse(self) -> T:\n")
	code.WriteString("        \"\"\"Parse and return the typed payload.\"\"\"\n")
	code.WriteString("        ...\n\n\n")
	code.WriteString("#\n# typedefs\n#\n\n")
	code.WriteString(forgepyschema.FormatTypedefs(&py.Typedefs))
	code.WriteString("\n#\n# apis\n#\n\n")
	code.WriteString(forgepyapi.Write(ctx, py.Apis))
	code.WriteString("\n#\n# sdkfuncs\n#\n\n")
	code.WriteString(forgepysdkfunc.Write(ctx, py.SdkFuncs))
	return code.String()
}

// WriteFiles generates a map of filename to content for multi-file package output.
func WriteFiles(ctx context.Context, opts Options) map[string]string {
	var py pygenerator
	py.Prepare(ctx, opts)

	files := make(map[string]string)

	// 1. types.py
	files["types.py"] = forgepyschema.WriteModuleTypedefs(ctx, py.Typedefs)

	// 2. apis.py
	var apisCode strings.Builder
	apisCode.WriteString("\"\"\"Generated API endpoint methods for i6 services.\"\"\"\n\n")
	apisCode.WriteString("from __future__ import annotations\n\n")
	apisCode.WriteString("from typing import Any, Dict, List, Optional, Union\n")
	apisCode.WriteString("from urllib.parse import quote\n\n")
	apisCode.WriteString("from ..invoker import Resp, ServiceOptions\n")
	apisCode.WriteString("from .types import *\n\n\n")
	apisCode.WriteString(forgepyapi.Write(ctx, py.Apis))
	apisCode.WriteString("\n\n__all__ = [\"Apis\"]\n")
	files["apis.py"] = apisCode.String()

	// 3. sdkfuncs.py
	var sdkfuncsCode strings.Builder
	sdkfuncsCode.WriteString("\"\"\"Generated SDK workflow functions.\"\"\"\n\n")
	sdkfuncsCode.WriteString("from __future__ import annotations\n\n")
	sdkfuncsCode.WriteString("from typing import Any, Callable, Dict, List, Optional, Union\n\n")
	sdkfuncsCode.WriteString("from ..invoker import Resp, ServiceOptions\n")
	sdkfuncsCode.WriteString("from .types import *\n\n\n")
	sdkfuncsCode.WriteString(forgepysdkfunc.Write(ctx, py.SdkFuncs))
	sdkfuncsCode.WriteString("\n\n__all__ = [\"SdkFuncs\"]\n")
	files["sdkfuncs.py"] = sdkfuncsCode.String()

	// 4. __init__.py
	var initCode strings.Builder
	initCode.WriteString("\"\"\"Generated i6 SDK APIs, schemas, and workflows.\"\"\"\n\n")
	initCode.WriteString("from .apis import Apis\n")
	initCode.WriteString("from .sdkfuncs import SdkFuncs\n")
	initCode.WriteString("from .types import *\n\n")
	initCode.WriteString("__all__ = [\n")
	initCode.WriteString("    \"Apis\",\n")
	initCode.WriteString("    \"SdkFuncs\",\n")
	initCode.WriteString("]\n")
	files["__init__.py"] = initCode.String()

	return files
}

// WriteDir writes the generated Python package files to the target directory.
func WriteDir(ctx context.Context, outDir string, opts Options) error {
	files := WriteFiles(ctx, opts)
	for relPath, content := range files {
		fullPath := filepath.Join(outDir, relPath)
		if err := filez.CreateParentDirs(fullPath); err != nil {
			return fmt.Errorf("failed to create parent dirs for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), os.ModePerm); err != nil {
			return fmt.Errorf("failed to write %s: %w", fullPath, err)
		}
		logger.Info(ctx, "sdkapis file generated", map[string]any{"path": fullPath, "len": len(content)})
	}
	return nil
}
