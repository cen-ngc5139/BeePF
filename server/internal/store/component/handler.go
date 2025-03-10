package component

import (
	"github.com/cen-ngc5139/BeePF/server/internal/database"
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cen-ngc5139/BeePF/server/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GetComponentByID 根据ID获取组件及其所有关联数据
func (s *Store) Get(id uint64) (*models.Component, error) {
	var componentDB models.ComponentDB
	result := database.DB.Preload("Programs").
		Preload("Programs.Spec").
		Preload("Programs.Properties").
		Preload("Maps").
		Preload("Maps.Spec").
		Preload("Maps.Properties").
		Where("deleted = 0").
		First(&componentDB, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return componentDB.ToComponent(), nil
}

// ListComponents 获取组件列表
func (s *Store) List(query *utils.Query) (total int64, components []*models.Component, err error) {
	var componentsDB []models.ComponentDB

	// 计算总数
	result := database.DB.Model(&models.ComponentDB{}).Where("deleted = 0").Count(&total)
	if result.Error != nil {
		err = result.Error
		return
	}

	// 构建查询
	db := database.DB.Where("deleted = 0")

	// 应用分页
	if query != nil && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		if offset < 0 {
			offset = 0
		}
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 执行查询
	result = db.Find(&componentsDB)
	if result.Error != nil {
		err = result.Error
		return
	}

	components = make([]*models.Component, len(componentsDB))
	for i, c := range componentsDB {
		var current *models.Component
		current, err = s.Get(c.ID)
		if err != nil {
			err = errors.Wrapf(err, "获取组件 %d 失败", c.ID)
			return
		}

		components[i] = current
	}

	return total, components, nil
}

// CreateComponent 创建组件
func (s *Store) Create(component *models.Component) (*models.Component, error) {
	// 开启事务
	return component, database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建组件
		componentDB := &models.ComponentDB{
			Name:       component.Name,
			ClusterID:  uint64(component.ClusterId),
			BinaryPath: component.BinaryPath,
		}

		if err := tx.Create(componentDB).Error; err != nil {
			return err
		}

		// 创建程序
		for _, program := range component.Programs {
			programDB := &models.ProgramDB{
				ComponentID: componentDB.ID,
				Name:        program.Name,
				Description: program.Description,
			}

			if err := tx.Create(programDB).Error; err != nil {
				return err
			}

			// 创建程序规格
			programSpecDB := &models.ProgramSpecDB{
				ProgramID:     programDB.ID,
				Name:          program.Spec.Name,
				Type:          program.Spec.Type,
				AttachType:    program.Spec.AttachType,
				AttachTo:      program.Spec.AttachTo,
				SectionName:   program.Spec.SectionName,
				Flags:         program.Spec.Flags,
				License:       program.Spec.License,
				KernelVersion: program.Spec.KernelVersion,
			}

			if err := tx.Create(programSpecDB).Error; err != nil {
				return err
			}

			// 创建程序属性
			programPropertiesDB := &models.ProgramPropertiesDB{
				ProgramID:      programDB.ID,
				PropertiesJSON: models.JSONProgramProperties(program.Properties),
			}

			if err := tx.Create(programPropertiesDB).Error; err != nil {
				return err
			}
		}

		// 创建 Maps
		for _, m := range component.Maps {
			mapDB := &models.MapDB{
				ComponentID: componentDB.ID,
				Name:        m.Name,
				Description: m.Description,
			}

			if err := tx.Create(mapDB).Error; err != nil {
				return err
			}

			// 创建 Map 规格
			mapSpecDB := &models.MapSpecDB{
				MapID:      mapDB.ID,
				Name:       m.Spec.Name,
				Type:       m.Spec.Type,
				KeySize:    m.Spec.KeySize,
				ValueSize:  m.Spec.ValueSize,
				MaxEntries: m.Spec.MaxEntries,
				Flags:      m.Spec.Flags,
				Pinning:    m.Spec.Pinning.String(),
			}

			if err := tx.Create(mapSpecDB).Error; err != nil {
				return err
			}

			// 创建 Map 属性
			mapPropertiesDB := &models.MapPropertiesDB{
				MapID:          mapDB.ID,
				PropertiesJSON: models.JSONMapProperties(m.Properties),
			}

			if err := tx.Create(mapPropertiesDB).Error; err != nil {
				return err
			}
		}

		// 重新查询完整的组件
		return tx.Preload("Programs").
			Preload("Programs.Spec").
			Preload("Programs.Properties").
			Preload("Maps").
			Preload("Maps.Spec").
			Preload("Maps.Properties").
			First(componentDB, componentDB.ID).Error
	})
}

// DeleteComponent 删除组件
func (s *Store) Delete(component *models.Component) (err error) {
	// 使用事务确保操作的原子性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取组件关联的所有程序
		var programIDs []uint64
		if err := tx.Model(&models.ProgramDB{}).
			Where("component_id = ? AND deleted = 0", component.Id).
			Pluck("id", &programIDs).Error; err != nil {
			return err
		}

		// 2. 获取组件关联的所有映射
		var mapIDs []uint64
		if err := tx.Model(&models.MapDB{}).
			Where("component_id = ? AND deleted = 0", component.Id).
			Pluck("id", &mapIDs).Error; err != nil {
			return err
		}

		// 3. 标记程序为已删除
		if len(programIDs) > 0 {
			// 3.1 标记程序为已删除
			if err := tx.Model(&models.ProgramDB{}).
				Where("id IN ?", programIDs).
				Update("deleted", 1).Error; err != nil {
				return err
			}

			// 3.2 标记程序规格为已删除
			if err := tx.Model(&models.ProgramSpecDB{}).
				Where("program_id IN ?", programIDs).
				Update("deleted", 1).Error; err != nil {
				return err
			}

			// 3.3 标记程序属性为已删除
			if err := tx.Model(&models.ProgramPropertiesDB{}).
				Where("program_id IN ?", programIDs).
				Update("deleted", 1).Error; err != nil {
				return err
			}
		}

		// 4. 标记映射为已删除
		if len(mapIDs) > 0 {
			// 4.1 标记映射为已删除
			if err := tx.Model(&models.MapDB{}).
				Where("id IN ?", mapIDs).
				Update("deleted", 1).Error; err != nil {
				return err
			}

			// 4.2 标记映射规格为已删除
			if err := tx.Model(&models.MapSpecDB{}).
				Where("map_id IN ?", mapIDs).
				Update("deleted", 1).Error; err != nil {
				return err
			}

			// 4.3 标记映射属性为已删除
			if err := tx.Model(&models.MapPropertiesDB{}).
				Where("map_id IN ?", mapIDs).
				Update("deleted", 1).Error; err != nil {
				return err
			}
		}

		// 5. 标记组件为已删除
		return tx.Model(&models.ComponentDB{}).
			Where("id = ?", component.Id).
			Update("deleted", 1).Error
	})
}
