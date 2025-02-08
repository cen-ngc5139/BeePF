package container

import (
	"os"
	"testing"
)

func TestNewBtfContainerFromBinary(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
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

			_, err = NewBTFContainerFromBinary(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBTFContainerFromBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
