package container_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/container"
)

func TestNewElfContainerFromBinary(t *testing.T) {

	type args struct {
		file string
	}
	tests := []struct {
		name     string
		args     args
		sections int
		wantErr  bool
	}{
		{
			name: "normal",
			args: args{
				file: "../../../testdata/rewrite.elf",
			},
			sections: 27,
			wantErr:  false,
		},
		{
			name: "shepherd",
			args: args{
				file: "../../../testdata/shepherd_x86_bpfel.o",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := os.ReadFile(tt.args.file)
			if err != nil {
				t.Errorf("NewElfContainerFromBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := container.NewElfContainerFromBinary(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewElfContainerFromBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got.File().Sections), tt.sections) {
				t.Errorf("NewElfContainerFromBinary() got elf file sections = %v, want %v", got, tt.sections)
			}
		})
	}
}
