package store

import (
	keystruct "basic/util/KeyStruct"
	"sync"
	"time"
)

const (
	MAP = iota
	RBTREE
	SKIPLIST
)

//DefaultDuration the default duration
const DefaultDuration = 15 * time.Second

//inner struct
const DefaultService = MAP

//sklist level
const SKListLevel = 16

//StoreService memservice
type StoreService interface {
	GetValue(key keystruct.KeyStruct) (value interface{}, err error)
	SetValue(key keystruct.KeyStruct, value interface{}, expire time.Duration)
	GetRange(keyL keystruct.KeyStruct, keyH keystruct.KeyStruct) (values []interface{}, err error)
}

type MemItem struct {
	Value    interface{}
	Expire   int64
	duration time.Duration
}

var svs StoreService
var dbonce sync.Once

func getService() StoreService {
	switch DefaultService {
	case MAP:
		return getServiceMap()
	case RBTREE:
		return getServiceRBtree()
	case SKIPLIST:
		return getServiceSkList()
	}
	return nil
}

//GetValue GetValue
func GetValue(key keystruct.KeyStruct) (value interface{}, err error) {
	return getService().GetValue(key)
}

//SetValue SetValue
func SetValue(key keystruct.KeyStruct, value interface{}, expire time.Duration) {
	getService().SetValue(key, value, expire)
}

//GetValueDefault  with default
func GetValueDefault(key keystruct.KeyStruct, getV func() interface{}) interface{} {
	if val, err := GetValue(key); err == nil {
		return val
	}

	val := getV()
	SetValue(key, val, time.Minute*3)
	return val
}
