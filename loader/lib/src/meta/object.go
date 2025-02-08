package meta

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

// ComposedObjectInner 内部复合对象
// 用于 JSON 序列化和反序列化的内部结构
type ComposedObjectInner struct {
	// BpfObject base64编码的、zlib压缩的对象文件
	BpfObject string `json:"bpf_object"`

	// BpfObjectSize 未压缩对象文件的大小（字节）
	BpfObjectSize uint `json:"bpf_object_size"`

	// Meta 元数据对象
	Meta EunomiaObjectMeta `json:"meta"`
}

// ComposedObject 完整的 eunomia JSON 对象
// 包含 eBPF 对象文件和元数据信息
// 原始 JSON 格式示例：
//
//	{
//	   "bpf_object": "", // base64编码、zlib压缩的对象文件
//	   "bpf_object_size": 0, // 未压缩的对象文件大小（字节）
//	   "meta": {} // 元数据对象
//	}
type ComposedObject struct {
	// BpfObject 对象二进制数据
	BpfObject []byte

	// Meta 元数据信息
	Meta EunomiaObjectMeta
}

// MarshalJSON 实现 json.Marshaler 接口
// 将 ComposedObject 序列化为 JSON
func (c *ComposedObject) MarshalJSON() ([]byte, error) {
	// 压缩对象文件
	compressed, err := CompressZlib(c.BpfObject)
	if err != nil {
		return nil, fmt.Errorf("compress object error: %w", err)
	}

	// 转换为 base64
	base64Data := base64.StdEncoding.EncodeToString(compressed)

	// 创建内部对象
	inner := ComposedObjectInner{
		BpfObject:     base64Data,
		BpfObjectSize: uint(len(c.BpfObject)),
		Meta:          c.Meta,
	}

	return json.Marshal(inner)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
// 从 JSON 解析 ComposedObject
func (c *ComposedObject) UnmarshalJSON(data []byte) error {
	var inner ComposedObjectInner
	if err := json.Unmarshal(data, &inner); err != nil {
		return fmt.Errorf("unmarshal json error: %w", err)
	}

	// 解码 base64
	compressed, err := base64.StdEncoding.DecodeString(inner.BpfObject)
	if err != nil {
		return fmt.Errorf("decode base64 error: %w", err)
	}

	// 解压数据
	decompressed, err := DecompressZlib(compressed)
	if err != nil {
		return fmt.Errorf("decompress error: %w", err)
	}

	// 验证大小
	if uint(len(decompressed)) != inner.BpfObjectSize {
		return fmt.Errorf("size mismatch: got %d, want %d", len(decompressed), inner.BpfObjectSize)
	}

	c.BpfObject = decompressed
	c.Meta = inner.Meta
	return nil
}

// CompressZlib 使用 zlib 压缩数据
func CompressZlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecompressZlib 使用 zlib 解压数据
func DecompressZlib(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
