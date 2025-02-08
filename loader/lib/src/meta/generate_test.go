package meta

import "testing"

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
			got, err := GenerateMeta(tt.args.objectFile)
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
