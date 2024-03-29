package ddd_mongodb

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filterItem struct {
	parent *filterItem
	name   string
	value  interface{}
	items  []*filterItem
}

func newFilterItem(parent *filterItem, name string) *filterItem {
	n := name
	if n == "id" {
		n = "_id"
	}
	return &filterItem{
		name:   n,
		parent: parent,
		value:  nil,
		items:  make([]*filterItem, 0),
	}
}

func (i *filterItem) addChildItem(name string, value interface{}) *filterItem {
	newItem := newFilterItem(i, name)
	newItem.value = value
	i.items = append(i.items, newItem)
	return newItem
}

func (i *filterItem) getAndItem() {
}

func (i *filterItem) getValues(data map[string]interface{}) {
	if len(i.items) != 0 {
		array := make([]interface{}, len(i.items))
		for i, v := range i.items {
			item := ddd_utils.NewMap()
			item[v.name] = v.value
			if len(v.items) > 0 {
				m := ddd_utils.NewMap()
				v.getValues(m)
				item[v.name] = m[v.name]
			}
			array[i] = item
		}
		data[i.name] = array
	} else if i.value != nil {
		data[i.name] = i.value
	}
}

func (i *filterItem) setValue(name string, value interface{}) {
	i.name = name
	i.value = value
}

type MongoProcess struct {
	item    *filterItem
	current *filterItem
}

func NewMongoProcess() *MongoProcess {
	m := &MongoProcess{
		item: newFilterItem(nil, "$and"),
	}
	m.init()
	return m
}

func (m *MongoProcess) init() {
	m.current = m.item
}

func (m *MongoProcess) GetFilter(tenantId string) map[string]interface{} {
	data := make(map[string]interface{})
	m.item.getValues(data)
	m1, ok := data[""]
	if ok {
		d1 := m1.(map[string]interface{})
		d1[TenantIdField] = tenantId
	} else if len(data) == 0 {
		data[TenantIdField] = tenantId
	} else {
		m1, ok := data["$and"]
		d1, ok := m1.(map[string]interface{})
		if ok {
			d1[TenantIdField] = tenantId
		}
		d2, ok := m1.([]interface{})
		if ok {
			item := ddd_utils.NewMap()
			item[TenantIdField] = tenantId
			d2 := append(d2, item)
			data["$and"] = d2
		}
	}
	return data
}

func (m *MongoProcess) OnAndItem() {
	m.current.name = "$and"
}

func (m *MongoProcess) OnAndStart() {
	m.current = m.current.addChildItem("$and", nil)
}

func (m *MongoProcess) OnAndEnd() {
	m.current = m.current.parent
}

func (m *MongoProcess) OnOrItem() {
	m.current.name = "$or"
}

func (m *MongoProcess) OnOrStart() {
	m.current = m.current.addChildItem("$or", nil)
}

func (m *MongoProcess) OnOrEnd() {
	m.current = m.current.parent
}

func (m *MongoProcess) OnEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, rsql.GetValue(rValue))
}

func (m *MongoProcess) OnNotEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$ne", rsql.GetValue(rValue)}})
}

func (m *MongoProcess) OnLike(name string, value interface{}, rValue rsql.Value) {
	pattern := fmt.Sprintf("%s", rsql.GetValue(rValue))
	m.current.addChildItem(name, primitive.Regex{Pattern: pattern, Options: "im"})
}

func (m *MongoProcess) OnNotLike(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", rsql.GetValue(rValue)}})
}

func (m *MongoProcess) OnGreaterThan(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$gt", rsql.GetValue(rValue)}})
}

func (m *MongoProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$gte", rsql.GetValue(rValue)}})
}

func (m *MongoProcess) OnLessThan(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lt", rsql.GetValue(rValue)}})
}

func (m *MongoProcess) OnLessThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	m.current.addChildItem(name, bson.D{{"$lte", rsql.GetValue(rValue)}})
}

func (m *MongoProcess) OnIn(name string, value interface{}, rValue rsql.Value) {
	listValue, _ := rValue.(rsql.ListValue)
	values := rsql.GetValueList(listValue)
	m.current.addChildItem(name, bson.M{"$in": values})
}

func (m *MongoProcess) OnNotIn(name string, value interface{}, rValue rsql.Value) {
	listValue, _ := rValue.(rsql.ListValue)
	values := rsql.GetValueList(listValue)
	m.current.addChildItem(name, bson.M{"$nin": values})
}
