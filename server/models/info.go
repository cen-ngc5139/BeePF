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
	LoadTime         time.Time

	Maps []ebpf.MapID
}

type ProgramDetail struct {
	ProgramInfoWrapper
	MapsDetail []MapInfoWrapper
}

type MapInfoWrapper struct {
	// Type of the map.
	Type ebpf.MapType
	// KeySize is the size of the map key in bytes.
	KeySize uint32
	// ValueSize is the size of the map value in bytes.
	ValueSize uint32
	// MaxEntries is the maximum number of entries the map can hold. Its meaning
	// is map-specific.
	MaxEntries uint32
	// Flags used during map creation.
	Flags uint32
	// Name as supplied by user space at load time. Available from 4.15.
	Name string

	ID       ebpf.MapID
	BTF      btf.ID
	MapExtra uint64
	Memlock  uint64
	Frozen   bool
}
