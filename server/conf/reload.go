package conf

import (
	"fmt"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func reloadConfig() {
	ParseConfig(ConfigFile, true)
	fmt.Println("Reload config complete")
}

func Reload() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		reloadConfig()
		<-ticker.C
	}
}
