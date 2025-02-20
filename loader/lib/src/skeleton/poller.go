package skeleton

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
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
	GetPollFunc() PollFunc
}

// RingBufPoller ring buffer 轮询器
type RingBufPoller struct {
	Reader    *ringbuf.Reader
	Processor EventProcessor
	Timeout   time.Duration
}

// PerfEventPoller perf event 轮询器
type PerfEventPoller struct {
	Reader    *perf.Reader
	Processor EventProcessor
	ErrorFlag atomic.Bool
	Timeout   time.Duration
}

// SampleMapPoller map 采样轮询器
type SampleMapPoller struct {
	BpfMap       *ebpf.Map
	Processor    SampleMapProcessor
	SampleConfig *MapSampleConfig
}

// MapSampleConfig map 采样配置
type MapSampleConfig struct {
	Interval int  `json:"interval"`
	ClearMap bool `json:"clear_map"`
}

type ProgramPoller struct {
	// 轮询控制
	stopChan chan struct{}
	wg       sync.WaitGroup

	// 错误处理
	errChan chan error

	// 轮询配置
	interval time.Duration
}

// NewProgramPoller 创建新的轮询器
func NewProgramPoller(interval time.Duration) *ProgramPoller {
	return &ProgramPoller{
		stopChan: make(chan struct{}),
		errChan:  make(chan error, 1),
		interval: interval,
	}
}

// PollFunc 轮询函数类型
type PollFunc func() error

// StartPolling 开始轮询
func (p *ProgramPoller) StartPolling(
	name string,
	pollFn PollFunc,
	errorHandler func(error),
) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

		for {
			select {
			case <-p.stopChan:
				return

			case <-ticker.C:
				// 执行轮询函数
				if err := pollFn(); err != nil {
					if errorHandler != nil {
						errorHandler(err)
					}
					// 发送错误到错误通道
					select {
					case p.errChan <- fmt.Errorf("poll %s error: %w", name, err):
					default:
						// 错误通道已满，记录日志
						log.Printf("Error polling %s: %v", name, err)
					}
				}
			}
		}
	}()
}

// Stop 停止轮询
func (p *ProgramPoller) Stop() {
	close(p.stopChan)
	p.wg.Wait()
}

// Errors 返回错误通道
func (p *ProgramPoller) Errors() <-chan error {
	return p.errChan
}

// NewRingBufPoller 创建 ring buffer 轮询器
func NewRingBufPoller(bpfMap *ebpf.Map, processor EventProcessor, timeoutMs uint64) (*RingBufPoller, error) {
	reader, err := ringbuf.NewReader(bpfMap)
	if err != nil {
		return nil, fmt.Errorf("create ringbuf reader error: %w", err)
	}

	return &RingBufPoller{
		Reader:    reader,
		Processor: processor,
		Timeout:   time.Duration(timeoutMs) * time.Millisecond,
	}, nil
}

func (p *RingBufPoller) GetPollFunc() PollFunc {
	return func() error {
		return p.Poll()
	}
}

// Poll 实现轮询方法
func (p *RingBufPoller) Poll() error {
	record, err := p.Reader.Read()
	if err != nil {
		return fmt.Errorf("read ringbuf error: %w", err)
	}

	err = saveRawBytes("test.bin", record.RawSample)
	if err != nil {
		return err
	}

	if err := p.Processor.HandleEvent(record.RawSample); err != nil {
		return fmt.Errorf("handle event error: %w", err)
	}

	return nil
}

func (p *RingBufPoller) Close() error {
	return p.Reader.Close()
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
		Reader:    reader,
		Processor: processor,
		Timeout:   time.Duration(timeoutMs) * time.Millisecond,
	}, nil
}

// Poll 实现轮询方法
func (p *PerfEventPoller) Poll() error {
	record, err := p.Reader.Read()
	if err != nil {
		return fmt.Errorf("read perf event error: %w", err)
	}

	if err := p.Processor.HandleEvent(record.RawSample); err != nil {
		p.ErrorFlag.Store(true)
		return fmt.Errorf("handle event error: %w", err)
	}

	return nil
}

func (p *PerfEventPoller) GetPollFunc() PollFunc {
	return func() error {
		return p.Poll()
	}
}

func (p *PerfEventPoller) Close() error {
	return p.Reader.Close()
}

// NewSampleMapPoller 创建 map 采样轮询器
func NewSampleMapPoller(bpfMap *ebpf.Map, processor SampleMapProcessor, config *MapSampleConfig) *SampleMapPoller {
	return &SampleMapPoller{
		BpfMap:       bpfMap,
		Processor:    processor,
		SampleConfig: config,
	}
}

// Poll 实现轮询方法
func (p *SampleMapPoller) Poll() error {
	var key []byte
	var value []byte

	iter := p.BpfMap.Iterate()
	for iter.Next(&key, &value) {
		if err := p.Processor.HandleEvent(key, value); err != nil {
			return fmt.Errorf("handle event error: %w", err)
		}
	}

	time.Sleep(time.Duration(p.SampleConfig.Interval) * time.Millisecond)
	return nil
}

func (p *SampleMapPoller) GetPollFunc() PollFunc {
	return func() error {
		return p.Poll()
	}
}

// Close 清理资源
func (p *SampleMapPoller) Close() error {
	if p.SampleConfig.ClearMap {
		// 清理 map
		var key []byte
		iter := p.BpfMap.Iterate()
		for iter.Next(&key, nil) {
			if err := p.BpfMap.Delete(key); err != nil {
				return fmt.Errorf("delete map entry error: %w", err)
			}
		}
	}
	return nil
}

// findStructType 递归查找 BTF 中的结构体类型
func FindStructType(t btf.Type) (*btf.Struct, error) {
	switch v := t.(type) {
	case *btf.Struct:
		// 直接是结构体类型
		return v, nil
	case *btf.Pointer:
		// 如果是指针，查找其目标类型
		return FindStructType(v.Target)
	case *btf.Typedef:
		// 如果是类型别名，查找原始类型
		return FindStructType(v.Type)
	case *btf.Volatile:
		// 如果是 volatile 修饰，查找基础类型
		return FindStructType(v.Type)
	case *btf.Const:
		// 如果是 const 修饰，查找基础类型
		return FindStructType(v.Type)
	case *btf.Var:
		return FindStructType(v.Type)
	default:
		return nil, fmt.Errorf("unexpected type %T, cannot find struct", t)
	}
}
