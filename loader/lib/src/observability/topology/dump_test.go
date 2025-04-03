package topology

import (
	"reflect"
	"testing"

	"github.com/cilium/ebpf"
)

func TestGetProgDumpXlated(t *testing.T) {
	type args struct {
		progID ebpf.ProgramID
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				progID: 18585,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProgDumpXlated(tt.args.progID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProgDumpXlated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProgDumpXlated() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProgDumpJited(t *testing.T) {
	type args struct {
		progID ebpf.ProgramID
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				progID: 18585,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProgDumpJited(tt.args.progID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProgDumpJited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProgDumpJited() got = %v, want %v", got, tt.want)
			}
		})
	}
}
