package utils

const _bingoDaemonFlag = "__BINGOCHILDPROCESS"

// Daemon 进程改为守护方式执行.
func Daemon(nochdir, noclose int) error {
	// already a daemon
	// if len(os.Args) > 0 && os.Args[len(os.Args)-1] == _bingoDaemonFlag {
	// 	/* Change the file mode mask */
	// 	syscall.Umask(0)

	// 	if nochdir == 0 {
	// 		if err := os.Chdir("/"); err != nil {
	// 			xlog.Error("daemon chdir err:%s", err.Error())
	// 		}
	// 	}

	// 	return nil
	// }

	// curFileNum := 3
	// maxFileNum := 6
	// files := make([]*os.File, curFileNum, maxFileNum)
	// if noclose == 0 {
	// 	nullDev, err := os.OpenFile("/dev/null", 0, 0)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	files[0], files[1], files[2] = nullDev, nullDev, nullDev
	// } else {
	// 	files[0], files[1], files[2] = os.Stdin, os.Stdout, os.Stderr
	// }

	// dir, err := os.Getwd()
	// if err != nil {
	// 	xlog.Error("daemon Getwd err:%s", err.Error())
	// }
	// sysattrs := syscall.SysProcAttr{Setsid: true}
	// attrs := os.ProcAttr{Dir: dir, Env: os.Environ(), Files: files, Sys: &sysattrs}

	// os.Args = append(os.Args, _bingoDaemonFlag)
	// proc, err := os.StartProcess(os.Args[0], os.Args, &attrs)
	// if err != nil {
	// 	return fmt.Errorf("can't create process=%s err:%w", os.Args[0], err)
	// }
	// err = proc.Release()
	// if err != nil {
	// 	xlog.Error("daemon proc Release err:%s", err.Error())
	// }
	// os.Exit(-1)
	return nil
}
