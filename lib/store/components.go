package store

type ComponentStoreInstance struct {
	Store func() map[string][]byte
	Get func(key string) []byte
	Set func(key string, value []byte) []byte
	SafeSet func(key string, value []byte) []byte
	Delete func(key string) map[string][]byte
}

var global_byte_map_store = make(map[string][]byte)

func GlobalByteMap() ComponentStoreInstance {
	return ComponentStoreInstance{
		Store: func() map[string][]byte {
			return global_byte_map_store
		},
		Get: func(key string) []byte {
			return global_byte_map_store[key]
		},
		Set: func(key string, keyValue []byte) []byte {
			global_byte_map_store[key] = keyValue
			return global_byte_map_store[key]
		},
		SafeSet: func(key string, keyValue []byte) []byte {
			_, exists := global_byte_map_store[key]
			if !exists {
				global_byte_map_store[key] = keyValue
				return global_byte_map_store[key]
			}
			return nil
		},
		Delete: func(key string) map[string][]byte {
			delete(global_byte_map_store, key)
			return global_byte_map_store
		},
	}
}