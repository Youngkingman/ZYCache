package main

const MAX_TEST_COUNT = 2000

func main() {
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go func() {
	// 	for i := 0; i < MAX_TEST_COUNT; i++ {
	// 		loopqueue.LogItemPush(loopqueue.DataItem{loopqueue.GET, keystruct.DefaultKey{}, i, time.Now().UnixNano()})
	// 	}
	// 	wg.Done()
	// }()

	// go func() {
	// 	for i := 0; i < MAX_TEST_COUNT; i++ {
	// 		has, item := loopqueue.LogItemPop()
	// 		fmt.Println(has, item)
	// 	}
	// 	wg.Done()
	// }()

	// wg.Wait()
	// has, item := loopqueue.LogItemPop()
	// fmt.Println(has, item, "ending")
}
