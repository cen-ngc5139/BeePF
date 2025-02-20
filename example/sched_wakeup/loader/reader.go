package loader

import (
	"fmt"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
)

// MapReader 定义 map 读取器接口
type MapReader interface {
	Read() (*perf.Record, error)
	Close() error
}

// PerfMapReader perf event reader 实现
type PerfMapReader struct {
	reader *perf.Reader
}

func NewPerfMapReader(m *ebpf.Map) (*PerfMapReader, error) {
	reader, err := perf.NewReader(m, os.Getpagesize())
	if err != nil {
		return nil, fmt.Errorf("create perf reader failed: %w", err)
	}
	return &PerfMapReader{reader: reader}, nil
}

func (r *PerfMapReader) Read() (*perf.Record, error) {
	return r.reader.Read()
}

func (r *PerfMapReader) Close() error {
	return r.reader.Close()
}

// RingMapReader ring buffer reader 实现
type RingMapReader struct {
	reader *ringbuf.Reader
}

func NewRingMapReader(m *ebpf.Map) (*RingMapReader, error) {
	reader, err := ringbuf.NewReader(m)
	if err != nil {
		return nil, fmt.Errorf("create ring buffer reader failed: %w", err)
	}
	return &RingMapReader{reader: reader}, nil
}

func (r *RingMapReader) Read() (*perf.Record, error) {
	record, err := r.reader.Read()
	if err != nil {
		return nil, err
	}
	// 转换为统一的 perf.Record 格式
	return &perf.Record{
		RawSample: record.RawSample,
		CPU:       record.CPU,
		Lost:      0,
	}, nil
}

func (r *RingMapReader) Close() error {
	return r.reader.Close()
}
