// pkg/tests/tests.go

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// ExampleTestStruct 用于测试的示例结构体
type ExampleTestStruct struct {
	Arr1   [2][3][4]int32 `json:"arr1"`
	Str    string         `json:"str"`
	StrArr [10]string     `json:"str_arr"`
	Ft     float32        `json:"ft"`
	Dbl    float64        `json:"dbl"`
	U8v    uint8          `json:"u8v"`
	I8v    int8           `json:"i8v"`
	U16v   uint16         `json:"u16v"`
	I16v   int16          `json:"i16v"`
	U32v   uint32         `json:"u32v"`
	I32v   int32          `json:"i32v"`
	U64v   uint64         `json:"u64v"`
	I64v   int64          `json:"i64v"`
	E      string         `json:"e"`
}

// TestWithExampleData 测试示例数据
func (s *ExampleTestStruct) TestWithExampleData(t *testing.T) {
	// 测试浮点数
	if s.Ft != 1.23 {
		t.Errorf("ft: expected 1.23, got %f", s.Ft)
	}
	if s.Dbl != 4.56 {
		t.Errorf("dbl: expected 4.56, got %f", s.Dbl)
	}

	// 测试整数
	if s.U8v != 0x12 {
		t.Errorf("u8v: expected 0x12, got 0x%x", s.U8v)
	}
	if s.U16v != 0x1234 {
		t.Errorf("u16v: expected 0x1234, got 0x%x", s.U16v)
	}
	if s.U32v != 0x12345678 {
		t.Errorf("u32v: expected 0x12345678, got 0x%x", s.U32v)
	}
	if s.U64v != 0x123456789abcdef0 {
		t.Errorf("u64v: expected 0x123456789abcdef0, got 0x%x", s.U64v)
	}

	// 测试有符号整数
	if s.I8v != -0x12 {
		t.Errorf("i8v: expected -0x12, got %d", s.I8v)
	}
	if s.I16v != -0x1234 {
		t.Errorf("i16v: expected -0x1234, got %d", s.I16v)
	}
	if s.I32v != -0x12345678 {
		t.Errorf("i32v: expected -0x12345678, got %d", s.I32v)
	}
	if s.I64v != -0x123456789abcdef0 {
		t.Errorf("i64v: expected -0x123456789abcdef0, got %d", s.I64v)
	}

	// 测试字符串
	if s.E != "E_A(0)" {
		t.Errorf("e: expected E_A(0), got %s", s.E)
	}
	if s.Str != "A-String" {
		t.Errorf("str: expected A-String, got %s", s.Str)
	}

	// 测试多维数组
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 4; k++ {
				expected := int32((i << 16) + (j << 8) + k)
				if s.Arr1[i][j][k] != expected {
					t.Errorf("arr1[%d][%d][%d]: expected %d, got %d",
						i, j, k, expected, s.Arr1[i][j][k])
				}
			}
		}
	}

	// 测试字符串数组
	for i := 0; i < 10; i++ {
		expected := fmt.Sprintf("hello %d", i)
		if s.StrArr[i] != expected {
			t.Errorf("str_arr[%d]: expected %s, got %s",
				i, expected, s.StrArr[i])
		}
	}
}

// GetAssetsDir 获取测试资源目录
func GetAssetsDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		assetsPath := filepath.Join(dir, "assets")
		if _, err := os.Stat(assetsPath); err == nil {
			return assetsPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("Could not find assets directory")
		}
		dir = parent
	}
}
