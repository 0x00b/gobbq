package proto

import (
	"github.com/0x00b/gobbq/cmd/bbq/proto/proto/add"
	"github.com/0x00b/gobbq/cmd/bbq/proto/proto/gen"
	"github.com/spf13/cobra"
)

// CmdProto represents the proto command.
var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files.",
	Run:   run,
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(gen.CmdGen)
}

func run(cmd *cobra.Command, args []string) {

}
