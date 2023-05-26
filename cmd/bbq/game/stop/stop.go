package stop

import (
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdStop = &cobra.Command{
	Use:   "stop",
	Short: "stop game",
	Long:  "stop game. Example: bbq game stop project",
	Run:   stop,
}

// var repoURL string
// var branch string

func init() {
	// if repoURL = os.Getenv("BUTC_LAYOUT_REPO"); repoURL == "" {
	// 	repoURL = "https://git.code..com/gobbq/gobbq-template.git"
	// }
	// CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	// CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
}

func stop(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Help()
		return
	}
	// _stop(com.ServerID(args[0]))
}

// func _stop(sid com.ServerID) {
// 	stopWithSignal(sid, com.StopSignal)
// }

// func stopWithSignal(sid com.ServerID, signal syscall.Signal) {
// 	err := os.Chdir(com.RunEnv.GoBBQRoot)
// 	com.CheckErrorOrQuit(err, "chdir to gobbq directory failed")

// 	ss := status.DetectServerStatus()
// 	status.ShowServerStatus(ss)
// 	if !ss.IsRunning() {
// 		// server is not running
// 		com.ShowMsgAndQuit("no server is running currently")
// 	}

// 	if ss.ServerID != "" && ss.ServerID != sid {
// 		com.ShowMsgAndQuit("another server is running: %s", ss.ServerID)
// 	}

// 	stopGates(ss, signal)
// 	StopGames(ss, signal)
// 	stopDispatcher(ss, signal)
// }

// func StopGames(ss *status.ServerStatus, signal syscall.Signal) {
// 	if ss.NumGamesRunning == 0 {
// 		return
// 	}

// 	com.ShowMsg("stop %d games ...", ss.NumGamesRunning)
// 	for _, proc := range ss.GameProcs {
// 		stopProc(proc, signal)
// 	}
// }

// func stopDispatcher(ss *status.ServerStatus, signal syscall.Signal) {
// 	if ss.NumDispatcherRunning == 0 {
// 		return
// 	}

// 	com.ShowMsg("stop dispatcher ...")
// 	for _, proc := range ss.DispatcherProcs {
// 		stopProc(proc, signal)
// 	}
// }

// func stopGates(ss *status.ServerStatus, signal syscall.Signal) {
// 	if ss.NumGatesRunning == 0 {
// 		return
// 	}

// 	com.ShowMsg("stop %d gates ...", ss.NumGatesRunning)
// 	for _, proc := range ss.GateProcs {
// 		stopProc(proc, signal)
// 	}
// }

// func stopProc(proc process.Process, signal syscall.Signal) {
// 	com.ShowMsg("stop process %s pid=%d", proc.Executable(), proc.Pid())

// 	proc.Signal(signal)
// 	for {
// 		time.Sleep(time.Millisecond * 100)
// 		if !checkProcessRunning(proc) {
// 			break
// 		}
// 	}
// }

// func checkProcessRunning(proc process.Process) bool {
// 	pid := proc.Pid()
// 	procs, err := process.Processes()
// 	com.CheckErrorOrQuit(err, "list processes failed")
// 	for _, _proc := range procs {
// 		if _proc.Pid() == pid {
// 			return true
// 		}
// 	}
// 	return false
// }
