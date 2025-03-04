package export

import (
	"os"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf/btf"
	"go.uber.org/zap/zaptest"
)

func TestEventExporterBuilder_BuildForSingleValueWithTypeDescriptor(t *testing.T) {
	type fields struct {
		progFile           string
		binFile            string
		ExportFormat       ExportFormatType
		ExportEventHandler meta.EventHandler
		UserCtx            *meta.UserContext
	}
	tests := []struct {
		name    string
		fields  fields
		want    *EventExporter
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				progFile:           "../../../../testdata/simple_prog.bpf.o",
				binFile:            "../../../../testdata/dumper_test.bin",
				ExportFormat:       FormatJson,
				ExportEventHandler: &MyCustomHandler{Logger: zaptest.NewLogger(t)},
				UserCtx:            meta.NewUserContext(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(tt.fields.progFile)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			bin, err := os.ReadFile(tt.fields.binFile)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			btfContainer, err := container.NewBTFContainerFromBinary(raw)
			if err != nil {
				t.Errorf("NewBTFContainer() error = %v", err)
				return
			}

			generateMeta, err := meta.GenerateMeta(raw, meta.Properties{})
			if err != nil {
				t.Errorf("GenerateComposedObject() error = %v", err)
				return
			}

			var targetType btf.Type
			if ptr, ok := generateMeta.ExportTypes[0].Type.(*btf.Var); ok {
				targetType = ptr.Type
			} else {
				t.Errorf("expected pointer type, got %T", generateMeta.ExportTypes[0].Type)
				return
			}

			structType, ok := targetType.(*btf.Pointer)
			if !ok {
				t.Errorf("expected pointer type, got %T", targetType)
				return
			}

			typ, ok := structType.Target.(*btf.Struct)
			if !ok {
				t.Errorf("expected struct type, got %T", structType.Target)
				return
			}

			typeDesc := &BTFTypeDescriptor{
				Type: typ,
				Name: "S",
			}

			b := &EventExporterBuilder{
				ExportFormat:       tt.fields.ExportFormat,
				ExportEventHandler: tt.fields.ExportEventHandler,
				UserCtx:            tt.fields.UserCtx,
			}
			got, err := b.BuildForSingleValueWithTypeDescriptor(typeDesc, btfContainer)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildForSingleValueWithTypeDescriptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = got.InternalImpl.Process(bin)
			if err != nil {
				t.Errorf("Process() error = %v", err)
				return
			}
		})
	}
}
