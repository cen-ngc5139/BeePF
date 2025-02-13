package export

import (
	"os"
	"reflect"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
)

func TestCheckExportTypesBtf(t *testing.T) {
	type args struct {
		progFile string
	}
	tests := []struct {
		name    string
		args    args
		want    []CheckedExportedMember
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				progFile: "../../../../testdata/shepherd_x86_bpfel.o",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(tt.args.progFile)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			generateMeta, err := meta.GenerateMeta(raw)
			if err != nil {
				t.Errorf("GenerateComposedObject() error = %v", err)
				return
			}

			got, err := CheckExportTypesBtf(generateMeta.ExportTypes[0])
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExportTypesBtf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckExportTypesBtf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
