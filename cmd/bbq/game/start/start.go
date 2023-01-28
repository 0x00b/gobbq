package start

import (
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "start",
	Short: "start game",
	Long:  "start game. Example: bbq game start project [-d true]",
	Run:   run,
}

// var repoURL string
var daemonMode bool

func init() {
	// if repoURL = os.Getenv("BUTC_LAYOUT_REPO"); repoURL == "" {
	// 	repoURL = "https://git.code.tencent.com/gobbq/gobbq-template.git"
	// }
	// CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().BoolVarP(&daemonMode, "daemonMode", "d", false, "daemon mode")
}

func run(cmd *cobra.Command, args []string) {
	// start(com.ServerID(args[0]), daemonMode)
}

// func start(sid com.ServerID, daemonMode bool) {
// 	err := os.Chdir(com.RunEnv.GoBBQRoot)
// 	com.CheckErrorOrQuit(err, "chdir to gobbq directory failed")

// 	ss := status.DetectServerStatus()
// 	if ss.NumDispatcherRunning > 0 || ss.NumGatesRunning > 0 {
// 		status.Status()
// 		com.ShowMsgAndQuit("server is already running, can not start multiple servers")
// 	}

// 	startDispatchers(daemonMode)
// 	StartGames(sid, false, daemonMode)
// 	startGates(daemonMode)
// }

// func startDispatchers(daemonMode bool) {
// 	com.ShowMsg("start dispatchers ...")
// 	dispatcherIds := config.GetDispatcherIDs()
// 	com.ShowMsg("dispatcher ids: %v", dispatcherIds)
// 	for _, dispid := range dispatcherIds {
// 		startDispatcher(dispid, daemonMode)
// 	}
// }

// func startDispatcher(dispid uint16, daemonMode bool) {
// 	cfg := config.GetDispatcher(dispid)
// 	args := []string{"-dispid", strconv.Itoa(int(dispid))}
// 	if daemonMode {
// 		args = append(args, "-d")
// 	}
// 	cmd := exec.Command(com.RunEnv.GetDispatcherBinary(), args...)
// 	err := runCmdUntilTag(cmd, cfg.LogFile, consts.DISPATCHER_STARTED_TAG, time.Second*10)
// 	com.CheckErrorOrQuit(err, "start dispatcher failed, see dispatcher.log for error")
// }

// func StartGames(sid com.ServerID, isRestore bool, daemonMode bool) {
// 	com.ShowMsg("start games ...")
// 	desiredGames := config.GetDeployment().DesiredGames
// 	com.ShowMsg("desired games = %d", desiredGames)
// 	for gameid := uint16(1); int(gameid) <= desiredGames; gameid++ {
// 		startGame(sid, gameid, isRestore, daemonMode)
// 	}
// }

// func startGame(sid com.ServerID, gameid uint16, isRestore bool, daemonMode bool) {
// 	com.ShowMsg("start game %d ...", gameid)

// 	gameExePath := filepath.Join(sid.Path(), sid.Name()+com.BinaryExtension)
// 	args := []string{"-gid", strconv.Itoa(int(gameid))}
// 	if isRestore {
// 		args = append(args, "-restore")
// 	}
// 	if daemonMode {
// 		args = append(args, "-d")
// 	}
// 	cmd := exec.Command(gameExePath, args...)
// 	err := runCmdUntilTag(cmd, config.GetGame(gameid).LogFile, consts.GAME_STARTED_TAG, time.Second*600)
// 	com.CheckErrorOrQuit(err, "start game failed, see game.log for error")
// }

// func startGates(daemonMode bool) {
// 	com.ShowMsg("start gates ...")
// 	desiredGates := config.GetDeployment().DesiredGates
// 	com.ShowMsg("desired gates = %d", desiredGates)
// 	for gateid := uint16(1); int(gateid) <= desiredGates; gateid++ {
// 		startGate(gateid, daemonMode)
// 	}
// }

// func startGate(gateid uint16, daemonMode bool) {
// 	com.ShowMsg("start gate %d ...", gateid)

// 	args := []string{"-gid", strconv.Itoa(int(gateid))}
// 	if daemonMode {
// 		args = append(args, "-d")
// 	}
// 	cmd := exec.Command(com.RunEnv.GetGateBinary(), args...)
// 	err := runCmdUntilTag(cmd, config.GetGate(gateid).LogFile, consts.GATE_STARTED_TAG, time.Second*10)
// 	com.CheckErrorOrQuit(err, "start gate failed, see gate.log for error")
// }

// func runCmdUntilTag(cmd *exec.Cmd, logFile string, tag string, timeout time.Duration) (err error) {
// 	clearLogFile(logFile)
// 	err = cmd.Start()
// 	if err != nil {
// 		return
// 	}

// 	timeoutTime := time.Now().Add(timeout)
// 	for time.Now().Before(timeoutTime) {
// 		time.Sleep(time.Millisecond * 200)
// 		if isTagInFile(logFile, tag) {
// 			cmd.Process.Release()
// 			return
// 		}
// 	}

// 	err = errors.Errorf("wait started tag timeout")
// 	return
// }

// func clearLogFile(logFile string) {
// 	ioutil.WriteFile(logFile, []byte{}, 0644)
// }

// func isTagInFile(filename string, tag string) bool {
// 	data, err := ioutil.ReadFile(filename)
// 	com.CheckErrorOrQuit(err, "read file error")
// 	return strings.Contains(string(data), tag)
// }
