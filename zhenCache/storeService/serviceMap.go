package store

import (
	"basic/yinLog/logger"
	"errors"
	"fmt"
	"sync"
	"time"
)

type serviceMap struct {
	Store  map[string]*MemItem
	ticker *time.Ticker
	keyRL  *sync.RWMutex //for default map is concurrency unsafe
}

func getServiceMap() StoreService {
	dbonce.Do(func() {
		s := &serviceMap{
			Store:  map[string]*MemItem{},
			ticker: time.NewTicker(CHECK_EXPIRE_TIME),
			keyRL:  new(sync.RWMutex),
		}
		svs = s
		go func() {
			if LOG_ENABLE == true {
				//to start the servise
				logger.LogItemPush(logger.DataItem{
					Commandtype: logger.INITMESSAGE,
					Key:         "",
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
					needDeletedKey := []string{}
					loggerItems := make([]logger.DataItem, 0)
					s.keyRL.RLock()
					for key, mItem := range s.Store {
						if mItem.Expire < now {
							needDeletedKey = append(needDeletedKey, key)
							loggerItems = append(loggerItems, logger.DataItem{
								Commandtype: logger.SET,
								Key:         key,
								Value:       mItem.Value,
								Expire:      mItem.Expire,
								TimeStamp:   now,
							})
						}
					}
					s.keyRL.RUnlock()

					go logger.RdbLog(loggerItems)

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

func (svs *serviceMap) GetValue(key string) (value interface{}, err error) {
	defer svs.keyRL.RUnlock()
	svs.keyRL.RLock()
	for k, v := range svs.Store {
		fmt.Println(k, v)
	}
	if mitem, ok := svs.Store[key]; ok {
		if mitem.Expire >= time.Now().Unix() {
			mitem.Expire = time.Now().Add(mitem.duration).Unix()
			return mitem.Value, nil
		}
		return nil, errors.New("expire")
	}
	return nil, errors.New("no value")
}

func (svs *serviceMap) SetValue(key string, value interface{}, expire time.Duration) {
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

// func (svs *serviceMap) GetRange(keyL string, keyH keystruct.KeyStruct) (values []interface{}, err error) {
// 	return nil, errors.New("no support for unordered map")
// }
