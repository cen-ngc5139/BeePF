package meta

import (
	"fmt"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// AttachProgram 根据程序类型选择合适的 attach 方式
func (p *ProgMeta) AttachProgram(spec *ebpf.ProgramSpec, program *ebpf.Program) (link.Link, error) {
	switch spec.Type {
	case ebpf.UnspecifiedProgram:
		return nil, fmt.Errorf("error:%v, %s", ErrSectionFormat, "invalid program type, make sure to use the right section prefix")
	case ebpf.Kprobe:
		return p.attachKprobe(program)
	case ebpf.TracePoint:
		return p.attachTracepoint(program)
	case ebpf.CGroupDevice, ebpf.CGroupSKB, ebpf.CGroupSock, ebpf.SockOps, ebpf.CGroupSockAddr, ebpf.CGroupSockopt, ebpf.CGroupSysctl:
		return p.attachCGroup()
	case ebpf.SocketFilter:
		return p.attachSocket()
	case ebpf.SchedCLS:
		return p.attachTCCLS()
	case ebpf.XDP:
		return p.attachXDP()
	case ebpf.RawTracepoint:
		return p.attachRawTracepoint(program)
	case ebpf.Tracing:
		return p.attachTracing(program)
	case ebpf.LSM:
		return p.attachLsm()
	default:
		return nil, fmt.Errorf("program type %s not implemented yet", spec.Type)
	}
}

func (p *ProgMeta) attachTracing(program *ebpf.Program) (link.Link, error) {
	tracing, err := link.AttachTracing(link.TracingOptions{Program: program})
	if err != nil {
		return nil, fmt.Errorf("error:%v, couldn's activate tracing %s, matchFuncName:%s", err, p.Attach, p.Name)
	}
	return tracing, nil
}

func (p *ProgMeta) attachKprobe(program *ebpf.Program) (link.Link, error) {
	// Prepare kprobe_events line parameters
	var err error
	funcName := p.Name
	isRet := false
	if strings.HasPrefix(p.Attach, "kretprobe/") {
		isRet = true
	} else if strings.HasPrefix(p.Attach, "kprobe/") {
		isRet = false
	} else {
		// this might actually be a Uprobe
		return p.attachUprobe(program)
	}

	var kp link.Link
	if isRet {
		kp, err = link.Kretprobe(funcName, program, nil)
	} else {
		kp, err = link.Kprobe(funcName, program, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("opening Kprobe: %s, funcName:%s, isRet:%t, section:%s", err, funcName, isRet, p.Attach)
	}

	return kp, nil
}

func (p *ProgMeta) attachUprobe(program *ebpf.Program) (link.Link, error) {
	return nil, nil
}

func (p *ProgMeta) attachTracepoint(program *ebpf.Program) (link.Link, error) {
	// Parse section
	traceGroup := strings.SplitN(p.Attach, "/", 3)
	if len(traceGroup) != 3 {
		return nil, fmt.Errorf("error:%v, expected SEC(\"tracepoint/[category]/[name]\") got %s", ErrSectionFormat, p.Attach)
	}
	category := traceGroup[1]
	name := traceGroup[2]

	kp, err := link.Tracepoint(category, name, program, nil)
	if err != nil {
		return nil, fmt.Errorf("error:%v , couldn's activate tracepoint %s, matchFuncName:%s", err, p.Attach, p.Name)
	}

	return kp, nil
}

func (p *ProgMeta) attachCGroup() (link.Link, error) {
	return nil, nil
}

func (p *ProgMeta) attachSocket() (link.Link, error) {
	return nil, nil
}

func (p *ProgMeta) attachTCCLS() (link.Link, error) {
	return nil, nil
}

func (p *ProgMeta) attachXDP() (link.Link, error) {
	return nil, nil
}

func (p *ProgMeta) attachRawTracepoint(program *ebpf.Program) (link.Link, error) {
	name := strings.TrimLeft(p.Attach, "raw_tracepoint/")
	link, err := link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    name,
		Program: program,
	})
	if err != nil {
		return nil, fmt.Errorf("error:%v , couldn's activate raw_tracepoint %s, matchFuncName:%s", err, p.Attach, p.Name)
	}

	return link, nil
}

func (p *ProgMeta) attachLsm() (link.Link, error) {
	return nil, nil
}
