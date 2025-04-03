package main

import (
	"fmt"
	"log"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/observability/topology"
)

func main() {

	dump, err := topology.GetProgDumpJited(18585)
	if err != nil {
		log.Fatalf("Failed to get prog dump: %v", err)
	}

	fmt.Println(string(dump))
}
