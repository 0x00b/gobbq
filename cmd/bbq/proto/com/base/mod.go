package base

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"
)

// ModName 获取当前项目中的go mod name
// 简单实现，当前文件夹找不到go.mod，则在父目录中找
//
//	@return string
func ModName() string {
	fileName := "go.mod"
	var modBytes []byte
	loopCnt := 0
	for {
		loopCnt++
		fileName = "../" + fileName
		var err error
		if modBytes, err = ioutil.ReadFile(fileName); err == nil {
			break
		}
		if loopCnt > 10 {
			fmt.Println("not found go.mod")
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

// ModulePath returns go module path.
func ModulePath(filename string) (string, error) {
	modBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}

// ModuleVersion returns module version.
func ModuleVersion(path string) (string, error) {
	stdout := &bytes.Buffer{}
	fd := exec.Command("go", "mod", "graph")
	fd.Stdout = stdout
	fd.Stderr = stdout
	if err := fd.Run(); err != nil {
		return "", err
	}
	rd := bufio.NewReader(stdout)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			return "", err
		}
		str := string(line)
		i := strings.Index(str, "@")
		if strings.Contains(str, path+"@") && i != -1 {
			return path + str[i:], nil
		}
	}
}

// // ButcMod returns bbq mod.
// func ButcMod() string {
// 	// go 1.15+ read from env GOMODCACHE
// 	cacheOut, _ := exec.Command("go", "env", "GOMODCACHE").Output()
// 	cachePath := strings.Trim(string(cacheOut), "\n")
// 	pathOut, _ := exec.Command("go", "env", "GOPATH").Output()
// 	gopath := strings.Trim(string(pathOut), "\n")
// 	if cachePath == "" {
// 		cachePath = filepath.Join(gopath, "pkg", "mod")
// 	}
// 	if path, err := ModuleVersion("github.com/go-butc/butc/v2"); err == nil {
// 		// $GOPATH/pkg/mod/github.com/go-butc/butc@v2
// 		return filepath.Join(cachePath, path)
// 	}
// 	// $GOPATH/src/github.com/go-butc/butc
// 	return filepath.Join(gopath, "src", "github.com", "go-butc", "butc")
// }
