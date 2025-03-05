package models

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf"
	"github.com/pkg/errors"
)

type Component struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Programs []Program `json:"programs"`
	Maps     []Map     `json:"maps"`
}

type Program struct {
	Id          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Spec        ProgramSpec            `json:"spec"`
	Properties  meta.ProgramProperties `json:"properties"`
}

type ProgramSpec struct {
	// Name is passed to the kernel as a debug aid. Must only contain
	// alpha numeric and '_' characters.
	Name string

	// Type determines at which hook in the kernel a program will run.
	Type ebpf.ProgramType

	// AttachType of the program, needed to differentiate allowed context
	// accesses in some newer program types like CGroupSockAddr.
	//
	// Available on kernels 4.17 and later.
	AttachType ebpf.AttachType

	// Name of a kernel data structure or function to attach to. Its
	// interpretation depends on Type and AttachType.
	AttachTo string

	// The name of the ELF section this program originated from.
	SectionName string

	// Flags is passed to the kernel and specifies additional program
	// load attributes.
	Flags uint32

	// License of the program. Some helpers are only available if
	// the license is deemed compatible with the GPL.
	//
	// See https://www.kernel.org/doc/html/latest/process/license-rules.html#id1
	License string

	// Version used by Kprobe programs.
	//
	// Deprecated on kernels 5.0 and later. Leave empty to let the library
	// detect this value automatically.
	KernelVersion uint32
}

type Map struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Spec        MapSpec            `json:"spec"`
	Properties  meta.MapProperties `json:"properties"`
}

type MapSpec struct {

	// Name is passed to the kernel as a debug aid. Must only contain
	// alpha numeric and '_' characters.
	Name       string
	Type       ebpf.MapType
	KeySize    uint32
	ValueSize  uint32
	MaxEntries uint32

	// Flags is passed to the kernel and specifies additional map
	// creation attributes.
	Flags uint32

	// Automatically pin and load a map from MapOptions.PinPath.
	// Generates an error if an existing pinned map is incompatible with the MapSpec.
	Pinning ebpf.PinType
}

func (c *Component) Validate() error {
	if c.Name == "" {
		return errors.New("组件名称不能为空")
	}
	return nil
}
