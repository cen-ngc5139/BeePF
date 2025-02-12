package export

import (
	"fmt"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf/btf"
)

// CheckExportTypesBtf 检查导出类型与 BTF 信息的匹配性
func CheckExportTypesBtf(structMeta meta.ExportedTypesStructMeta, spec *btf.Spec) ([]CheckedExportedMember, error) {
	// 获取实际的结构体类型
	st, err := getActualStructType(structMeta.Type)
	if err != nil {
		return nil, err
	}

	// 验证类型名称
	if st.Name != structMeta.Name {
		return nil, fmt.Errorf(
			"type names don't match: `%s` from btf, but `%s` from struct_meta",
			st.Name, structMeta.Name,
		)
	}

	// 4. 验证成员数量
	if len(st.Members) != len(structMeta.Members) {
		return nil, fmt.Errorf(
			"unmatched member count: `%d` from btf, but `%d` from struct_meta",
			len(st.Members), len(structMeta.Members),
		)
	}

	// 5. 检查并构建导出成员
	var result []CheckedExportedMember
	for i, btfMem := range st.Members {
		metaMem := structMeta.Members[i]

		// 验证成员名称
		if btfMem.Name != metaMem.Name {
			continue
		}

		// 验证位域
		if btfMem.BitfieldSize > 0 {
			return nil, fmt.Errorf(
				"bitfield is not supported. Member %s, bit_offset=%d, bit_sz=%d",
				btfMem.Name, btfMem.Offset, btfMem.BitfieldSize,
			)
		}

		// 构建检查后的成员
		result = append(result, CheckedExportedMember{
			FieldName:          metaMem.Name,
			Type:               btfMem.Type,
			BitOffset:          btfMem.Offset, // 转换为 bit 偏移
			Size:               btfMem.BitfieldSize,
			OutputHeaderOffset: 0,
		})
	}

	return result, nil
}

// getActualStructType 递归解析类型直到找到实际的结构体类型
func getActualStructType(t btf.Type) (*btf.Struct, error) {
	for {
		switch v := t.(type) {
		case *btf.Struct:
			return v, nil
		case *btf.Pointer:
			t = v.Target
		case *btf.Var:
			t = v.Type
		case *btf.Typedef:
			t = v.Type
		case *btf.Volatile:
			t = v.Type
		case *btf.Const:
			t = v.Type
		default:
			return nil, fmt.Errorf("unexpected type %T, expected struct", t)
		}
	}
}
