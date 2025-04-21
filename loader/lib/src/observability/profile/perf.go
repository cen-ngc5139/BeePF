package profile

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"golang.org/x/sys/unix"
)

type PerfEvent struct {
	prog  *ebpf.Program
	links []link.Link
	attrs []unix.PerfEventAttr
	fds   []int
}

func NewPerfEvent(prog *ebpf.Program) (*PerfEvent, error) {
	// 采样率
	simpleRate := 200

	// 创建 perf event 实例
	pe := &PerfEvent{
		prog:  prog,
		links: make([]link.Link, runtime.NumCPU()),
		attrs: make([]unix.PerfEventAttr, runtime.NumCPU()),
		fds:   make([]int, runtime.NumCPU()),
	}

	// 设置采样属性
	attr := unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Size:        uint32(unsafe.Sizeof(unix.PerfEventAttr{})),
		Config:      unix.PERF_COUNT_SW_CPU_CLOCK,
		Sample:      uint64(simpleRate),
		Sample_type: unix.PERF_SAMPLE_RAW,
		Bits:        unix.PerfBitDisabled | unix.PerfBitFreq,
	}

	// 在每个 CPU 上创建 perf event
	for cpu := 0; cpu < runtime.NumCPU(); cpu++ {
		// 打开 perf event
		fd, err := unix.PerfEventOpen(
			&attr,
			-1, // 所有进程
			cpu,
			-1,
			unix.PERF_FLAG_FD_CLOEXEC,
		)
		if err != nil {
			pe.Close()
			return nil, fmt.Errorf("failed to open perf event on CPU %d: %v", cpu, err)
		}
		pe.fds[cpu] = fd

		// 附加 eBPF 程序
		rawLink, err := link.AttachRawLink(link.RawLinkOptions{
			Program: prog,
			Target:  fd,
			Attach:  ebpf.AttachPerfEvent,
		})
		if err != nil {
			pe.Close()
			return nil, fmt.Errorf("failed to attach perf event on CPU %d: %v", cpu, err)
		}
		pe.links[cpu] = rawLink
		pe.attrs[cpu] = attr
	}

	// 启用所有 perf event
	for _, fd := range pe.fds {
		if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_ENABLE, 0); err != nil {
			pe.Close()
			return nil, fmt.Errorf("failed to enable perf event: %v", err)
		}
	}

	return pe, nil
}

// Close 清理资源
func (pe *PerfEvent) Close() error {
	var errs []error

	// 关闭所有 links
	for i, link := range pe.links {
		if link != nil {
			if err := link.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close link %d: %v", i, err))
			}
		}
	}

	// 关闭所有 perf event fd
	for i, fd := range pe.fds {
		if fd != 0 {
			if err := unix.Close(fd); err != nil {
				errs = append(errs, fmt.Errorf("failed to close fd %d: %v", i, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing perf event: %v", errs)
	}
	return nil
}
