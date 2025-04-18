// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64

package binary

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type cgroup_skbSpanInfo struct {
	Timestamp    uint64
	StartTime    uint64
	EndTime      uint64
	TraceId      uint64
	SpanId       uint64
	ParentSpanId uint64
	Pid          uint32
	Tid          uint32
	Name         [32]int8
	Pod          [100]int8
	Container    [100]int8
	DevFileId    uint64
}

// loadCgroup_skb returns the embedded CollectionSpec for cgroup_skb.
func loadCgroup_skb() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Cgroup_skbBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load cgroup_skb: %w", err)
	}

	return spec, err
}

// loadCgroup_skbObjects loads cgroup_skb and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*cgroup_skbObjects
//	*cgroup_skbPrograms
//	*cgroup_skbMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadCgroup_skbObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadCgroup_skb()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// cgroup_skbSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type cgroup_skbSpecs struct {
	cgroup_skbProgramSpecs
	cgroup_skbMapSpecs
	cgroup_skbVariableSpecs
}

// cgroup_skbProgramSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type cgroup_skbProgramSpecs struct {
	CountEgressPackets *ebpf.ProgramSpec `ebpf:"count_egress_packets"`
}

// cgroup_skbMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type cgroup_skbMapSpecs struct {
	LinkBegin *ebpf.MapSpec `ebpf:"link_begin"`
	PktCount  *ebpf.MapSpec `ebpf:"pkt_count"`
}

// cgroup_skbVariableSpecs contains global variables before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type cgroup_skbVariableSpecs struct {
	UnusedSpanInfo *ebpf.VariableSpec `ebpf:"unused_span_info"`
}

// cgroup_skbObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadCgroup_skbObjects or ebpf.CollectionSpec.LoadAndAssign.
type cgroup_skbObjects struct {
	cgroup_skbPrograms
	cgroup_skbMaps
	cgroup_skbVariables
}

func (o *cgroup_skbObjects) Close() error {
	return _Cgroup_skbClose(
		&o.cgroup_skbPrograms,
		&o.cgroup_skbMaps,
	)
}

// cgroup_skbMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadCgroup_skbObjects or ebpf.CollectionSpec.LoadAndAssign.
type cgroup_skbMaps struct {
	LinkBegin *ebpf.Map `ebpf:"link_begin"`
	PktCount  *ebpf.Map `ebpf:"pkt_count"`
}

func (m *cgroup_skbMaps) Close() error {
	return _Cgroup_skbClose(
		m.LinkBegin,
		m.PktCount,
	)
}

// cgroup_skbVariables contains all global variables after they have been loaded into the kernel.
//
// It can be passed to loadCgroup_skbObjects or ebpf.CollectionSpec.LoadAndAssign.
type cgroup_skbVariables struct {
	UnusedSpanInfo *ebpf.Variable `ebpf:"unused_span_info"`
}

// cgroup_skbPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadCgroup_skbObjects or ebpf.CollectionSpec.LoadAndAssign.
type cgroup_skbPrograms struct {
	CountEgressPackets *ebpf.Program `ebpf:"count_egress_packets"`
}

func (p *cgroup_skbPrograms) Close() error {
	return _Cgroup_skbClose(
		p.CountEgressPackets,
	)
}

func _Cgroup_skbClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed cgroup_skb_x86_bpfel.o
var _Cgroup_skbBytes []byte
