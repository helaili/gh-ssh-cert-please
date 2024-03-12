package main

import (
	"log"

	"github.com/helaili/gh-ssh-cert-please/cmd"
)

var (
	version = "next"
	commit  = ""
)

func main() {
	rootCmd := cmd.New(version, commit)
	if err := rootCmd.Execute(); err != nil {
		l := log.New(rootCmd.ErrOrStderr(), "", 0)
		l.Fatalln("🚫", err.Error())
	}
}
