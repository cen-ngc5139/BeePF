package btf

import (
	"fmt"
	"github.com/cilium/ebpf/btf"
)

// LoadSystemBTF 加载系统 BTF 信息
func LoadSystemBTF(path string) (spec *btf.Spec, err error) {
	// 直接使用 cilium/ebpf 的 API 加载 BTF
	if path != "" {
		// 从指定路径加载
		spec, err = btf.LoadSpec(path)
	} else {
		// 直接从内核加载
		spec, err = btf.LoadKernelSpec()
	}

	if err != nil {
		return nil, fmt.Errorf("load BTF error: %w", err)
	}

	return spec, nil
}
