package store

import (
	"basic/yinLog/logger"
	keystruct "basic/zhenCache/innerDB/keystruct"
	skiplist "basic/zhenCache/innerDB/skipList"
	"errors"
	"time"
)

type serviceSkList struct {
	Store  skiplist.SkipList
	ticker *time.Ticker
}

func getServiceSkList() StoreService {
	dbonce.Do(func() {
		s := &serviceSkList{
			Store:  skiplist.New(SKListLevel),
			ticker: time.NewTicker(CHECK_EXPIRE_TIME),
		}
		svs = s
		go func() {
			if LOG_ENABLE == true {
				//to start the servise
				logger.LogItemPush(logger.DataItem{
					Commandtype: logger.INITMESSAGE,
					Key:         keystruct.DefaultKey{},
					Value:       nil,
					Expire:      time.Microsecond,
					TimeStamp:   time.Now().Unix(),
				})
				defer logger.ShutLog()
			}
			for {
				select {
				case <-s.ticker.C:
					now := time.Now().Unix()
					needDeletedKey := s.Store.Range(func(i interface{}) bool {
						return i.(*MemItem).Expire < now
					})

					for _, key := range needDeletedKey {
						s.Store.Delete(key)
					}
				}
			}
		}()
	})
	return svs
}

func (svs *serviceSkList) GetValue(key keystruct.KeyStruct) (value interface{}, err error) {
	if mitem, ok := svs.Store.Search(key); ok {
		if mitem.(*MemItem).Expire >= time.Now().Unix() {
			mitem.(*MemItem).Expire = time.Now().Add(mitem.(*MemItem).duration).Unix()
			return mitem.(*MemItem).Value, nil
		}
	}
	return nil, errors.New("expire")
}

func (svs *serviceSkList) SetValue(key keystruct.KeyStruct, value interface{}, expire time.Duration) {
	if mitem, ok := svs.Store.Search(key); ok {
		mitem.(*MemItem).Expire = time.Now().Add(expire).Unix()
		mitem.(*MemItem).Value = value
		return
	}

	m := MemItem{
		value,
		time.Now().Add(expire).Unix(),
		expire,
	}
	svs.Store.InsertElement(key, &m)
}

func (svs *serviceSkList) GetRange(keyL keystruct.KeyStruct, keyH keystruct.KeyStruct) (values []interface{}, err error) {
	//TODO
	return
}
