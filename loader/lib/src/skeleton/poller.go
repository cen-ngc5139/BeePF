package skeleton

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
)

// EventProcessor 事件处理器接口
type EventProcessor interface {
	HandleEvent(data []byte) error
}

// SampleMapProcessor map 采样处理器接口
type SampleMapProcessor interface {
	HandleEvent(key []byte, value []byte) error
}

// Poller 轮询器接口
type Poller interface {
	Poll() error
	Close() error
}

// RingBufPoller ring buffer 轮询器
type RingBufPoller struct {
	reader    *ringbuf.Reader
	processor EventProcessor
	timeout   time.Duration
}

// PerfEventPoller perf event 轮询器
type PerfEventPoller struct {
	reader    *perf.Reader
	processor EventProcessor
	errorFlag atomic.Bool
	timeout   time.Duration
}

// SampleMapPoller map 采样轮询器
type SampleMapPoller struct {
	bpfMap       *ebpf.Map
	processor    SampleMapProcessor
	sampleConfig *MapSampleConfig
}

// MapSampleConfig map 采样配置
type MapSampleConfig struct {
	Interval int  `json:"interval"`
	ClearMap bool `json:"clear_map"`
}

// NewRingBufPoller 创建 ring buffer 轮询器
func NewRingBufPoller(bpfMap *ebpf.Map, processor EventProcessor, timeoutMs uint64) (*RingBufPoller, error) {
	reader, err := ringbuf.NewReader(bpfMap)
	if err != nil {
		return nil, fmt.Errorf("create ringbuf reader error: %w", err)
	}

	return &RingBufPoller{
		reader:    reader,
		processor: processor,
		timeout:   time.Duration(timeoutMs) * time.Millisecond,
	}, nil
}

// Poll 实现轮询方法
func (p *RingBufPoller) Poll() error {
	record, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("read ringbuf error: %w", err)
	}

	err = saveRawBytes("test.bin", record.RawSample)
	if err != nil {
		return err
	}

	if err := p.processor.HandleEvent(record.RawSample); err != nil {
		return fmt.Errorf("handle event error: %w", err)
	}

	return nil
}

// 直接将字节数组保存为文件的辅助函数
func saveRawBytes(filename string, data []byte) error {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	// 直接写入原始字节
	return os.WriteFile(filename, data, 0644)
}

// NewPerfEventPoller 创建 perf event 轮询器
func NewPerfEventPoller(bpfMap *ebpf.Map, processor EventProcessor, timeoutMs uint64) (*PerfEventPoller, error) {
	reader, err := perf.NewReader(bpfMap, os.Getpagesize())
	if err != nil {
		return nil, fmt.Errorf("create perf reader error: %w", err)
	}

	return &PerfEventPoller{
		reader:    reader,
		processor: processor,
		timeout:   time.Duration(timeoutMs) * time.Millisecond,
	}, nil
}

// Poll 实现轮询方法
func (p *PerfEventPoller) Poll() error {
	record, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("read perf event error: %w", err)
	}

	if err := p.processor.HandleEvent(record.RawSample); err != nil {
		p.errorFlag.Store(true)
		return fmt.Errorf("handle event error: %w", err)
	}

	return nil
}

// NewSampleMapPoller 创建 map 采样轮询器
func NewSampleMapPoller(bpfMap *ebpf.Map, processor SampleMapProcessor, config *MapSampleConfig) *SampleMapPoller {
	return &SampleMapPoller{
		bpfMap:       bpfMap,
		processor:    processor,
		sampleConfig: config,
	}
}

// Poll 实现轮询方法
func (p *SampleMapPoller) Poll() error {
	var key []byte
	var value []byte

	iter := p.bpfMap.Iterate()
	for iter.Next(&key, &value) {
		if err := p.processor.HandleEvent(key, value); err != nil {
			return fmt.Errorf("handle event error: %w", err)
		}
	}

	time.Sleep(time.Duration(p.sampleConfig.Interval) * time.Millisecond)
	return nil
}

// Close 清理资源
func (p *SampleMapPoller) Close() error {
	if p.sampleConfig.ClearMap {
		// 清理 map
		var key []byte
		iter := p.bpfMap.Iterate()
		for iter.Next(&key, nil) {
			if err := p.bpfMap.Delete(key); err != nil {
				return fmt.Errorf("delete map entry error: %w", err)
			}
		}
	}
	return nil
}

// findStructType 递归查找 BTF 中的结构体类型
func findStructType(t btf.Type) (*btf.Struct, error) {
	switch v := t.(type) {
	case *btf.Struct:
		// 直接是结构体类型
		return v, nil
	case *btf.Pointer:
		// 如果是指针，查找其目标类型
		return findStructType(v.Target)
	case *btf.Typedef:
		// 如果是类型别名，查找原始类型
		return findStructType(v.Type)
	case *btf.Volatile:
		// 如果是 volatile 修饰，查找基础类型
		return findStructType(v.Type)
	case *btf.Const:
		// 如果是 const 修饰，查找基础类型
		return findStructType(v.Type)
	case *btf.Var:
		return findStructType(v.Type)
	default:
		return nil, fmt.Errorf("unexpected type %T, cannot find struct", t)
	}
}
