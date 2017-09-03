package main

import (
	"runtime"

	"github.com/cryptounicorns/gluttony/cli"
)

func init() { runtime.GOMAXPROCS(runtime.NumCPU()) }
func main() { cli.Execute() }
