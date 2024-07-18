package store

type ComponentStoreInstance struct {
	Store func() map[string][]byte
	Get func(key string) []byte
	Set func(key string, value []byte) []byte
	SafeSet func(key string, value []byte) []byte
	Delete func(key string) map[string][]byte
}

var globalByteMapStore = make(map[string][]byte)

func GlobalByteMap() ComponentStoreInstance {
	return ComponentStoreInstance{
		Store: func() map[string][]byte {
			return globalByteMapStore
		},
		Get: func(key string) []byte {
			return globalByteMapStore[key]
		},
		Set: func(key string, keyValue []byte) []byte {
			globalByteMapStore[key] = keyValue
			return globalByteMapStore[key]
		},
		SafeSet: func(key string, keyValue []byte) []byte {
			_, exists := globalByteMapStore[key]
			if !exists {
				globalByteMapStore[key] = keyValue
				return globalByteMapStore[key]
			}
			return nil
		},
		Delete: func(key string) map[string][]byte {
			delete(globalByteMapStore, key)
			return globalByteMapStore
		},
	}
}