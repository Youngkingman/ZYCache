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

//every hour create a new log file
const LOG_SPLIT_TIME = time.Hour

var lt logTimer

func startLogServe() {
	lt = logTimer{
		ticker: time.NewTicker(LOG_SPLIT_TIME),
		stop:   make(chan struct{}),
	}
	defer close(lt.stop)

	oname := fmt.Sprintf("zy_log_out-%d.zyl",
		time.Now().Unix())
	ofile, err := os.Create("../logbin/" + oname)
	defer ofile.Close()
	if err != nil {
		panic("can't create log file, check your authority")
	}
	for {
		//update the file and create a new one
		select {
		case <-lt.ticker.C:
			ofile.Close()
			oname = fmt.Sprintf("zy_log_out_%d.zyl",
				time.Now().Unix())
			ofile, err = os.Create("../logbin/" + oname)
			if err != nil {
				panic("can't create log file, check your authority")
			}
		case <-lt.stop:
			for {
				has, item := LogItemPop()
				if !has {
					break
				}
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
			shutdown <- struct{}{}
			break
		default:
		}
		has, item := LogItemPop()
		if has {
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
	}
}
