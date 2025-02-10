package skeleton

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// LoadAndAttach 加载并附加 eBPF 程序
func (p *PreLoadBpfSkeleton) LoadAndAttach() (*BpfSkeleton, error) {
	// 直接加载 BPF 对象集合，cilium/ebpf 会自动处理 .rodata 和 .bss
	coll, err := ebpf.NewCollection(p.Spec)
	if err != nil {
		return nil, fmt.Errorf("load collection error: %w", err)
	}

	// 附加程序
	var links []link.Link
	for _, progMeta := range p.Meta.BpfSkel.Progs {
		prog := coll.Programs[progMeta.Name]
		if prog == nil {
			return nil, fmt.Errorf("program %s not found", progMeta.Name)
		}

		progSpec := p.Spec.Programs[progMeta.Name]
		if progSpec == nil {
			return nil, fmt.Errorf("program %s not found", progMeta.Name)
		}

		if !progMeta.Link {
			continue // 跳过不需要 link 的程序
		}

		// 根据不同的 AttachType 使用对应的 attach 方式

		link, err := progMeta.AttachProgram(progSpec, prog)
		if err != nil {
			return nil, fmt.Errorf("attach program %s error: %w", progMeta.Name, err)
		}
		links = append(links, link)
	}

	return &BpfSkeleton{
		Meta:       p.Meta,
		ConfigData: p.ConfigData,
		Links:      links,
		Collection: coll,
	}, nil
}
