package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/observability/topology"
)

func main() {

	dump, err := topology.GetProgDumpJited(18585)
	if err != nil {
		log.Fatalf("Failed to get prog dump: %v", err)
	}

	fmt.Println(string(dump))

	collector, err := topology.NewNodeMetricsCollector(time.Second, nil)
	if err != nil {
		log.Fatalf("fail to create metrics collector: %v", err)
	}

	if err := collector.Start(); err != nil {
		log.Fatalf("fail to start metrics collector: %v", err)
	}

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	err = collector.Stop()
	if err != nil {
		log.Fatalf("fail to stop metrics collection: %v", err)
	}

	log.Println("正常关闭")

}
