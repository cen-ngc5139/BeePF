package container

import (
	"bytes"
	"fmt"

	"github.com/cilium/ebpf/btf"
)

// BTFContainer 用于解析和存储 eBPF 程序中的 BTF 信息
type BTFContainer struct {
	// spec 存储从 eBPF ELF 文件中解析出的 BTF 信息
	spec *btf.Spec
	// elfContainer 保存原始 ELF 文件数据
	elfContainer *ElfContainer
}

// NewBTFContainerFromBinary 从 eBPF ELF 二进制数据创建 BTF 容器
func NewBTFContainerFromBinary(elfData []byte) (*BTFContainer, error) {
	// 创建 ELF 容器
	elfContainer, err := NewElfContainerFromBinary(elfData)
	if err != nil {
		return nil, fmt.Errorf("create elf container error: %w", err)
	}

	// 从 ELF 文件中解析 BTF 信息
	spec, err := btf.LoadSpecFromReader(bytes.NewReader(elfData))
	if err != nil {
		return nil, fmt.Errorf("parse BTF from ELF error: %w", err)
	}

	return &BTFContainer{
		spec:         spec,
		elfContainer: elfContainer,
	}, nil
}

// GetSpec 获取 BTF 规范
func (b *BTFContainer) GetSpec() *btf.Spec {
	return b.spec
}

// GetElfContainer 获取 ELF 容器
func (b *BTFContainer) GetElfContainer() *ElfContainer {
	return b.elfContainer
}
