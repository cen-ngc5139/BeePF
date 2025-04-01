package models

import (
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
)

type ProgramInfoWrapper struct {
	Type ebpf.ProgramType
	ID   ebpf.ProgramID
	// Truncated hash of the BPF bytecode. Available from 4.13.
	Tag string
	// Name as supplied by user space at load time. Available from 4.15.
	Name string

	CreatedByUID     uint32
	HaveCreatedByUID bool
	BTF              btf.ID
	LoadTime         time.Duration

	Maps []ebpf.MapID
}

type MapInfoWrapper struct {
	ID   ebpf.MapID
	Name string
}
