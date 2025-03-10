package component

import (
	"bytes"
	"fmt"
	"os"
	"time"

	loader "github.com/cen-ngc5139/BeePF/loader/lib/src/cli"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/metrics"
	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cilium/ebpf"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (o *Operator) UploadBinary() (err error) {
	if o.Binary == nil || len(o.Binary) == 0 {
		return errors.New("二进制文件为空")
	}

	// 校验文件是否为 ELF 文件
	if _, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(o.Binary)); err != nil {
		return errors.Wrap(err, "加载 ELF 文件失败")
	}

	// 使用 uuid 生成文件名
	fileName := uuid.New().String() + ".o"

	// 指定目录
	dir := "./binary"
	// 如果目录不存在，则创建目录
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	binaryPath := dir + "/" + fileName
	if err := os.WriteFile(binaryPath, o.Binary, 0644); err != nil {
		return errors.Wrap(err, "写入二进制文件失败")
	}

	o.Component.BinaryPath = binaryPath

	// 初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("初始化日志失败", zap.Error(err))
		return
	}
	defer logger.Sync()

	config := &loader.Config{
		ObjectPath:  binaryPath,
		Logger:      logger,
		StructName:  "sched_latency_t",
		PollTimeout: 100 * time.Millisecond,
		Properties: meta.Properties{
			Maps: map[string]*meta.Map{
				"sched_events": {
					Name:          "sched_events",
					ExportHandler: &export.MyCustomHandler{Logger: logger},
				},
			},
			Stats: &meta.Stats{
				Interval: 1 * time.Second,
				Handler:  metrics.NewDefaultHandler(logger),
			},
		},
	}

	bpfLoader := loader.NewBPFLoader(config)

	err = bpfLoader.Init()
	if err != nil {
		logger.Fatal("初始化 BPF 加载器失败", zap.Error(err))
		return
	}

	defer func() {
		if bpfLoader.Collection != nil {
			bpfLoader.Collection.Close()
		}
	}()

	// 将 bpfLoader.PreLoadSkeleton.Spec 中的数据转换为 component 组件
	component, err := o.convertSpecToComponent(bpfLoader.PreLoadSkeleton.Spec)
	if err != nil {
		logger.Error("转换 Spec 到组件失败", zap.Error(err))
		return errors.Wrap(err, "转换 Spec 到组件失败")
	}

	// 设置组件
	o.Component = component

	// 检查组件
	if err = o.checkComponent(); err != nil {
		return errors.Wrap(err, "组件校验失败")
	}

	// 创建组件
	if err = o.Create(); err != nil {
		return errors.Wrap(err, "创建组件失败")
	}

	return nil
}

// convertSpecToComponent 将 ebpf.CollectionSpec 转换为 models.Component
func (o *Operator) convertSpecToComponent(spec *ebpf.CollectionSpec) (*models.Component, error) {
	if spec == nil {
		return nil, errors.New("Spec 为空")
	}

	// 创建组件
	component := &models.Component{
		Name:       o.Component.Name,      // 使用第一个程序的名称作为组件名称
		ClusterId:  o.Component.ClusterId, // 默认集群ID，可以根据需要修改
		BinaryPath: o.Component.BinaryPath,
		Programs:   make([]models.Program, 0),
		Maps:       make([]models.Map, 0),
	}

	// 处理程序
	for name, progSpec := range spec.Programs {
		program := models.Program{
			Name:        name,
			Description: fmt.Sprintf("Program %s", name),
			Spec: models.ProgramSpec{
				Name:          progSpec.Name,
				Type:          progSpec.Type,
				AttachType:    progSpec.AttachType,
				AttachTo:      progSpec.AttachTo,
				SectionName:   progSpec.SectionName,
				Flags:         progSpec.Flags,
				License:       progSpec.License,
				KernelVersion: progSpec.KernelVersion,
			},
			Properties: meta.ProgramProperties{
				// 设置程序属性
				CGroupPath:  "",
				PinPath:     "",
				LinkPinPath: "",
				Tc:          nil,
			},
		}
		component.Programs = append(component.Programs, program)
	}

	// 处理映射
	for name, mapSpec := range spec.Maps {
		mapObj := models.Map{
			Name:        name,
			Description: fmt.Sprintf("Map %s", name),
			Spec: models.MapSpec{
				Name:       mapSpec.Name,
				Type:       mapSpec.Type,
				KeySize:    mapSpec.KeySize,
				ValueSize:  mapSpec.ValueSize,
				MaxEntries: mapSpec.MaxEntries,
				Flags:      mapSpec.Flags,
				Pinning:    mapSpec.Pinning,
			},
			Properties: meta.MapProperties{
				// 设置映射属性
				PinPath: "",
			},
		}
		component.Maps = append(component.Maps, mapObj)
	}

	return component, nil
}

// getFirstKey 获取 map 的第一个键
func getFirstKey(m map[string]*ebpf.ProgramSpec) string {
	for k := range m {
		return k
	}
	return ""
}
