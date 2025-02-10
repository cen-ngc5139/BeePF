package meta

import (
	"os"
	"reflect"
	"testing"
)

func TestGenerateMeta(t *testing.T) {
	type args struct {
		objectFile string
	}
	tests := []struct {
		name    string
		args    args
		want    []ExportedTypesStructMeta
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				objectFile: "../../../testdata/shepherd_x86_bpfel.o",
			},
			want: []ExportedTypesStructMeta{
				{
					Name: "sched_latency_t",
					Members: []ExportedTypesStructMemberMeta{
						{Name: "pid", Type: "__u32"},
						{Name: "tid", Type: "__u32"},
						{Name: "delay_ns", Type: "__u64"},
						{Name: "ts", Type: "__u64"},
						{Name: "preempted_pid", Type: "__u32"},
						{Name: "preempted_comm", Type: "char[16]"},
						{Name: "is_preempt", Type: "__u64"},
						{Name: "comm", Type: "char[16]"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objectFile, err := os.ReadFile(tt.args.objectFile)
			if err != nil {
				t.Errorf("Failed to read object file: %v", err)
				return
			}

			got, err := GenerateMeta(objectFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.ExportTypes) != len(tt.want) {
				t.Errorf("GenerateMeta() got = %v, want %v", got, tt.want)
			}
			for i, v := range got.ExportTypes {
				if v.Name != tt.want[i].Name {
					t.Errorf("GenerateMeta() got = %v, want %v", v, tt.want[i])
				}

				for j, v := range v.Members {
					if v.Name != tt.want[i].Members[j].Name {
						t.Errorf("GenerateMeta() got = %v, want %v", v, tt.want[i].Members[j])
					}
				}
			}
		})
	}
}

func TestGenerateComposedObject(t *testing.T) {
	type args struct {
		objectFile string
	}
	tests := []struct {
		name    string
		args    args
		want    *ComposedObject
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				objectFile: "../../../testdata/shepherd_x86_bpfel.o",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateComposedObject(tt.args.objectFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateComposedObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			inner, err := got.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			// 将结果写入JSON文件
			jsonFile := "../../../testdata/shepherd_x86_bpfel.json"

			if err := os.WriteFile(jsonFile, inner, 0644); err != nil {
				t.Errorf("Failed to write JSON file: %v", err)
				return
			}

			unmarshal, err := os.ReadFile(jsonFile)
			if err != nil {
				t.Errorf("Failed to read JSON file: %v", err)
				return
			}

			got2 := &ComposedObject{}
			err = got2.UnmarshalJSON(unmarshal)
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, got2) {
				t.Errorf("GenerateComposedObject() got = %v, want %v", got, got2)
			}

		})
	}
}
