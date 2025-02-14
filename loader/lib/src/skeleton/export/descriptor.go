package export

import (
	"fmt"

	"github.com/cilium/ebpf/btf"
)

// TypeDescriptor 描述导出数据的类型信息
type TypeDescriptor interface {
	// BuildCheckedExportedMembers 构建已检查的导出成员列表
	BuildCheckedExportedMembers() ([]CheckedExportedMember, error)
}

// BTFTypeDescriptor BTF类型描述符
type BTFTypeDescriptor struct {
	Type btf.Type
	Name string
}

func (b *BTFTypeDescriptor) BuildCheckedExportedMembers() ([]CheckedExportedMember, error) {
	// 检查类型名称是否匹配
	if b.Type.TypeName() != b.Name {
		return nil, fmt.Errorf("type name mismatch: %s != %s", b.Type.TypeName(), b.Name)
	}

	// 获取结构体信息
	st, ok := b.Type.(*btf.Struct)
	if !ok {
		return nil, fmt.Errorf("type %s is not struct", b.Type.TypeName())
	}

	// 构建导出成员
	var result []CheckedExportedMember
	for _, member := range st.Members {
		// 检查位域
		if member.BitfieldSize%8 != 0 {
			return nil, fmt.Errorf("bitfield is not supported for member %s", member.Name)
		}

		size, err := btf.Sizeof(member.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to get size of member %s: %w", member.Name, err)
		}

		result = append(result, CheckedExportedMember{
			FieldName:          member.Name,
			Type:               member.Type,
			BitOffset:          btf.Bits(member.Offset),
			Size:               btf.Bits(size),
			OutputHeaderOffset: 0,
		})
	}

	return result, nil
}

// NewBTFTypeDescriptor 创建BTF类型描述符
func NewBTFTypeDescriptor(typ btf.Type, name string) TypeDescriptor {
	return &BTFTypeDescriptor{
		Type: typ,
		Name: name,
	}
}
