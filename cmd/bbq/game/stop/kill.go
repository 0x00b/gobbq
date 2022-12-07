package stop

import (
	"github.com/0x00b/gobbq/cmd/bbq/game/com"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdKill = &cobra.Command{
	Use:   "kill",
	Short: "kill game",
	Long:  "kill game. Example: bbq game kill xxxx",
	Run:   kill,
}

func init() {
}

func kill(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Help()
		return
	}
	_kill(com.ServerID(args[0]))
}

func _kill(sid com.ServerID) {
	// stopWithSignal(sid, syscall.SIGKILL)
}
