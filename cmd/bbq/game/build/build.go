package build

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/0x00b/gobbq/cmd/bbq/game/com"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "build",
	Short: "build project",
	Long:  "build project. Example: bbq game build project",
	Run:   run,
}

func init() {
}

func run(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Help()
		return
	}
	build(com.ServerID(args[0]))
}

func build(sid com.ServerID) {
	com.ShowMsg("building server %s ...", sid)

	buildServer(sid)
	buildDispatcher()
	buildGate()
}

func buildServer(sid com.ServerID) {
	serverPath := sid.Path()
	com.ShowMsg("server directory is %s ...", serverPath)
	if !com.IsDir(serverPath) {
		com.ShowMsgAndQuit("wrong server id: %s, using '\\' instead of '/'?", sid)
	}

	com.ShowMsg("go build %s ...", sid)
	buildDirectory(serverPath)
}

func buildDispatcher() {
	com.ShowMsg("go build dispatcher ...")
	buildDirectory(filepath.Join(com.RunEnv.GoBBQRoot, "components", "dispatcher"))
}

func buildGate() {
	com.ShowMsg("go build gate ...")
	buildDirectory(filepath.Join(com.RunEnv.GoBBQRoot, "components", "gate"))
}

func buildDirectory(dir string) {
	var err error
	var curdir string
	curdir, err = os.Getwd()
	com.CheckErrorOrQuit(err, "")

	err = os.Chdir(dir)
	com.CheckErrorOrQuit(err, "")

	defer os.Chdir(curdir)

	cmd := exec.Command("go", "build", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	com.CheckErrorOrQuit(err, "")
	return
}
