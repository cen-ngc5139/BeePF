package export

import (
	"os"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
	"github.com/cilium/ebpf/btf"
)

type MyCustomHandler struct {
	logger *zap.Logger
}

// 实现 EventHandler 接口
func (h *MyCustomHandler) HandleEvent(ctx *UserContext, data *ReceivedEventData) error {
	switch data.Type {
	case TypeJsonText:
		h.logger.Info("received json data",
			zap.String("data", data.JsonText))
	case TypePlainText:
		h.logger.Info("received plain text",
			zap.String("data", data.Text))
	}
	return nil
}

func TestEventExporterBuilder_BuildForSingleValueWithTypeDescriptor(t *testing.T) {
	type fields struct {
		progFile           string
		binFile            string
		ExportFormat       ExportFormatType
		ExportEventHandler EventHandler
		UserCtx            *UserContext
	}
	type args struct {
		typeDesc     TypeDescriptor
		btfContainer *container.BTFContainer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *EventExporter
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				progFile:           "../../../../testdata/simple_prog.bpf.o",
				binFile:            "../../../../testdata/dumper_test.bin",
				ExportFormat:       FormatJson,
				ExportEventHandler: &MyCustomHandler{logger: zaptest.NewLogger(t)},
				UserCtx:            NewUserContext(0),
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

			generateMeta, err := meta.GenerateMeta(raw)
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
