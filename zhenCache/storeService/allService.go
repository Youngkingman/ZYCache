package store

import (
	"basic/yinLog/loopqueue"
	keystruct "basic/zhenCache/innerDB/keystruct"
	"errors"
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

//logger flag
const LOG_ENABLE = true

//check expire time, may bring some unexpected effect
//set this const carfully
const CHECK_EXPIRE_TIME = 20 * time.Minute

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
var setonce sync.Once
var current_svs int = DefaultService

func getService(service int) StoreService {
	switch service {
	case MAP:
		return getServiceMap()
	case RBTREE:
		return getServiceRBtree()
	case SKIPLIST:
		return getServiceSkList()
	default:
		return getService(DefaultService)
	}
}

//GetValue GetValue
func GetValue(key keystruct.KeyStruct) (value interface{}, err error) {
	if LOG_ENABLE {
		//put current command into cache
		logitem := loopqueue.DataItem{
			Commandtype: loopqueue.GET,
			Key:         key,
			Value:       nil,
			Expire:      time.Duration(time.Now().Unix()),
			TimeStamp:   time.Now().UnixNano(),
		}
		loopqueue.LogItemPush(logitem)
	}
	return getService(current_svs).GetValue(key)
}

//SetValue SetValue
func SetValue(key keystruct.KeyStruct, value interface{}, expire time.Duration) {
	if LOG_ENABLE {
		//put current command into cache
		logitem := loopqueue.DataItem{
			Commandtype: loopqueue.SET,
			Key:         key,
			Value:       value,
			Expire:      expire,
			TimeStamp:   time.Now().UnixNano(),
		}
		loopqueue.LogItemPush(logitem)
	}
	getService(current_svs).SetValue(key, value, expire)
}

//Set Service, or it will be default
func SetStoreService(service int) error {
	setonce.Do(func() {
		current_svs = service
	})
	if current_svs == service {
		return nil
	}
	return errors.New("already set")
}
