package meta

import (
	"bytes"
	"os"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// TestProgMeta_AttachProgram 测试 ProgMeta 的 AttachProgram 方法
// ⚠️ 该单元测试需要使用 IDE 远程单元测试，通过 ssh 连接到目标 linux 机器上运行
func TestProgMeta_AttachProgram(t *testing.T) {
	type fields struct {
		BinaryPath string
		Name       string
		Attach     string
		Link       bool
	}
	type args struct {
		spec    *ebpf.ProgramSpec
		program *ebpf.Program
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    link.Link
		wantErr bool
	}{
		{
			name: "sched_wakeup",
			fields: fields{
				BinaryPath: "../../../testdata/shepherd_x86_bpfel.o",
				Name:       "sched_wakeup",
				Attach:     "tp_btf/sched_wakeup",
				Link:       true,
			},
			args: args{
				spec:    &ebpf.ProgramSpec{},
				program: &ebpf.Program{},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectRaw, err := os.ReadFile(tt.fields.BinaryPath)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			// 从对象文件创建 CollectionSpec
			spec, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(objectRaw))
			if err != nil {
				t.Errorf("load collection spec error: %v", err)
				return
			}

			collection, err := ebpf.NewCollection(spec)
			if err != nil {
				t.Errorf("NewCollection() error = %v", err)
				return
			}

			progSpec := spec.Programs[tt.fields.Name]
			if progSpec == nil {
				t.Errorf("program %s not found", tt.fields.Name)
				return
			}

			prog := collection.Programs[tt.fields.Name]
			if prog == nil {
				t.Errorf("program %s not found", tt.fields.Name)
				return
			}

			p := &ProgMeta{
				Name:   tt.fields.Name,
				Attach: tt.fields.Attach,
				Link:   tt.fields.Link,
			}
			got, err := p.AttachProgram(progSpec, prog)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttachProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("AttachProgram() got = %v, want %v", got, tt.want)
			}
		})
	}
}
