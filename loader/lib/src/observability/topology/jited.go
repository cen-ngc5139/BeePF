// Copyright 2024 Leon Hwang.
// SPDX-License-Identifier: Apache-2.0
package topology

import (
	"fmt"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/knightsc/gapstone"
)

type BpfProgInfo struct {
	Progs []*BpfProgAddrLineInfo

	IsLbrProg bool
}

type BpfProgKaddrRange struct {
	Start, End uintptr
}

type BpfProgAddrLineInfo struct {
	KaddrRange BpfProgKaddrRange
	FuncName   string

	JitedLineInfo []uintptr        // ordered
	LineInfos     []btf.LineOffset // mapping 1:1 with jitedLineInfo
}

type BpfProgLineInfo struct {
	FuncName string
	KsymAddr uintptr

	FileName string
	FileLine uint32
}

type BpfProgs struct {
	Progs map[string]*ebpf.Program
	Funcs map[uintptr]*BpfProgInfo // func IP -> prog info
}

func NewBPFProgInfo(prog *ebpf.Program, engine gapstone.Engine) (*BpfProgInfo, error) {
	// 从 ebpf prog 中获取信息，包括：类型、名称、附加点等
	pinfo, err := prog.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get prog info: %w", err)
	}

	// 从 ebpf prog 中获取函数信息，包括：函数名、函数地址等
	funcInfos, err := pinfo.FuncInfos()
	if err != nil {
		return nil, fmt.Errorf("failed to get func infos: %w", err)
	}

	// 获取 ebpf prog 每行指令信息
	lines, err := pinfo.LineInfos()
	if err != nil {
		return nil, fmt.Errorf("failed to get line infos: %w", err)
	}

	// 获取 ebpf prog 每行指令的 jited 指令
	jitedInsns, _ := pinfo.JitedInsns()
	// 获取 ebpf prog 每行指令的 jited ksyms
	jitedKsyms, _ := pinfo.JitedKsymAddrs()
	// 获取 ebpf prog 每行指令的 jited 函数长度
	jitedFuncLens, _ := pinfo.JitedFuncLens()
	// 获取 ebpf prog 每行指令的 jited line infos
	jitedLineInfos, _ := pinfo.JitedLineInfos()

	if len(funcInfos) != len(jitedFuncLens) {
		return nil, fmt.Errorf("func info number %d != jited func lens number %d", len(funcInfos), len(jitedFuncLens))
	}

	if len(jitedKsyms) != len(jitedFuncLens) {
		return nil, fmt.Errorf("jited ksyms number %d != jited func lens number %d", len(jitedKsyms), len(jitedFuncLens))
	}

	if len(jitedLineInfos) != len(lines) {
		return nil, fmt.Errorf("line info number %d != jited line info number %d", len(lines), len(jitedLineInfos))
	}

	jited2li := make(map[uint64]btf.LineOffset, len(jitedLineInfos))
	for i, kaddr := range jitedLineInfos {
		jited2li[kaddr] = lines[i]
	}

	var progInfo BpfProgInfo
	progInfo.Progs = make([]*BpfProgAddrLineInfo, 0, len(jitedFuncLens))

	insns := jitedInsns
	for i, funcLen := range jitedFuncLens {
		ksym := uint64(jitedKsyms[i])
		fnInsns := insns[:funcLen]
		pc := uint64(0)

		var info BpfProgAddrLineInfo
		info.KaddrRange.Start = jitedKsyms[i]
		info.KaddrRange.End = info.KaddrRange.Start + uintptr(funcLen)
		info.FuncName = strings.TrimSpace(funcInfos[i].Func.Name)

		for len(fnInsns) > 0 {
			kaddr := ksym + pc
			if li, ok := jited2li[kaddr]; ok {
				info.JitedLineInfo = append(info.JitedLineInfo, uintptr(kaddr))
				info.LineInfos = append(info.LineInfos, li)
			}

			inst, err := engine.Disasm(fnInsns, kaddr, 1)
			if err != nil {
				return nil, fmt.Errorf("failed to disasm instruction: %w", err)
			}

			instSize := uint64(inst[0].Size)
			pc += instSize
			fnInsns = fnInsns[instSize:]
		}

		progInfo.Progs = append(progInfo.Progs, &info)

		insns = insns[funcLen:]
	}

	return &progInfo, nil
}
