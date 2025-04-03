package topology

import (
	"bytes"
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/knightsc/gapstone"
)

func GetProgDumpXlated(progID ebpf.ProgramID) ([]byte, error) {
	prog, err := GetProgInfo(progID)
	if err != nil {
		return nil, err
	}

	insns, err := prog.Instructions()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// 写入程序头部信息
	fmt.Fprintf(&buf, "程序ID: %d\n", progID)
	fmt.Fprintf(&buf, "程序名称: %s\n", prog.Name)
	fmt.Fprintf(&buf, "程序类型: %s\n", prog.Type)
	fmt.Fprintf(&buf, "程序标签: %s\n\n", prog.Tag)

	// 写入指令标题
	fmt.Fprintf(&buf, "指令列表:\n")
	fmt.Fprintf(&buf, "---------------------------------------\n")

	// 遍历并格式化所有指令
	for i, insn := range insns {
		// 生成地址
		address := i * 8

		// 获取指令的原始表示
		insnStr := fmt.Sprintf("%v", insn)

		// 获取源代码信息
		srcInfo := ""
		if src := insn.Source(); src != nil {
			srcInfo = fmt.Sprintf("# %s", src)
		}

		// 使用更加可读的格式输出
		fmt.Fprintf(&buf, "%04x: %-40s %s\n", address, insnStr, srcInfo)
	}

	return buf.Bytes(), nil
}

func GetProgDumpJited(progID ebpf.ProgramID) ([]byte, error) {
	// 根据 ID 获取程序信息
	prog, err := ebpf.NewProgramFromID(progID)
	if err != nil {
		return nil, fmt.Errorf("failed to get program from id: %w", err)
	}

	defer prog.Close()

	engine, err := gapstone.New(int(gapstone.CS_ARCH_X86), int(gapstone.CS_MODE_64))
	if err != nil {
		return nil, fmt.Errorf("failed to create gapstone engine: %w", err)
	}
	defer engine.Close()

	bpfProgInfo, err := NewBPFProgInfo(prog, engine)
	if err != nil {
		return nil, fmt.Errorf("failed to create bpf prog info: %w", err)
	}

	progInfo, err := prog.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get program info: %w", err)
	}

	fullName, err := GetFullName(progInfo)
	if err == nil {
		progInfo.Name = fullName
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "程序ID: %d\n", progID)
	fmt.Fprintf(&buf, "程序名称: %s\n", progInfo.Name)
	fmt.Fprintf(&buf, "程序类型: %s\n", progInfo.Type)
	fmt.Fprintf(&buf, "程序标签: %s\n\n", progInfo.Tag)

	for _, prog := range bpfProgInfo.Progs {
		fmt.Fprintf(&buf, "---------------------------------------\n")
		for _, lineInfo := range prog.LineInfos {
			// fmt.Fprintf(&buf, "文件:%s ", lineInfo.Line.FileName())
			// fmt.Fprintf(&buf, "行号:%d ", lineInfo.Line.LineNumber())
			fmt.Fprintf(&buf, "%s\n", lineInfo.Line.String())
		}
	}

	return buf.Bytes(), nil
}
