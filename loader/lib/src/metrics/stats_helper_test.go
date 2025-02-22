package metrics

import "testing"

func TestIsBPFStatsEnabled(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "linux-x86",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := EnableBPFStats(); err != nil {
				t.Errorf("fail to enable bpf stats")
			}

		})
	}
}
