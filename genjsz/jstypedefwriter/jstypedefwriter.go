package jstypedefwriter

import "go.code.infinity6.ai/sdkgo/genjsz/jstypez"

func WritePipedefs(pipedefs jstypez.Typedefs) {
	for _, pipedef := range pipedefs.All() {
		// genarete jsdocs pipedef
		_ = pipedef
	}
}
