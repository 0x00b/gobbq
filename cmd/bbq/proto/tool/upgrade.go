package tool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/0x00b/gobbq/cmd/bbq/proto/com/base"
	"github.com/spf13/cobra"
)

// CmdUpgrade represents the new command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade tools",
	Long:  "upgrade tools. Example: bbq install",
	Run:   upgrade,
}

func upgrade(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"github.com/0x00b/gobbq/cmd/protoc-gen-gobbq",
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("bbq install")
	installCmd := exec.Command("bbq", "install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		panic(err)
	}
}
