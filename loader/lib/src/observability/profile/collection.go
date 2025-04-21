package profile

import (
	"fmt"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/observability/add2line"
)

func BuildAddr2Line(vmlinux string) (*add2line.Addr2Line, error) {
	kallsyms, err := add2line.NewKallsyms()
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/kallsyms: %v", err)
	}

	if vmlinux == "" {
		vmlinux, err = add2line.FindVmlinux()
		if err != nil {
			return nil, fmt.Errorf("failed to find vmlinux: %v", err)
		}
	}

	textAddr, err := add2line.ReadTextAddrFromVmlinux(vmlinux)
	if err != nil {
		return nil, fmt.Errorf("failed to read .text address from vmlinux: %v", err)
	}

	kaslrOffset := textAddr - kallsyms.Stext()
	addr2line, err := add2line.NewAddr2Line(vmlinux, kaslrOffset, kallsyms.SysBPF())
	if err != nil {
		return nil, fmt.Errorf("failed to create addr2line: %v", err)
	}

	return addr2line, nil
}
