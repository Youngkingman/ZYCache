package store

import (
	"basic/yinLog/logger"
	keystruct "basic/zhenCache/innerDB/keystruct"
	"errors"
	"sync"
	"time"
)

type serviceMap struct {
	Store  map[keystruct.KeyStruct]*MemItem
	ticker *time.Ticker
	keyRL  *sync.RWMutex //for default map is concurrency unsafe
}

func getServiceMap() StoreService {
	dbonce.Do(func() {
		s := &serviceMap{
			Store:  map[keystruct.KeyStruct]*MemItem{},
			ticker: time.NewTicker(CHECK_EXPIRE_TIME),
			keyRL:  new(sync.RWMutex),
		}
		svs = s
		go func() {
			if LOG_ENABLE == true {
				//to start the servise
				logger.LogItemPush(logger.DataItem{
					Commandtype: logger.INITMESSAGE,
					Key:         keystruct.DefaultKey{},
					Value:       nil,
					Expire:      0,
					TimeStamp:   time.Now().Unix(),
				})
				defer logger.ShutLog()
			}
			for {
				select {
				case <-s.ticker.C:
					now := time.Now().Unix()
					needDeletedKey := []keystruct.KeyStruct{}
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

func (svs *serviceMap) GetValue(key keystruct.KeyStruct) (value interface{}, err error) {
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

func (svs *serviceMap) SetValue(key keystruct.KeyStruct, value interface{}, expire time.Duration) {
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

func (svs *serviceMap) GetRange(keyL keystruct.KeyStruct, keyH keystruct.KeyStruct) (values []interface{}, err error) {
	return nil, errors.New("no support for unordered map")
}
