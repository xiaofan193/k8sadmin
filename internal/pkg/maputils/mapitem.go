package maputils

import "github.com/xiaofan193/k8sadmin/internal/types"

func ToMap(items []types.ListMapItem) map[string]string {
	dataMap := make(map[string]string)
	for _, item := range items {
		dataMap[item.Key] = item.Value
	}
	return dataMap
}

func ToList(data map[string]string) []types.ListMapItem {
	list := make([]types.ListMapItem, 0)
	for k, v := range data {
		list = append(list, types.ListMapItem{
			Key:   k,
			Value: v,
		})
	}
	return list
}
func ToListWithMapByte(data map[string][]byte) []types.ListMapItem {
	list := make([]types.ListMapItem, 0)
	for k, v := range data {
		list = append(list, types.ListMapItem{
			Key:   k,
			Value: string(v),
		})
	}
	return list
}
