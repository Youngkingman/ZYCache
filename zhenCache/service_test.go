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

// func BenchmarkSet(b *testing.B) {
//     addr := "localhost:8001"
//     for i := 0; i < b.N; i++ {
//         url := fmt.Sprintf("http://%s/key/%d", addr, i)
//         post(url, []byte(url), b)
//     }
// }

// func BenchmarkGet(b *testing.B) {
//     addr := "localhost:8001"
//     for i := 0; i < b.N; i++ {
//         url := fmt.Sprintf("http://%s/key/%d", addr, i)
//         get(url, b)
//     }
// }

// func BenchmarkPSet(b *testing.B) {
//     addr := "localhost:8001"
//     b.RunParallel(func(pb *testing.PB) {
//         i := 0
//         for pb.Next() {
//             i++
//             url := fmt.Sprintf("http://%s/key/%d", addr, i)
//             post(url, []byte(url), b)
//         }
//     })
// }

// func BenchmarkPGet(b *testing.B) {
//     addr := "localhost:8001"
//     b.RunParallel(func(pb *testing.PB) {
//         i := 0
//         for pb.Next() {
//             i++
//             url := fmt.Sprintf("http://%s/key/%d", addr, i)
//             get(url, b)
//         }
//     })
// }

// func post(url string, data []byte, l testing.TB) {
//     ret, err := httpClient.Post(url, "", bytes.NewReader(data))
//     if err != nil {
//         l.Fatal(err)
//     }
//     defer ret.Body.Close()
//     io.Copy(ioutil.Discard, ret.Body)
// }

// func get(url string, l testing.TB) {
//     ret, err := httpClient.Get(url)
//     if err != nil {
//         l.Fatal(err)
//     }
//     defer ret.Body.Close()
//     io.Copy(ioutil.Discard, ret.Body)
// }
