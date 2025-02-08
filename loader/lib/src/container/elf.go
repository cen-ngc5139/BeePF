package container

import (
	"bytes"
	"debug/elf"
	"fmt"
	"io"
)

// ElfContainer ELF容器
// 用于解决 ELF 文件引用问题的辅助结构体
// 包含原始 ELF 文件的二进制数据和 ELF 文件结构
type ElfContainer struct {
	// bin 保存原始二进制数据
	bin []byte
	// file 保存解析后的 ELF 文件
	file *elf.File
	// reader 提供二进制数据的读取接口
	reader *bytes.Reader
}

// NewElfContainerFromBinary 从二进制数据创建 ELF 容器
func NewElfContainerFromBinary(bin []byte) (*ElfContainer, error) {
	// 创建二进制数据读取器
	reader := bytes.NewReader(bin)

	// 解析 ELF 文件
	file, err := elf.NewFile(reader)
	if err != nil {
		return nil, fmt.Errorf("parse elf error: %w", err)
	}

	return &ElfContainer{
		bin:    bin,
		file:   file,
		reader: bytes.NewReader(bin),
	}, nil
}

// File 获取 ELF 文件结构
func (e *ElfContainer) File() *elf.File {
	return e.file
}

// Reader 获取二进制数据读取器
func (e *ElfContainer) Reader() io.ReaderAt {
	return e.reader
}

// Binary 获取原始二进制数据
func (e *ElfContainer) Binary() []byte {
	return e.bin
}

// Close 关闭并清理资源
func (e *ElfContainer) Close() error {
	if e.file != nil {
		return e.file.Close()
	}
	return nil
}

// Reset 重置读取器位置
func (e *ElfContainer) Reset() {
	e.reader.Reset(e.bin)
}

// GetElfInfo 获取 ELF 文件的基本信息
func (e *ElfContainer) GetElfInfo() map[string]interface{} {
	return map[string]interface{}{
		"类型":    e.file.Type.String(),
		"机器架构":  e.file.Machine.String(),
		"入口点地址": fmt.Sprintf("0x%x", e.file.Entry),
		"字节序":   e.file.ByteOrder.String(),
	}
}

// GetSections 获取所有段信息
func (e *ElfContainer) GetSections() []map[string]interface{} {
	sections := make([]map[string]interface{}, 0)
	for _, section := range e.file.Sections {
		sections = append(sections, map[string]interface{}{
			"名称":  section.Name,
			"类型":  section.Type.String(),
			"大小":  section.Size,
			"偏移量": section.Offset,
			"权限":  section.Flags.String(),
		})
	}
	return sections
}

// GetSymbols 获取所有符号信息
func (e *ElfContainer) GetSymbols() ([]map[string]interface{}, error) {
	symbols, err := e.file.Symbols()
	if err != nil {
		return nil, fmt.Errorf("获取符号表失败: %w", err)
	}

	symbolInfo := make([]map[string]interface{}, 0)
	for _, sym := range symbols {
		symbolInfo = append(symbolInfo, map[string]interface{}{
			"名称":  sym.Name,
			"类型":  sym.Info,
			"地址":  sym.Value,
			"大小":  sym.Size,
			"段索引": sym.Section,
		})
	}
	return symbolInfo, nil
}
