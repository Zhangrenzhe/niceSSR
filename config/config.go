package config

import (
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
)

var RootPath,_ = filepath.Abs(".")
var ShadowPath string

func init()  {
	if runtime.GOOS == "windows" {
		ShadowPath = RootPath + "\\atools\\shadowsocks2-win64.exe"
	}else if runtime.GOOS == "linux" {
		ShadowPath = RootPath + "/atools/shadowsocks2-linux-amd64"
		cmd := exec.Command("chmod","+x", ShadowPath)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}else {
		ShadowPath = RootPath + "/atools/shadowsocks2-macos-amd64"
		cmd := exec.Command("chmod","+x", ShadowPath)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}
