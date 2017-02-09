// Please use library.
package main

import (
	_ "github.com/ossrs/go-oryx-lib/http"
	_ "github.com/ossrs/go-oryx-lib/https"
	_ "github.com/ossrs/go-oryx-lib/json"
	_ "github.com/ossrs/go-oryx-lib/kxps"
	ol "github.com/ossrs/go-oryx-lib/logger"
	_ "github.com/ossrs/go-oryx-lib/options"
	"github.com/SnailTowardThesun/go-oryx-lib/rtmp"
	"time"
)

func main() {
	rtmpUrl := "rtmp://127.0.0.1:1935/live/livestream"
	rtmpClient, err := rtmp.NewSimpleRtmpClient(rtmpUrl)
	if err != nil {
		ol.E(nil, "create simple rtmp client failed. err is", err)
		return
	}
	defer rtmpClient.Close()
	
	time.Sleep(time.Duration(30) * time.Second)
	return
}
