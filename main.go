package main

import (
	"log"
	"time"

	"github.com/wl1524520/golang_common/lib"
)

func main() {
	if err := lib.InitModule("./conf/dev/", []string{"base", "database", "redis"}); err != nil {
		log.Fatal(err)
	}
	defer lib.Destroy()

	//todo sth
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagUndefind, map[string]interface{}{"message": "todo sth"})
	time.Sleep(time.Second)
}
