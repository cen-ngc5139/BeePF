package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"github.com/cilium/ebpf"
)

// ComponentDB 组件数据库模型
type ComponentDB struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"column:name;uniqueIndex" json:"name"`
	BinaryPath     string    `gorm:"column:binary_path" json:"binary_path"`
	Deleted        uint8     `gorm:"column:deleted;default:0" json:"deleted"`
	Creator        string    `gorm:"column:creator" json:"creator"`
	CreatedTime    time.Time `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`

	ClusterID uint64 `gorm:"column:cluster_id" json:"cluster_id"`

	// 关联关系
	Programs []ProgramDB `gorm:"foreignKey:ComponentID" json:"programs"`
	Maps     []MapDB     `gorm:"foreignKey:ComponentID" json:"maps"`
}

// TableName 指定表名
func (ComponentDB) TableName() string {
	return "beepf.component"
}

// ProgramDB 程序数据库模型
type ProgramDB struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ComponentID    uint64    `gorm:"column:component_id;index" json:"component_id"`
	Name           string    `gorm:"column:name" json:"name"`
	Description    string    `gorm:"column:description" json:"description"`
	Deleted        uint8     `gorm:"column:deleted;default:0" json:"deleted"`
	Creator        string    `gorm:"column:creator" json:"creator"`
	CreatedTime    time.Time `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`

	// 关联关系
	Spec       ProgramSpecDB       `gorm:"foreignKey:ProgramID" json:"spec"`
	Properties ProgramPropertiesDB `gorm:"foreignKey:ProgramID" json:"properties"`
}

// TableName 指定表名
func (ProgramDB) TableName() string {
	return "beepf.program"
}

// ProgramSpecDB 程序规格数据库模型
type ProgramSpecDB struct {
	ID             uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProgramID      uint64           `gorm:"column:program_id;uniqueIndex" json:"program_id"`
	Name           string           `gorm:"column:name" json:"name"`
	Type           ebpf.ProgramType `gorm:"column:type" json:"type"`
	AttachType     ebpf.AttachType  `gorm:"column:attach_type" json:"attach_type"`
	AttachTo       string           `gorm:"column:attach_to" json:"attach_to"`
	SectionName    string           `gorm:"column:section_name" json:"section_name"`
	Flags          uint32           `gorm:"column:flags" json:"flags"`
	License        string           `gorm:"column:license" json:"license"`
	KernelVersion  uint32           `gorm:"column:kernel_version" json:"kernel_version"`
	Deleted        uint8            `gorm:"column:deleted;default:0" json:"deleted"`
	CreatedTime    time.Time        `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time        `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`
}

// TableName 指定表名
func (ProgramSpecDB) TableName() string {
	return "beepf.program_spec"
}

// ProgramPropertiesDB 程序属性数据库模型
type ProgramPropertiesDB struct {
	ID             uint64                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProgramID      uint64                `gorm:"column:program_id;uniqueIndex" json:"program_id"`
	PropertiesJSON JSONProgramProperties `gorm:"column:properties_json" json:"properties_json"`
	Deleted        uint8                 `gorm:"column:deleted;default:0" json:"deleted"`
	CreatedTime    time.Time             `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time             `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`
}

// TableName 指定表名
func (ProgramPropertiesDB) TableName() string {
	return "beepf.program_properties"
}

// MapDB Map数据库模型
type MapDB struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ComponentID    uint64    `gorm:"column:component_id;index" json:"component_id"`
	Name           string    `gorm:"column:name" json:"name"`
	Description    string    `gorm:"column:description" json:"description"`
	Deleted        uint8     `gorm:"column:deleted;default:0" json:"deleted"`
	Creator        string    `gorm:"column:creator" json:"creator"`
	CreatedTime    time.Time `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`

	// 关联关系
	Spec       MapSpecDB       `gorm:"foreignKey:MapID" json:"spec"`
	Properties MapPropertiesDB `gorm:"foreignKey:MapID" json:"properties"`
}

// TableName 指定表名
func (MapDB) TableName() string {
	return "beepf.map"
}

// MapSpecDB Map规格数据库模型
type MapSpecDB struct {
	ID             uint64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MapID          uint64       `gorm:"column:map_id;uniqueIndex" json:"map_id"`
	Name           string       `gorm:"column:name" json:"name"`
	Type           ebpf.MapType `gorm:"column:type" json:"type"`
	KeySize        uint32       `gorm:"column:key_size" json:"key_size"`
	ValueSize      uint32       `gorm:"column:value_size" json:"value_size"`
	MaxEntries     uint32       `gorm:"column:max_entries" json:"max_entries"`
	Flags          uint32       `gorm:"column:flags" json:"flags"`
	Pinning        string       `gorm:"column:pinning" json:"pinning"`
	Deleted        uint8        `gorm:"column:deleted;default:0" json:"deleted"`
	CreatedTime    time.Time    `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time    `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`
}

// TableName 指定表名
func (MapSpecDB) TableName() string {
	return "beepf.map_spec"
}

// MapPropertiesDB Map属性数据库模型
type MapPropertiesDB struct {
	ID             uint64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MapID          uint64            `gorm:"column:map_id;uniqueIndex" json:"map_id"`
	PropertiesJSON JSONMapProperties `gorm:"column:properties_json" json:"properties_json"`
	Deleted        uint8             `gorm:"column:deleted;default:0" json:"deleted"`
	CreatedTime    time.Time         `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	LastUpdateTime time.Time         `gorm:"column:last_update_time;autoUpdateTime" json:"last_update_time"`
}

// TableName 指定表名
func (MapPropertiesDB) TableName() string {
	return "beepf.map_properties"
}

// JSONProgramProperties 用于存储 ProgramProperties 的 JSON 类型
type JSONProgramProperties meta.ProgramProperties

// Value 实现 driver.Valuer 接口
func (j JSONProgramProperties) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONProgramProperties) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &j)
}

// JSONMapProperties 用于存储 MapProperties 的 JSON 类型
type JSONMapProperties meta.MapProperties

// Value 实现 driver.Valuer 接口
func (j JSONMapProperties) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONMapProperties) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &j)
}

// 模型转换函数

// ToComponent 将数据库模型转换为业务模型
func (c *ComponentDB) ToComponent() *Component {
	component := &Component{
		Id:        int(c.ID),
		Name:      c.Name,
		ClusterId: int64(c.ClusterID),
	}

	// 转换 Programs
	if len(c.Programs) > 0 {
		component.Programs = make([]Program, len(c.Programs))
		for i, p := range c.Programs {
			component.Programs[i] = *p.ToProgram()
		}
	}

	// 转换 Maps
	if len(c.Maps) > 0 {
		component.Maps = make([]Map, len(c.Maps))
		for i, m := range c.Maps {
			component.Maps[i] = *m.ToMap()
		}
	}

	return component
}

// ToProgram 将数据库模型转换为业务模型
func (p *ProgramDB) ToProgram() *Program {
	return &Program{
		Id:          int(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Spec:        *p.Spec.ToProgramSpec(),
		Properties:  meta.ProgramProperties(p.Properties.PropertiesJSON),
	}
}

// ToProgramSpec 将数据库模型转换为业务模型
func (ps *ProgramSpecDB) ToProgramSpec() *ProgramSpec {
	return &ProgramSpec{
		Name:          ps.Name,
		Type:          ps.Type,
		AttachType:    ps.AttachType,
		AttachTo:      ps.AttachTo,
		SectionName:   ps.SectionName,
		Flags:         ps.Flags,
		License:       ps.License,
		KernelVersion: ps.KernelVersion,
	}
}

// ToMap 将数据库模型转换为业务模型
func (m *MapDB) ToMap() *Map {
	return &Map{
		Id:          int(m.ID),
		Name:        m.Name,
		Description: m.Description,
		Spec:        *m.Spec.ToMapSpec(),
		Properties:  meta.MapProperties(m.Properties.PropertiesJSON),
	}
}

// ToMapSpec 将数据库模型转换为业务模型
func (ms *MapSpecDB) ToMapSpec() *MapSpec {
	return &MapSpec{
		Name:       ms.Name,
		Type:       ms.Type,
		KeySize:    ms.KeySize,
		ValueSize:  ms.ValueSize,
		MaxEntries: ms.MaxEntries,
		Flags:      ms.Flags,
		Pinning:    ebpf.PinType(0), // 需要转换字符串到枚举
	}
}
