package store

import (
	keystruct "basic/util/KeyStruct"
	"basic/util/rbtree"
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
			ticker: time.NewTicker(time.Minute * 60),
		}
		svs = s
		go func() {
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

func (svs *serviceRBtree) GetValue(key keystruct.KeyStruct) (value interface{}, err error) {
	if mitem, ok := svs.Store.Search(key); ok {
		if mitem.(*MemItem).Expire >= time.Now().Unix() {
			mitem.(*MemItem).Expire = time.Now().Add(mitem.(*MemItem).duration).Unix()
			return mitem.(*MemItem).Value, nil
		}
	}
	return nil, errors.New("expire")
}

func (svs *serviceRBtree) SetValue(key keystruct.KeyStruct, value interface{}, expire time.Duration) {
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

func (svs *serviceRBtree) GetRange(keyL keystruct.KeyStruct, keyH keystruct.KeyStruct) (values []interface{}, err error) {
	//TODO
	return
}
