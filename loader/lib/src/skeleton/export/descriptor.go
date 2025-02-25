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

	// 处理结构体类型
	if st, ok := b.Type.(*btf.Struct); ok {
		return b.buildStructMembers(st)
	}

	// 处理非结构体类型（如 map 的 key 或 value）
	size, err := btf.Sizeof(b.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get size of type %s: %w", b.Name, err)
	}

	// 对于非结构体类型，创建单个成员表示整个类型
	result := []CheckedExportedMember{
		{
			FieldName:          b.Name,
			Type:               b.Type,
			BitOffset:          0,
			Size:               btf.Bits(size * 8), // 转换为比特
			OutputHeaderOffset: 0,
		},
	}

	return result, nil
}

// buildStructMembers 处理结构体类型的成员
func (b *BTFTypeDescriptor) buildStructMembers(st *btf.Struct) ([]CheckedExportedMember, error) {
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
			Size:               btf.Bits(size * 8), // 转换为比特
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
