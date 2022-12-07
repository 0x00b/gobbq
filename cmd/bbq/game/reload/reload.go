package reload

import (
	"os"

	"github.com/0x00b/gobbq/cmd/bbq/game/com"
	"github.com/0x00b/gobbq/cmd/bbq/game/start"
	"github.com/0x00b/gobbq/cmd/bbq/game/status"
	"github.com/0x00b/gobbq/cmd/bbq/game/stop"
	"github.com/0x00b/gobbq/engine/binutil"
	"github.com/0x00b/gobbq/engine/config"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "reload",
	Short: "reload game",
	Long:  "reload game. Example: bbq game reload xxx",
	Run:   run,
}

func init() {
}

func run(cmd *cobra.Command, args []string) {

}
func reload(sid com.ServerID, daemonMode bool) {
	err := os.Chdir(com.RunEnv.GoBBQRoot)
	com.CheckErrorOrQuit(err, "chdir to gobbq directory failed")

	ss := status.DetectServerStatus()
	status.ShowServerStatus(ss)
	if !ss.IsRunning() {
		// server is not running
		com.ShowMsgAndQuit("no server is running currently")
	}

	if ss.ServerID != "" && ss.ServerID != sid {
		com.ShowMsgAndQuit("another server is running: %s", ss.ServerID)
	}

	if ss.NumGamesRunning == 0 {
		com.ShowMsgAndQuit("no game is running")
	} else if ss.NumGamesRunning != config.GetDeployment().DesiredGames {
		com.ShowMsgAndQuit("found %d games, but should have %d", ss.NumGamesRunning, config.GetDeployment().DesiredGames)
	}

	stop.StopGames(ss, binutil.FreezeSignal)
	start.StartGames(sid, true, daemonMode)
}
