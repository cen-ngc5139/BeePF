package export

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/test"
	"github.com/cilium/ebpf/btf"
	"github.com/stretchr/testify/require"
)

func TestDumpToJson(t *testing.T) {
	type args struct {
		specFile string
		binFile  string
		spec     *btf.Spec
		typ      btf.Type
		data     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "test dump int",
			args: args{
				specFile: "../../../../testdata/simple_prog.bpf.o",
				binFile:  "../../../../testdata/dumper_test.bin",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 读取 BPF 对象文件
			spec, err := btf.LoadSpec(tt.args.specFile)
			require.NoError(t, err, "Failed to load BTF spec")

			// 读取二进制测试数据
			binData, err := os.ReadFile(tt.args.binFile)
			require.NoError(t, err, "Failed to read binary test data")

			typ, err := spec.TypeByID(1)
			require.NoError(t, err, "Failed to get type by ID")

			got, err := DumpToJson(typ, binData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var result test.ExampleTestStruct
			err = json.Unmarshal(got, &result)
			require.NoError(t, err, "Failed to unmarshal JSON")

			result.TestWithExampleData(t)
		})
	}
}
