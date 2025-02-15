package export

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
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

func TestDumpToJsonWithCheckedTypes(t *testing.T) {
	type args struct {
		specFile     string
		binFile      string
		checkedTypes []CheckedExportedMember
		data         []byte
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
			wantErr: false,
		},
		{
			name: "test dump sched_latency_t",
			args: args{
				specFile: "../../../../testdata/shepherd_x86_bpfel.o",
				binFile:  "../../../../testdata/shepherd_x86_bpfel.bin",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(tt.args.specFile)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			binData, err := os.ReadFile(tt.args.binFile)
			require.NoError(t, err, "Failed to read binary test data")

			generateMeta, err := meta.GenerateMeta(raw)
			if err != nil {
				t.Errorf("GenerateComposedObject() error = %v", err)
				return
			}

			checkedTypes, err := CheckExportTypesBtf(generateMeta.ExportTypes[0])
			if err != nil {
				t.Errorf("CheckExportTypesBtf() error = %v", err)
				return
			}

			got, err := DumpToJsonWithCheckedTypes(checkedTypes, binData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpToJsonWithCheckedTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var result test.ExampleTestStruct
			err = json.Unmarshal(got, &result)
			require.NoError(t, err, "Failed to unmarshal JSON")

			result.TestWithExampleData(t)
		})
	}
}

func TestJSONUnmarshalInt64Precision(t *testing.T) {
	// 测试用例：一个会导致精度丢失的大整数
	testValue := uint64(0x123456789abcdef0)

	// 将数值转换为 JSON
	jsonBytes, err := json.Marshal(testValue)
	if err != nil {
		t.Fatalf("failed to marshal test value: %v", err)
	}

	// 尝试解析回 interface{}
	var result interface{}
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	resultValue, err := result.(json.Number).Int64()
	if err != nil {
		t.Fatalf("failed to convert result to int64: %v", err)
	}
	if resultValue != int64(testValue) {
		t.Errorf("precision loss detected:\nexpected: %#x\ngot:      %#x\ndiff:     %#x",
			testValue, resultValue, testValue-uint64(resultValue))
	}

}

func TestDumpToStringWithCheckedTypes(t *testing.T) {
	type args struct {
		specFile     string
		binFile      string
		checkedTypes []CheckedExportedMember
		data         []byte
		out          *strings.Builder
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test dump int",
			args: args{
				specFile: "../../../../testdata/simple_prog.bpf.o",
				binFile:  "../../../../testdata/dumper_test.bin",
				out:      &strings.Builder{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(tt.args.specFile)
			if err != nil {
				t.Errorf("ReadFile() error = %v", err)
				return
			}

			generateMeta, err := meta.GenerateMeta(raw)
			if err != nil {
				t.Errorf("GenerateComposedObject() error = %v", err)
				return
			}

			checkedTypes, err := CheckExportTypesBtf(generateMeta.ExportTypes[0])
			if err != nil {
				t.Errorf("CheckExportTypesBtf() error = %v", err)
				return
			}

			binData, err := os.ReadFile(tt.args.binFile)
			require.NoError(t, err, "Failed to read binary test data")

			err = DumpToStringWithCheckedTypes(checkedTypes, binData, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpToStringWithCheckedTypes() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Log(tt.args.out.String())
		})
	}
}
