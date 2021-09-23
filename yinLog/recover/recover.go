package recover

import (
	"basic/yinLog/logger"
	store "basic/zhenCache/storeService"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

//used only once before the service start
//not fast enough now
//files can be get by package "path"
func Recover(files []string) {
	//gurantee the time sequence is in right ordered
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	//read log items and recover them into DataItem
	dataJson := make([]logger.DataItem, 0)
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("file can't open")
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("file can't read")
		}

		items := strings.Split(string(content), "\r\n")
		for _, jsonitem := range items {
			data := logger.DataItem{}

			err = json.Unmarshal([]byte(jsonitem), &data)
			if err != nil {
				log.Println("unmarshal some data fail")
			}

			dataJson = append(dataJson, data)
		}
	}

	//recover the dataset status
	for _, data := range dataJson {
		if data.Commandtype == logger.SET {
			now := time.Now().Unix()
			if now < data.Expire {
				store.SetValue(data.Key, data.Value, time.Duration(data.Expire-now))
			}
		}
	}

}
