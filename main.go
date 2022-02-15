package main

import (
	"JDC-Monitor/cmd"
	"github.com/mizuki1412/go-core-kit/init/initkit"
)

func main() {
	initkit.LoadConfig()
	cmd.Execute()
}
