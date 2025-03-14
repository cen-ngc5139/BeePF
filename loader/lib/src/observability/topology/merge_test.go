package topology

import (
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"reflect"
	"testing"
)

func TestMergeTopology(t *testing.T) {
	tests := []struct {
		name    string
		want    meta.Topology
		wantErr bool
	}{
		{
			name: "fake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeTopology()
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeTopology() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeTopology() got = %v, want %v", got, tt.want)
			}
		})
	}
}
