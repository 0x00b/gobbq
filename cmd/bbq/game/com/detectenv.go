package com

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/0x00b/gobbq/engine/config"
)

// Env represents environment variables
type Env struct {
	GoBBQRoot string
}

// GetDispatcherDir returns the path to the dispatcher
func (env *Env) GetDispatcherDir() string {
	return filepath.Join(env.GoBBQRoot, "components", "dispatcher")
}

// GetGateDir returns the path to the gate
func (env *Env) GetGateDir() string {
	return filepath.Join(env.GoBBQRoot, "components", "gate")
}

// GetDispatcherBinary returns the path to the dispatcher binary
func (env *Env) GetDispatcherBinary() string {
	return filepath.Join(env.GetDispatcherDir(), "dispatcher"+BinaryExtension)
}

// GetGateBinary returns the path to the gate binary
func (env *Env) GetGateBinary() string {
	return filepath.Join(env.GetGateDir(), "gate"+BinaryExtension)
}

var RunEnv Env

func getGoSearchPaths() []string {
	var paths []string
	goroot := os.Getenv("GOROOT")
	if goroot != "" {
		paths = append(paths, goroot)
	}

	gopath := os.Getenv("GOPATH")
	for _, p := range strings.Split(gopath, string(os.PathListSeparator)) {
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

type ModuleInfo struct {
	Path      string `json:"Path"`
	Main      bool   `json:"Main"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
}

func goListModule() (*ModuleInfo, error) {
	cmd := exec.Command("go", "list", "-m", "-json")

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(r)
	var mi ModuleInfo
	err = d.Decode(&mi)
	if err != nil {
		return nil, err
	}

	cmd.Wait()
	return &mi, err
}

func _DetectGoBBQPath() string {
	mi, err := goListModule()
	if err == nil {
		ShowMsg("go list -m -json: %+v", *mi)
		return mi.Dir
	}

	searchPaths := getGoSearchPaths()
	ShowMsg("go search paths: %s", strings.Join(searchPaths, string(os.PathListSeparator)))
	for _, sp := range searchPaths {
		gobbqPath := filepath.Join(sp, "src", "github.com", "0x00b", "gobbq")
		if IsDir(gobbqPath) {
			return gobbqPath
		}
	}
	return ""
}

func DetectGoBBQPath() {
	RunEnv.GoBBQRoot = _DetectGoBBQPath()
	if RunEnv.GoBBQRoot == "" {
		ShowMsgAndQuit("gobbq directory is not detected")
	}

	ShowMsg("gobbq directory found: %s", RunEnv.GoBBQRoot)
	configFile := filepath.Join(RunEnv.GoBBQRoot, "gobbq.ini")
	config.SetConfigFile(configFile)
}
