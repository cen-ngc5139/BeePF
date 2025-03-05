package component

import (
	"github.com/cen-ngc5139/BeePF/server/internal/database"
	"github.com/cen-ngc5139/BeePF/server/models"
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
func (s *Store) List() (total int64, components []*models.Component, err error) {
	var componentsDB []models.ComponentDB
	result := database.DB.Where("deleted = 0").Find(&componentsDB)

	if result.Error != nil {
		err = result.Error
		return
	}

	components = make([]*models.Component, len(componentsDB))
	for i, c := range componentsDB {
		components[i] = c.ToComponent()
	}

	return int64(len(componentsDB)), components, nil
}

// CreateComponent 创建组件
func (s *Store) Create(component *models.Component) (*models.Component, error) {
	// 开启事务
	return component, database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建组件
		componentDB := &models.ComponentDB{
			Name: component.Name,
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
				Type:          program.Spec.Type.String(),
				AttachType:    program.Spec.AttachType.String(),
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
				Type:       m.Spec.Type.String(),
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
