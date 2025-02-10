package meta

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"golang.org/x/sys/unix"
)

// GenerateMeta 生成元数据
func GenerateMeta(objectFile string) (*EunomiaObjectMeta, error) {
	// 打开 ELF 文件
	spec, err := ebpf.LoadCollectionSpec(objectFile)
	if err != nil {
		return nil, fmt.Errorf("加载 ELF 文件失败: %v", err)
	}

	// 解析数据段信息
	dataSections := make([]DataSectionMeta, 0)

	exportTypes := make([]ExportedTypesStructMeta, 0)
	// 处理全局变量
	for name, varSpec := range spec.Variables {
		section := DataSectionMeta{
			Name:      name,
			Variables: make([]DataSectionVariableMeta, 0),
		}

		// 获取类型信息
		btfType := varSpec.Type().Type
		typeName := getActualTypeName(btfType)

		// 如果是指针类型，获取目标类型
		if ptr, isPtr := btfType.(*btf.Pointer); isPtr {
			if target, isStruct := ptr.Target.(*btf.Struct); isStruct {
				// 获取结构体信息
				structName := target.Name
				structSize := target.Size

				// 构建成员信息
				members := make([]ExportedTypesStructMemberMeta, 0)
				for _, member := range target.Members {
					memberMeta := ExportedTypesStructMemberMeta{
						Name: member.Name,
						Type: getActualTypeName(member.Type),
					}
					members = append(members, memberMeta)
				}

				// 添加到导出类型
				exportType := ExportedTypesStructMeta{
					Name:    structName,
					Size:    uint32(structSize),
					Members: members,
				}
				exportTypes = append(exportTypes, exportType)

				typeName = structName // 更新类型名为结构体名
			}
		}

		// 添加变量
		varMeta := DataSectionVariableMeta{
			Name: name,
			Type: typeName,
		}
		section.Variables = append(section.Variables, varMeta)
		dataSections = append(dataSections, section)
	}

	// 创建元数据结构
	meta := EunomiaObjectMeta{
		ExportTypes: exportTypes,
		BpfSkel: BpfSkeletonMeta{
			Maps:         convertMaps(spec.Maps),
			Progs:        convertProgs(spec.Programs),
			DataSections: dataSections,
		},
		PerfBufferPages:  64,  // 默认值
		PerfBufferTimeMs: 10,  // 默认值
		PollTimeoutMs:    100, // 默认值
	}

	return &meta, nil
}

// 获取类型的实际名称
func getActualTypeName(t btf.Type) string {
	// 如果类型为空，返回 "unknown"
	if t == nil {
		return "unknown"
	}

	// 获取当前类型名
	typeName := t.TypeName()

	// 如果类型名不为空，直接返回
	if typeName != "" {
		return typeName
	}

	// 根据不同类型继续查找
	switch v := t.(type) {
	case *btf.Array:
		return fmt.Sprintf("%s[%d]", getActualTypeName(v.Type), v.Nelems)
	case *btf.Pointer:
		return fmt.Sprintf("*%s", getActualTypeName(v.Target))
	case *btf.Int:
		// 对于整型，返回具体的类型名（如 char, int 等）
		return v.Name
	// 可以根据需要添加其他类型的处理
	default:
		return "unknown"
	}
}

// 转换 Maps
func convertMaps(maps map[string]*ebpf.MapSpec) map[string]*MapMeta {
	result := make(map[string]*MapMeta)
	for name, mapSpec := range maps {
		meta := &MapMeta{
			Name:   name,
			Ident:  name,
			Mmaped: mapSpec.Flags&unix.BPF_F_MMAPABLE != 0,
		}
		result[name] = meta
	}
	return result
}

// 转换 Programs
func convertProgs(progs map[string]*ebpf.ProgramSpec) map[string]*ProgMeta {
	result := make(map[string]*ProgMeta)
	for name, progSpec := range progs {
		meta := &ProgMeta{
			Name:   name,
			Attach: progSpec.Type.String(), // 使用程序类型作为附加点
			Link:   true,                   // 默认生成 bpf_link
		}
		result[name] = meta
	}
	return result
}
