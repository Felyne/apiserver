package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"time"

	"apiserver/config"
	"apiserver/model"
	ver "apiserver/pkg/version"
	"apiserver/router"
	"apiserver/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "c", "", "apiserver config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title Apiserver Example API
// @version 1.0
// @description apiserver demo

// @contact.name lkong
// @contact.url http://www.swagger.io/support
// @contact.email 2316082691@qq.com

// @host localhost:8080
// @BasePath /v1
func main() {
	pflag.Parse()
	if *version {
		v := ver.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
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
		//middleware.Throttle(100),
		middleware.Logging(),
		middleware.RequestId(),
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()
	go countGoroutines()
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		tlsAddr := viper.GetString("tls.addr")
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", tlsAddr)
			log.Info(http.ListenAndServeTLS(tlsAddr, cert, key, g).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())

}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}


func countGoroutines() {
	http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
		num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
		_, _ = w.Write([]byte(num))
	})
	http.ListenAndServe("localhost:6060", nil)
	log.Info("goroutine stats and pprof listen on 6060")
}