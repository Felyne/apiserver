package main

import (
	"apiserver/config"
	"apiserver/model"
	v "apiserver/pkg/version"
	"apiserver/router"
	"apiserver/router/middleware"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *version {
		value := v.Get()
		info, err := json.MarshalIndent(&value, "", " ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(info))
		return
	}

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	model.DB.Init()
	defer model.DB.Close()

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	router.Load(
		g,
		middleware.RequestId(), //全局中间件
		middleware.Logging(),
	)

	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		tlsAddr := viper.GetString("tls.addr")
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", tlsAddr)
			log.Info(http.ListenAndServeTLS(tlsAddr, cert, key, g).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s")
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// make sure the router is working.
func pingServer() error {
	count := viper.GetInt("max_ping_count")
	url := viper.GetString("url")

	for i := 0; i < count; i++ {
		resp, err := http.Get(url + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("cannot connect to the router")
}
