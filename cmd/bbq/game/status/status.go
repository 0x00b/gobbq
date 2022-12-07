package status

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/0x00b/gobbq/cmd/bbq/game/com"
	"github.com/0x00b/gobbq/cmd/bbq/game/process"
	"github.com/0x00b/gobbq/engine/config"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "status",
	Short: "show game status",
	Long:  "show game status. Example: bbq game status",
	Run:   run,
}

func init() {
}

func run(cmd *cobra.Command, args []string) {
	Status()
}

// ServerStatus represents the status of a server
type ServerStatus struct {
	NumDispatcherRunning int
	NumGatesRunning      int
	NumGamesRunning      int

	DispatcherProcs []process.Process
	GateProcs       []process.Process
	GameProcs       []process.Process
	ServerID        com.ServerID
}

// IsRunning returns if a server is running
func (ss *ServerStatus) IsRunning() bool {
	return ss.NumDispatcherRunning > 0 || ss.NumGatesRunning > 0 || ss.NumGamesRunning > 0
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}
	return true
}

func DetectServerStatus() *ServerStatus {
	ss := &ServerStatus{}
	procs, err := process.Processes()
	com.CheckErrorOrQuit(err, "list processes failed")
	for _, proc := range procs {
		path, err := proc.Path()
		if err != nil {
			continue
		}

		if !IsExists(path) {
			cmdline, err := proc.CmdlineSlice()
			if err != nil {
				continue
			}
			path = cmdline[0]
			if !filepath.IsAbs(path) {
				cwd, err := proc.Cwd()
				if err != nil {
					continue
				}
				path = filepath.Join(cwd, path)
			}

		}

		relpath, err := filepath.Rel(com.RunEnv.GoWorldRoot, path)
		if err != nil || strings.HasPrefix(relpath, "..") {
			continue
		}

		dir, file := filepath.Split(relpath)

		if file == "dispatcher"+com.BinaryExtension {
			ss.NumDispatcherRunning++
			ss.DispatcherProcs = append(ss.DispatcherProcs, proc)
		} else if file == "gate"+com.BinaryExtension {
			ss.NumGatesRunning++
			ss.GateProcs = append(ss.GateProcs, proc)
		} else {
			if strings.HasSuffix(dir, string(filepath.Separator)) {
				dir = dir[:len(dir)-1]
			}
			serverid := com.ServerID(strings.Join(strings.Split(dir, string(filepath.Separator)), "/"))
			if strings.HasPrefix(string(serverid), "cmd/") || strings.HasPrefix(string(serverid), "components/") || string(serverid) == "examples/test_client" {
				// this is a cmd or a component, not a game
				continue
			}
			ss.NumGamesRunning++
			ss.GameProcs = append(ss.GameProcs, proc)
			if ss.ServerID == "" {
				ss.ServerID = serverid
			}
		}
	}

	return ss
}

func Status() {
	ss := DetectServerStatus()
	ShowServerStatus(ss)
}

func ShowServerStatus(ss *ServerStatus) {
	com.ShowMsg("%d dispatcher running, %d/%d gates running, %d/%d games (%s) running", ss.NumDispatcherRunning,
		ss.NumGatesRunning, config.GetDeployment().DesiredGates,
		ss.NumGamesRunning, config.GetDeployment().DesiredGames,
		ss.ServerID,
	)

	var listProcs []process.Process
	listProcs = append(listProcs, ss.DispatcherProcs...)
	listProcs = append(listProcs, ss.GameProcs...)
	listProcs = append(listProcs, ss.GateProcs...)
	for _, proc := range listProcs {
		cmdlineSlice, err := proc.CmdlineSlice()
		var cmdline string
		if err == nil {
			cmdline = strings.Join(cmdlineSlice, " ")
		} else {
			cmdline = fmt.Sprintf("get cmdline failed: %e", err)
		}

		com.ShowMsg("\t%-10d%-16s%s", proc.Pid(), proc.Executable(), cmdline)
	}
}
