package cmdserver

import (
	"github.com/infinity6-ai/i6sdk/sdkweb/internal/sdkwebserver"
	"github.com/spf13/cobra"
)

func prepareServerStart(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			s := sdkwebserver.Start(sdkwebserver.Options{
				Context: cmd.Context(),
			})
			s.Serve()
		},
	}
	parent.AddCommand(cmd)
}

func Prepare(parent *cobra.Command) {
	cmd := &cobra.Command{Use: "server"}
	parent.AddCommand(cmd)

	prepareServerStart(cmd)

}
