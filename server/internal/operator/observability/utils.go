package observability

import (
	"errors"
	"time"

	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/cilium/ebpf"
	"github.com/shirou/gopsutil/host"
)

func ConvertProgToWrapper(prog *ebpf.ProgramInfo) (models.ProgramInfoWrapper, error) {
	if prog == nil {
		return models.ProgramInfoWrapper{}, errors.New("program is nil")
	}

	wrapper := models.ProgramInfoWrapper{
		Name: prog.Name,
		Type: prog.Type,
		Tag:  prog.Tag,
	}

	id, ok := prog.ID()
	if ok {
		wrapper.ID = id
	}

	maps, ok := prog.MapIDs()
	if ok {
		wrapper.Maps = maps
	}

	btfID, ok := prog.BTFID()
	if ok {
		wrapper.BTF = btfID
	}

	loadTime, ok := prog.LoadTime()
	if ok {
		wrapper.LoadTime = getActualLoadTime(loadTime)
	}

	createdByUID, haveCreatedByUID := prog.CreatedByUID()
	if haveCreatedByUID {
		wrapper.CreatedByUID = createdByUID
	}

	return wrapper, nil
}

func ConvertMapToWrapper(mapInfo *ebpf.MapInfo) (models.MapInfoWrapper, error) {
	if mapInfo == nil {
		return models.MapInfoWrapper{}, errors.New("map info is nil")
	}

	wrapper := models.MapInfoWrapper{
		Name:       mapInfo.Name,
		Type:       mapInfo.Type,
		KeySize:    mapInfo.KeySize,
		ValueSize:  mapInfo.ValueSize,
		MaxEntries: mapInfo.MaxEntries,
		Flags:      mapInfo.Flags,
	}

	id, ok := mapInfo.ID()
	if ok {
		wrapper.ID = id
	}

	btfID, ok := mapInfo.BTFID()
	if ok {
		wrapper.BTF = btfID
	}

	mapExtra, ok := mapInfo.MapExtra()
	if ok {
		wrapper.MapExtra = mapExtra
	}

	memlock, ok := mapInfo.Memlock()
	if ok {
		wrapper.Memlock = memlock
	}

	wrapper.Frozen = mapInfo.Frozen()

	return wrapper, nil
}

func getActualLoadTime(loadTimeSinceBoot time.Duration) time.Time {
	// 获取系统启动时间
	bootTime, _ := host.BootTime()
	bootTimeUnix := time.Unix(int64(bootTime), 0)

	// 计算实际加载时间
	actualLoadTime := bootTimeUnix.Add(loadTimeSinceBoot)

	return actualLoadTime
}
