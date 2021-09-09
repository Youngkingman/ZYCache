package hashmap

import (
	"errors"
	"sync"
	"time"
)

//DefaultDuration the default duration
const DefaultDuration = 15 * time.Second

//IService memservice
type IService interface {
	GetValue(key string) (value interface{}, err error)
	SetValue(key string, value interface{}, expire time.Duration)
}

//MemItem the meme item
type MemItem struct {
	Value    interface{}
	Expire   int64
	duration time.Duration
}

type service struct {
	Store  map[string]*MemItem
	ticker *time.Ticker
	keyRL  *sync.RWMutex
}

var svs IService
var dbonce sync.Once

func getService() IService {
	dbonce.Do(func() {
		s := &service{
			Store:  map[string]*MemItem{},
			ticker: time.NewTicker(time.Minute * 10),
			keyRL:  new(sync.RWMutex),
		}
		svs = s
		go func() {
			for {
				select {
				case <-s.ticker.C:
					now := time.Now().Unix()
					needDeletedKey := []string{}
					s.keyRL.RLock()
					for key, mItem := range s.Store {
						if mItem.Expire < now {
							needDeletedKey = append(needDeletedKey, key)
						}
					}
					s.keyRL.RUnlock()

					s.keyRL.Lock()
					for _, key := range needDeletedKey {
						delete(s.Store, key)
					}
					s.keyRL.Unlock()
				}
			}
		}()
	})
	return svs
}

func (svs *service) GetValue(key string) (value interface{}, err error) {
	defer svs.keyRL.RUnlock()
	svs.keyRL.RLock()
	if mitem, ok := svs.Store[key]; ok {
		if mitem.Expire >= time.Now().Unix() {
			mitem.Expire = time.Now().Add(mitem.duration).Unix()
			return mitem.Value, nil
		}
		return nil, errors.New("expire")
	}
	return nil, errors.New("no value")
}

func (svs *service) SetValue(key string, value interface{}, expire time.Duration) {
	defer svs.keyRL.Unlock()
	svs.keyRL.Lock()
	if mitem, ok := svs.Store[key]; ok {
		mitem.Expire = time.Now().Add(expire).Unix()
		mitem.Value = value
		return
	}

	m := MemItem{
		value,
		time.Now().Add(expire).Unix(),
		expire,
	}
	svs.Store[key] = &m
}

//GetValue GetValue
func GetValue(key string) (value interface{}, err error) {
	return getService().GetValue(key)
}

//SetValue SetValue
func SetValue(key string, value interface{}, expire time.Duration) {
	getService().SetValue(key, value, expire)
}

//GetValueDefault  with default
func GetValueDefault(key string, getV func() interface{}) interface{} {
	if val, err := GetValue(key); err == nil {
		return val
	}

	val := getV()
	SetValue(key, val, time.Minute*3)
	return val
}
