package topology

import (
	"bytes"
	"fmt"

	"github.com/cilium/ebpf"
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
		// 输出指令编号
		fmt.Fprintf(&buf, "%4d: ", i)

		// 使用内置的 Format 方法格式化指令
		var insnBuf bytes.Buffer
		fmt.Fprintf(&insnBuf, "%v", insn)
		fmt.Fprintf(&buf, "%-40s", insnBuf.String())

		// 如果有源代码信息，也输出源代码
		if src := insn.Source(); src != nil {
			fmt.Fprintf(&buf, "  # %s", src)
		}

		fmt.Fprintf(&buf, "\n")
	}

	return buf.Bytes(), nil
}
