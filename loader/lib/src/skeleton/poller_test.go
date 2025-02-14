package skeleton

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"github.com/cilium/ebpf/ringbuf"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf"

	"github.com/cilium/ebpf/perf"
)

func TestPerfEventPoller_Poll(t *testing.T) {
	type fields struct {
		BinaryPath, BtfArchivePath string
		reader                     *perf.Reader
		processor                  EventProcessor
		errorFlag                  *atomic.Bool
		timeout                    time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "normol",
			fields: fields{
				BinaryPath:     "../../../testdata/shepherd_x86_bpfel.o",
				BtfArchivePath: "../../../testdata/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := meta.GenerateComposedObject(tt.fields.BinaryPath)
			if err != nil {
				t.Errorf("GenerateComposedObject() error = %v", err)
				return
			}

			preLoadBpfSkeleton, err := FromJsonPackage(pkg, tt.fields.BtfArchivePath).Build()
			if err != nil {
				t.Errorf("Build() error = %v", err)
				return
			}

			skeleton, err := preLoadBpfSkeleton.LoadAndAttach()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAndAttach() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, m := range skeleton.Collection.Maps {
				if m.Type() != ebpf.RingBuf {
					continue
				}

				perfReader, err := ringbuf.NewReader(m)
				if err != nil {
					t.Errorf("fail to new perf event array map  reader, err = %v", err)
				}

				for _, v := range skeleton.Collection.Variables {

					ee := export.NewEventExporterBuilder().
						SetExportFormat(export.FormatJson).
						SetUserContext(export.NewUserContext(0))

					structType, err := findStructType(v.Type())
					if err != nil {
						t.Fatalf("failed to find struct type: %v", err)
					}

					if structType.Name != "sched_latency_t" {
						continue
					}

					exporter, err := ee.BuildForSingleValueWithTypeDescriptor(
						&export.BTFTypeDescriptor{
							Type: structType,
							Name: structType.TypeName(),
						},
						skeleton.Btf,
					)
					if err != nil {
						t.Errorf("BuildForSingleValueWithTypeDescriptor() error = %v", err)
						return
					}

					jsonHandler := export.NewJsonExportEventHandler(exporter)

					p := &RingBufPoller{
						reader:    perfReader,
						processor: jsonHandler,
						timeout:   tt.fields.timeout,
					}
					if err := p.Poll(); (err != nil) != tt.wantErr {
						t.Errorf("Poll() error = %v, wantErr %v", err, tt.wantErr)
					}
				}

			}

		})
	}
}
