package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"sync"
)

type InstanceMap struct {
	lock sync.Mutex
	m    map[string][]model.Instance
}

func NewInstanceMap() *InstanceMap {
	return &InstanceMap{
		m: make(map[string][]model.Instance, 16),
	}
}

func (im *InstanceMap) OverrideLatestInstances(serviceName, groupName string, instances []model.Instance) {
	im.lock.Lock()
	defer im.lock.Unlock()

	key := instanceMapKey(serviceName, groupName)
	im.m[key] = instances
}

func (im *InstanceMap) GetAndRemoveLatestInstances(serviceName, groupName string) []model.Instance {
	im.lock.Lock()
	defer im.lock.Unlock()

	key := instanceMapKey(serviceName, groupName)
	latest, exists := im.m[key]
	if !exists {
		return nil
	}

	delete(im.m, key)
	return latest
}

func instanceMapKey(serviceName, groupName string) string {
	return serviceName + ":" + groupName
}
