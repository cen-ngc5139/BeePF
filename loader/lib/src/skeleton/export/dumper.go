package export

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/cilium/ebpf/btf"
)

// DumpToJson 将 BTF 类型数据转换为 JSON
func DumpToJson(typ btf.Type, data []byte) (json.RawMessage, error) {
	switch t := typ.(type) {
	case *btf.Int:
		return dumpInt(t, data)
	case *btf.Pointer:
		if _, ok := t.Target.(*btf.Struct); ok {
			return DumpToJson(t.Target, data)
		}
		return dumpPointer(data)
	case *btf.Array:
		return dumpArray(t, data)
	case *btf.Struct:
		return dumpStruct(t, data)
	case *btf.Enum:
		return dumpEnum(t, data)
	case *btf.Float:
		return dumpFloat(t, data)
	case *btf.Typedef:
		return handleTypedef(t, data)
	case *btf.Volatile:
		return DumpToJson(typ, data)
	case *btf.Const:
		return DumpToJson(typ, data)
	default:
		return nil, fmt.Errorf("unsupported type: %T", t)
	}
}

// DumpToJsonWithCheckedTypes 使用已检查的类型信息将数据转换为 JSON
func DumpToJsonWithCheckedTypes(
	checkedTypes []CheckedExportedMember,
	data []byte,
) (json.RawMessage, error) {
	// 创建结果 map
	result := make(map[string]interface{})

	// 处理每个成员
	for _, member := range checkedTypes {
		// 从 Type 中获取实际的 Size
		size, err := btf.Sizeof(member.Type)
		if err != nil {
			return nil, fmt.Errorf("get size error: %w", err)
		}

		offset := uint64(member.BitOffset) / 8
		if member.BitOffset%8 != 0 {
			return nil, fmt.Errorf("bit offset must be byte-aligned: %s", member.FieldName)
		}

		// 确保数据长度足够
		end := offset + uint64(size)
		if uint64(len(data)) < end {
			return nil, fmt.Errorf(
				"input buffer too small for field %s: need %d..%d bytes, got %d bytes",
				member.FieldName,
				offset,
				end,
				len(data),
			)
		}

		// 获取字段数据
		fieldData := data[offset:end]

		// 转换字段数据为 JSON
		fieldJson, err := DumpToJson(member.Type, fieldData)
		if err != nil {
			return nil, fmt.Errorf("failed to dump field %s: %w", member.FieldName, err)
		}

		// 由于 json.Unmarshal 会丢失精度，所以使用 json.NewDecoder 来解码
		var fieldValue interface{}
		decoder := json.NewDecoder(bytes.NewReader(fieldJson))
		decoder.UseNumber()
		if err := decoder.Decode(&fieldValue); err != nil {
			return nil, fmt.Errorf("failed to decode field %s JSON: %w", member.FieldName, err)
		}

		// 添加到结果 map
		result[member.FieldName] = fieldValue
	}

	// 编码最终结果
	return json.Marshal(result)
}

// dumpInt 处理整数类型
func dumpInt(t *btf.Int, data []byte) (json.RawMessage, error) {
	if t.Encoding == btf.Bool {
		return json.Marshal(data[0] != 0)
	}

	size := t.Size
	if len(data) < int(size) {
		return nil, fmt.Errorf("data too short for int: need %d, got %d", size, len(data))
	}

	var val interface{}
	switch size {
	case 1:
		if t.Encoding == btf.Signed {
			val = int8(data[0])
		} else {
			val = uint8(data[0])
		}
	case 2:
		if t.Encoding == btf.Signed {
			val = int16(binary.LittleEndian.Uint16(data))
		} else {
			val = binary.LittleEndian.Uint16(data)
		}
	case 4:
		if t.Encoding == btf.Signed {
			val = int32(binary.LittleEndian.Uint32(data))
		} else {
			val = binary.LittleEndian.Uint32(data)
		}
	case 8:
		if t.Encoding == btf.Signed {
			val = int64(binary.LittleEndian.Uint64(data))
		} else {
			val = binary.LittleEndian.Uint64(data)
		}
	default:
		return nil, fmt.Errorf("unsupported int size: %d", size)
	}

	return json.Marshal(val)
}

// dumpPointer 处理指针类型
func dumpPointer(data []byte) (json.RawMessage, error) {
	switch len(data) {
	case 4:
		return json.Marshal(binary.LittleEndian.Uint32(data))
	case 8:
		return json.Marshal(binary.LittleEndian.Uint64(data))
	default:
		return nil, fmt.Errorf("invalid pointer size: %d", len(data))
	}
}

// dumpArray 处理数组类型
func dumpArray(t *btf.Array, data []byte) (json.RawMessage, error) {
	elemType := t.Type
	// 处理字符串数组
	if _, ok := elemType.(*btf.Int); ok && elemType.TypeName() == "char" {
		// 查找字符串结束位置
		strLen := 0
		for strLen < len(data) && data[strLen] != 0 {
			strLen++
		}
		return json.Marshal(string(data[:strLen]))
	}

	// 处理普通数组
	elemSize, err := btf.Sizeof(elemType)
	if err != nil {
		return nil, fmt.Errorf("get element size error: %w", err)
	}

	result := make([]json.RawMessage, t.Nelems)
	elemSizeU32 := uint32(elemSize)

	for i := uint32(0); i < t.Nelems; i++ {
		start := i * elemSizeU32
		end := start + elemSizeU32
		if end > uint32(len(data)) {
			return nil, fmt.Errorf("array data too short")
		}

		elem, err := DumpToJson(elemType, data[start:end])
		if err != nil {
			return nil, fmt.Errorf("dump array element %d error: %w", i, err)
		}
		result[i] = elem
	}

	return json.Marshal(result)
}

// dumpStruct 处理结构体类型
func dumpStruct(t *btf.Struct, data []byte) (json.RawMessage, error) {
	result := make(map[string]interface{})
	result["__EUNOMIA_TYPE"] = "struct"
	result["__EUNOMIA_TYPE_NAME"] = t.Name

	for _, member := range t.Members {
		offset := uint32(member.Offset / 8)
		if member.Offset%8 != 0 {
			return nil, fmt.Errorf("bit fields not supported: %s", member.Name)
		}

		memberType := member.Type
		size, err := btf.Sizeof(memberType)
		if err != nil {
			return nil, fmt.Errorf("get member size error: %w", err)
		}

		if offset+uint32(size) > uint32(len(data)) {
			return nil, fmt.Errorf("data too short for member %s", member.Name)
		}

		memberValue, err := DumpToJson(memberType, data[offset:offset+uint32(size)])
		if err != nil {
			return nil, fmt.Errorf("dump member %s error: %w", member.Name, err)
		}

		result[member.Name] = memberValue
	}

	return json.Marshal(result)
}

// dumpEnum 处理枚举类型
func dumpEnum(t *btf.Enum, data []byte) (json.RawMessage, error) {
	size, err := btf.Sizeof(t)
	if err != nil {
		return nil, fmt.Errorf("get enum size error: %w", err)
	}

	var val int64
	switch size {
	case 1:
		val = int64(int8(data[0]))
	case 2:
		val = int64(binary.LittleEndian.Uint16(data))
	case 4:
		val = int64(binary.LittleEndian.Uint32(data))
	default:
		return nil, fmt.Errorf("unsupported enum size: %d", size)
	}

	// 查找枚举值
	for _, v := range t.Values {
		if int64(v.Value) == val {
			return json.Marshal(fmt.Sprintf("%s(%d)", v.Name, val))
		}
	}

	return json.Marshal(fmt.Sprintf("<UNKNOWN_VARIANT>(%d)", val))
}

// dumpFloat 处理浮点数类型
func dumpFloat(t *btf.Float, data []byte) (json.RawMessage, error) {
	size, err := btf.Sizeof(t)
	if err != nil {
		return nil, fmt.Errorf("get enum size error: %w", err)
	}
	switch size {
	case 4:
		bits := binary.LittleEndian.Uint32(data)
		val := math.Float32frombits(bits)
		return json.Marshal(val)
	case 8:
		bits := binary.LittleEndian.Uint64(data)
		val := math.Float64frombits(bits)
		return json.Marshal(val)
	default:
		return nil, fmt.Errorf("unsupported float size: %d", size)
	}
}

// DumpToString 将 BTF 类型数据转换为字符串
func DumpToString(typ btf.Type, data []byte) (string, error) {
	// 先转换为 JSON
	jsonVal, err := DumpToJson(typ, data)
	if err != nil {
		return "", fmt.Errorf("convert to json error: %w", err)
	}

	// 使用 json.decode 解码
	decoder := json.NewDecoder(bytes.NewReader(jsonVal))
	decoder.UseNumber()
	var val interface{}
	if err := decoder.Decode(&val); err != nil {
		return "", fmt.Errorf("decode json error: %w", err)
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%v", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// DumpToStringWithCheckedTypes 使用已检查的类型信息将数据转换为字符串
func DumpToStringWithCheckedTypes(
	checkedTypes []CheckedExportedMember,
	data []byte,
	out *strings.Builder,
) error {
	for _, member := range checkedTypes {
		// 处理输出头部偏移
		if uint64(out.Len()) < uint64(member.OutputHeaderOffset) {
			out.WriteString(strings.Repeat(" ",
				int(member.OutputHeaderOffset)-out.Len()))
		} else {
			out.WriteString(" ")
		}

		// 计算字段偏移量
		offset := member.BitOffset / 8
		if member.BitOffset%8 != 0 {
			return fmt.Errorf(
				"bit field found in member %s, but it is not supported now",
				member.FieldName,
			)
		}

		size, err := btf.Sizeof(member.Type)
		if err != nil {
			return fmt.Errorf("get size error: %w", err)
		}

		// 确保数据长度足够
		end := uint64(offset) + uint64(size)
		if uint64(len(data)) < end {
			return fmt.Errorf(
				"data too short for member %s: need %d bytes, got %d",
				member.FieldName, end, len(data),
			)
		}

		// 获取字段数据并转换为字符串
		fieldData := data[offset:end]
		str, err := DumpToString(member.Type, fieldData)
		if err != nil {
			return fmt.Errorf("dump member %s error: %w", member.FieldName, err)
		}

		// 写入结果
		out.WriteString(str)
	}

	return nil
}

// 处理 typedef 类型
func handleTypedef(typedef *btf.Typedef, data []byte) (json.RawMessage, error) {
	// 特别处理 __u32 类型
	if typedef.Name == "__u32" {
		if len(data) < 4 {
			return nil, fmt.Errorf("data too short for __u32: need 4 bytes, got %d", len(data))
		}
		// 使用小端序读取 uint32 值
		value := binary.LittleEndian.Uint32(data[:4])
		return json.Marshal(value)
	}
	// 对于其他 typedef，递归处理其底层类型
	return DumpToJson(typedef.Type, data)
}
