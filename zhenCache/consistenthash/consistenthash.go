package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type ConsisteHashMap struct {
	hash     Hash           //hash function
	replicas int            //number of virtual peers of one machine
	keys     []int          //sorted
	innerMap map[int]string //certain hash to real machine
}

func New(replicas int, fn Hash) *ConsisteHashMap {
	m := &ConsisteHashMap{
		replicas: replicas,
		hash:     fn,
		innerMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (chm *ConsisteHashMap) Add(realMachineKeys ...string) {
	for _, key := range realMachineKeys {
		for i := 0; i < chm.replicas; i++ {
			hash := int(chm.hash([]byte(strconv.Itoa(i) + key)))
			chm.keys = append(chm.keys, hash)
			chm.innerMap[hash] = key
		}
	}
}

func (chm *ConsisteHashMap) Get(key string) string {
	if len(chm.keys) == 0 {
		return ""
	}

	hash := int(chm.hash([]byte(key)))
	idx := sort.Search(len(chm.keys), func(i int) bool {
		return chm.keys[i] >= hash
	})

	return chm.innerMap[chm.keys[idx%len(chm.keys)]]
}
