package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type logTimer struct {
	ticker *time.Ticker
	stop   chan struct{}
}

//bring IO presure every LOG_SPLIT_TIME, be carefully set
const LOG_SPLIT_TIME = 15 * time.Minute

//the max item cout of a single log
const LOG_SPLIT_COUNT = 1024 * 1024

var lt logTimer

func startLogAppendServe() {
	lt = logTimer{
		ticker: time.NewTicker(LOG_SPLIT_TIME),
		stop:   make(chan struct{}),
	}
	current_count := 0
	log_id := 0

	defer close(lt.stop)

	oname := fmt.Sprintf("zy_log_out_%dt_%d.zyl", log_id,
		time.Now().UnixNano())
	ofile, err := os.Create("../logbin/" + oname)
	defer ofile.Close()

	if err != nil {
		panic("can't create log file, check your authority")
	}
	for {
		//update the file and create a new one
		select {
		case <-lt.ticker.C:
			ofile = splitFile(ofile, log_id)
			log_id++
		case <-lt.stop:
			for {
				has, item := LogItemPop()
				if !has {
					break
				}
				jsonWrite(ofile, item)
			}
			shutdown <- struct{}{}
			break
		default:
		}
		has, item := LogItemPop()
		if has {
			current_count++
			jsonWrite(ofile, item)
		}
		if current_count > LOG_SPLIT_COUNT {
			ofile = splitFile(ofile, log_id)
			log_id++
			current_count = 0
		}
	}
}

func splitFile(ofile *os.File, log_id int) *os.File {
	ofile.Close()
	oname := fmt.Sprintf("zy_log_out_%dt_%d.zyl", log_id,
		time.Now().UnixNano())
	newofile, err := os.Create("../logbin/" + oname)
	if err != nil {
		panic("can't create log file, check your authority")
	}
	return newofile
}

func jsonWrite(ofile *os.File, item interface{}) {
	itemJson, err := json.Marshal(item)
	itemJson = append(itemJson, []byte("\r\n")...)
	if err != nil {
		itemJson = []byte(fmt.Sprintf("fail insert marshal item at -%d \r\n",
			time.Now().Unix()))
	}
	if _, err := ofile.Write(itemJson); err != nil {
		panic(err)
	}
}

//currently can be used for map, but not RBTree & SkipList
func RdbLog(items []DataItem) {
	oname := fmt.Sprintf("rdb_copy_%d.zyl",
		time.Now().UnixNano())
	ofile, err := os.Create("../logbin/" + oname)
	if err != nil {
		panic(err)
	}
	defer ofile.Close()
	for _, item := range items {
		jsonWrite(ofile, item)
	}
}
