package topology

import (
	"reflect"
	"testing"
)

func TestListAllPrograms(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name: "fake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListAllPrograms()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAllPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListAllPrograms() got = %v, want %v", got, tt.want)
			}
		})
	}
}
