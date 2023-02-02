package com

import (
	"fmt"
	"os"
)

func ShowMsgAndQuit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "! "+format+"\n", a...)
	os.Exit(2)
}

func ShowMsg(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "> "+format+"\n", a...)
}

func CheckErrorOrQuit(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "! %s: %v\n", msg, err)
		os.Exit(2)
	}
}
