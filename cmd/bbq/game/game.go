package game

import (
	"github.com/0x00b/gobbq/cmd/bbq/game/build"
	"github.com/0x00b/gobbq/cmd/bbq/game/reload"
	"github.com/0x00b/gobbq/cmd/bbq/game/start"
	"github.com/0x00b/gobbq/cmd/bbq/game/status"
	"github.com/0x00b/gobbq/cmd/bbq/game/stop"
	"github.com/spf13/cobra"
)

// CmdGame represents the proto command.
var CmdGame = &cobra.Command{
	Use:   "game",
	Short: "game control",
	Long:  "game control.",
	Run:   run,
}

func init() {
	CmdGame.AddCommand(build.CmdNew)
	CmdGame.AddCommand(reload.CmdNew)
	CmdGame.AddCommand(start.CmdNew)
	CmdGame.AddCommand(status.CmdNew)
	CmdGame.AddCommand(stop.CmdKill)
	CmdGame.AddCommand(stop.CmdStop)
}

func run(cmd *cobra.Command, args []string) {
	// com.DetectGoBBQPath()
}
