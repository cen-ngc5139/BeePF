package topology

import (
	"fmt"
	"strings"

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

	// 创建一个字符串构建器
	var sb strings.Builder

	// 输出程序头信息
	sb.WriteString(fmt.Sprintf("程序ID: %d\n", progID))
	sb.WriteString(fmt.Sprintf("程序名称: %s\n", prog.Name))
	sb.WriteString(fmt.Sprintf("程序类型: %s\n", prog.Type))
	sb.WriteString(fmt.Sprintf("程序标签: %s\n", prog.Tag))
	sb.WriteString("\n指令:\n")

	for i, insn := range insns {
		// 输出指令编号和内容
		sb.WriteString(fmt.Sprintf("%4d: %v", i, insn))

		// 如果有源代码信息，也输出源代码
		if src := insn.Source(); src != nil {
			sb.WriteString(fmt.Sprintf("\n# %s", src))
		}

		sb.WriteString("\n")
	}

	return []byte(sb.String()), nil
}
