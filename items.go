package launchbar

import (
	"encoding/json"
	"fmt"
)

type Items []*Item

func NewItems() *Items {
	return &Items{}
}

func (items *Items) Add(i ...*Item) *Items {
	for _, item := range i {
		*items = append(*items, item)
	}
	return items
}

func (i *Items) setItems(items []*item) {
	for _, item := range items {
		i.Add(newItem(item))
	}
}

func (items *Items) getItems() []*item {
	a := make([]*item, len(*items))
	for i, item := range *items {
		a[i] = item.item
	}
	return a
}

type itemsByOrder Items

func (o itemsByOrder) Len() int           { return len(o) }
func (o itemsByOrder) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o itemsByOrder) Less(i, j int) bool { return o[i].item.Order < o[j].item.Order }

func (items *Items) Compile() string {
	if len(*items) == 0 {
		return ""
	}

	b, err := json.Marshal(items.getItems())
	if err != nil {
		return fmt.Sprintf(`[{"title": "%v","subtitle":"error"}]`, err)
	}
	return string(b)
}
