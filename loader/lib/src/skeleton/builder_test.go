package skeleton

import (
	"os"
	"reflect"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
)

func TestBpfSkeletonBuilder_Build(t *testing.T) {
	type fields struct {
		btfArchivePath string
		objectMeta     *meta.EunomiaObjectMeta
		objectFile     string
		bpfObject      []byte
		runnerConfig   *meta.RunnerConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    *PreLoadBpfSkeleton
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				btfArchivePath: "../../../testdata/",
				objectFile:     "../../../testdata/shepherd_x86_bpfel.o",
				objectMeta:     nil,
				bpfObject:      nil,
				runnerConfig:   nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectFile, err := os.ReadFile(tt.fields.objectFile)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			b := &BpfSkeletonBuilder{
				btfArchivePath: tt.fields.btfArchivePath,
				objectMeta:     tt.fields.objectMeta,
				bpfObject:      objectFile,
				runnerConfig:   tt.fields.runnerConfig,
			}
			got, err := b.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
		})
	}
}
