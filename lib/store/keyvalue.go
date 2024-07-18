package store

import (
	"fmt"
	"kato-studio/katoengine/lib/utils"
	"kato-studio/katoengine/types"
)

type StingMapStoreInstance struct {
	Name   string
	Store  func() types.StrMap
	Get    func(key string) string
	Set    func(key string, value string) string
	Delete func(key string) types.StrMap
}

var globalStingMapStore = make(map[string]map[string]string)

func GlobalStingMap(name string) StingMapStoreInstance {
	_, found := globalStingMapStore[name]
	if !found {
		fmt.Println("WARNING: Store not found, creating new store")
		fmt.Println("Store Name: ", name)
		return CreateGlobalStingMap(name, make(map[string]string))
	}

	return StingMapStoreInstance{
		Name: name,
		Store: func() map[string]string {
			return globalStingMapStore[name]
		},
		Get: func(key string) string {
			return globalStingMapStore[name][key]
		},
		Set: func(key string, keyValue string) string {
			// globalStingMapStore[name][key] = keyValue
			return globalStingMapStore[name][key]
		},
		Delete: func(key string) map[string]string {
			delete(globalStingMapStore[name], key)
			return globalStingMapStore[name]
		},
	}
}

func CreateGlobalStingMap(name string, value map[string]string) StingMapStoreInstance {
	globalStingMapStore[name] = value
	return GlobalStingMap(name)
}

func DeleteGlobalStingMap(name string) {
	delete(globalStingMapStore, name)
}

type SmallIntStoreInstance struct {
	Store  func() types.IntMap
	Get    func(key int) string
	Set    func(key int, value string) string
	Sorted func() types.IntMap
	Delete func(key int) types.IntMap
}

func SmallIntStore() SmallIntStoreInstance {
	value := make(types.IntMap)

	return SmallIntStoreInstance{
		Store: func() types.IntMap {
			return value
		},
		Get: func(key int) string {
			return value[key]
		},
		Set: func(key int, keyValue string) string {
			value[key] = keyValue
			return value[key]
		},
		Sorted: func() map[int]string {
			return utils.SortIntMap(value)
		},
		Delete: func(key int) types.IntMap {
			delete(value, key)
			return value
		},
	}
}

type SmallStrStoreInstance struct {
	Store  func() types.StrMap
	Get    func(key string) string
	Set    func(key string, value string) string
	Sorted func() types.StrMap
	Delete func(key string) types.StrMap
}

func SmallStrStore() SmallStrStoreInstance {
	value := make(types.StrMap)

	return SmallStrStoreInstance{
		Store: func() map[string]string {
			return value
		},
		Get: func(key string) string {
			return value[key]
		},
		Set: func(key string, keyValue string) string {
			value[key] = keyValue
			return value[key]
		},
		Sorted: func() map[string]string {
			return utils.SortStrMap(value)
		},
		Delete: func(key string) map[string]string {
			delete(value, key)
			return value
		},
	}
}
