package skeleton

import (
	"fmt"
	"os"
	"path"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// MergeMapProperties 合并 map 的配置
func (p *PreLoadBpfSkeleton) MergeMapProperties() (map[string]*ebpf.Map, error) {
	mergedMaps := make(map[string]*ebpf.Map)
	for _, mapMeta := range p.Meta.BpfSkel.Maps {
		if mapMeta.Properties == nil {
			continue
		}

		if mapMeta.Properties.PinPath == "" {
			continue
		}

		mapSpec := p.Spec.Maps[mapMeta.Name]
		if mapSpec == nil {
			return nil, fmt.Errorf("map %s not found", mapMeta.Name)
		}

		newMap, err := ebpf.NewMapWithOptions(mapSpec, ebpf.MapOptions{
			PinPath: mapMeta.Properties.PinPath,
		})
		if err != nil {
			return nil, fmt.Errorf("create map %s error: %w", mapMeta.Name, err)
		}

		mergedMaps[mapMeta.Name] = newMap
	}

	return mergedMaps, nil
}

// 检查 prog、map meta 是否设置了 pinPath，如果设置了，则将程序 pin 到文件系统
func (p *PreLoadBpfSkeleton) CheckPinPath(replaceMaps map[string]*ebpf.Map) error {
	for _, mapMeta := range p.Meta.BpfSkel.Maps {
		if mapMeta.Properties == nil || mapMeta.Properties.PinPath == "" {
			continue
		}

		pinnedMap, err := p.LoadPinMap(mapMeta)
		if err != nil {
			if err == meta.ErrPinnedObjectNotFound {
				continue
			}

			return err
		}

		replaceMaps[mapMeta.Name] = pinnedMap
	}

	return nil
}

func (p *PreLoadBpfSkeleton) LoadPinMap(mapMeta *meta.MapMeta) (*ebpf.Map, error) {
	// Check if the pinned object exists
	if _, err := os.Stat(mapMeta.Properties.PinPath); err != nil {
		return nil, meta.ErrPinnedObjectNotFound
	}

	pinPath := path.Join(mapMeta.Properties.PinPath, mapMeta.Name)
	pinnedMap, err := ebpf.LoadPinnedMap(pinPath, nil)
	if err != nil {
		return nil, fmt.Errorf("load map %s error: %w", mapMeta.Name, err)
	}

	return pinnedMap, nil
}

func (p *PreLoadBpfSkeleton) LoadPinProgram(progMeta *meta.ProgMeta) (*ebpf.Program, error) {
	// Check if the pinned object exists
	if _, err := os.Stat(progMeta.Properties.PinPath); err != nil {
		return nil, meta.ErrPinnedObjectNotFound
	}

	pinnedProg, err := ebpf.LoadPinnedProgram(progMeta.Properties.PinPath, nil)
	if err != nil {
		return nil, fmt.Errorf("load program %s error: %w", progMeta.Name, err)
	}

	delete(p.Spec.Programs, progMeta.Name)

	return pinnedProg, nil
}

// LoadAndAttach 加载并附加 eBPF 程序
func (p *PreLoadBpfSkeleton) LoadAndAttach() (*BpfSkeleton, error) {
	collectionOptions := ebpf.CollectionOptions{}
	mergedMaps, err := p.MergeMapProperties()
	if err != nil {
		return nil, fmt.Errorf("merge map properties error: %w", err)
	}

	if err := p.CheckPinPath(mergedMaps); err != nil {
		return nil, fmt.Errorf("check pin path error: %w", err)
	}

	collectionOptions.MapReplacements = mergedMaps

	// 直接加载 BPF 对象集合，cilium/ebpf 会自动处理 .rodata 和 .bss
	coll, err := ebpf.NewCollectionWithOptions(p.Spec, collectionOptions)
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

		// 如果设置了 pinPath，则将程序 pin 到文件系统
		linkPinPath := progMeta.Properties.LinkPinPath
		if linkPinPath != "" {
			if err := link.Pin(linkPinPath); err != nil {
				return nil, fmt.Errorf("pin program %s error: %w", progMeta.Name, err)
			}
		}

		links = append(links, link)
	}

	return &BpfSkeleton{
		Meta:       p.Meta,
		ConfigData: p.ConfigData,
		Links:      links,
		Collection: coll,
		Btf:        p.Btf,
	}, nil
}
