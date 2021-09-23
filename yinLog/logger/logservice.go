package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type logTimer struct {
	ticker *time.Ticker
}

//every hour create a new log file
const LOG_SPLIT_TIME = time.Hour

func startLogServe() {
	lt := logTimer{ticker: time.NewTicker(LOG_SPLIT_TIME)}
	oname := fmt.Sprintf("zy-log-out-%s.zyl",
		time.Now().Format("2006-01-02 15:04:05.000000"))
	ofile, err := os.Create(oname)
	defer ofile.Close()
	if err != nil {
		panic("can't create log file, check your authority")
	}
	for {
		//update the file and create a new one
		select {
		case <-lt.ticker.C:
			ofile.Close()
			oname = fmt.Sprintf("zy-log-out-%s.zyl",
				time.Now().Format("2006-01-02 15:04:05.000000"))
			ofile, err = os.Create(oname)
			if err != nil {
				panic("can't create log file, check your authority")
			}
		default:
		}
		has, item := LogItemPop()
		if has {
			itemJson, err := json.Marshal(item)
			itemJson = append(itemJson, []byte("\r\n")...)
			if err != nil {
				itemJson = []byte(fmt.Sprintf("fail insert marshal item at -%s \r\n",
					time.Now().Format("2006-01-02 15:04:05.000000")))
			}
			if _, err := ofile.Write(itemJson); err != nil {
				panic(err)
			}
		}
	}
}
