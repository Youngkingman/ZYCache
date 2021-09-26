package store

import (
	"basic/yinLog/logger"
	"basic/zhenCache/innerDB/rbtree"
	"errors"
	"time"
)

type serviceRBtree struct {
	Store  rbtree.RBTree
	ticker *time.Ticker
}

func getServiceRBtree() StoreService {
	dbonce.Do(func() {
		s := &serviceRBtree{
			Store:  rbtree.New(),
			ticker: time.NewTicker(CHECK_EXPIRE_TIME),
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

func (svs *serviceRBtree) GetValue(key string) (value interface{}, err error) {
	if mitem, ok := svs.Store.Search(key); ok {
		if mitem.(*MemItem).Expire >= time.Now().Unix() {
			mitem.(*MemItem).Expire = time.Now().Add(mitem.(*MemItem).duration).Unix()
			return mitem.(*MemItem).Value, nil
		}
	}
	return nil, errors.New("expire")
}

func (svs *serviceRBtree) SetValue(key string, value interface{}, expire time.Duration) {
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

// func (svs *serviceRBtree) GetRange(keyL string, keyH keystruct.KeyStruct) (values []interface{}, err error) {
// 	//TODO
// 	return
// }
