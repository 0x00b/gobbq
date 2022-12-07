package main

import (
	"log"

	"github.com/0x00b/gobbq/cmd/bbq/game"
	"github.com/0x00b/gobbq/cmd/bbq/proto/project"
	"github.com/0x00b/gobbq/cmd/bbq/proto/proto"
	"github.com/0x00b/gobbq/cmd/bbq/proto/tool"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "bbq",
	Short:   "bbq: An toolkit for Go microservices.",
	Long:    `bbq: An toolkit for Go microservices.`,
	Version: "",
}

func main() {

	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(tool.CmdInstall)
	rootCmd.AddCommand(tool.CmdUpgrade)
	rootCmd.AddCommand(game.CmdGame)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
