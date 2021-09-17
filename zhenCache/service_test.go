package zhenCache

import (
	keystruct "basic/zhenCache/innerDB/KeyStruct"
	"testing"
)

type TestKey struct {
	keystruct.DefaultKey
	key string
}

func (key TestKey) CompareBiggerThan(other keystruct.KeyStruct) bool {
	return key.key > other.KeyString()
}

func (key TestKey) KeyString() string {
	return key.key
}

func Test_Service_Access(t *testing.T) {

}
