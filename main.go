package main

import (
	"cnng/conf"
	"cnng/controller"
	"cnng/model"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
)


var confPath string

func init() {
	flag.StringVar(&confPath, "c", "./conf/default.yaml", "配置文件")
	flag.Parse()
}

func main() {
	var errChan = make(chan error)
	err := conf.UnmarshalYamlFile(confPath)
	if err != nil {
		log.Fatal("load main config fail:", err.Error())
		return
	}
	var m model.Manager
	go startHttp(&m, errChan)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	logrus.Infof("exit:%s", <-errChan)
}

func startHttp(m *model.Manager, errChan chan error) {
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	n := negroni.New(negroni.NewRecovery())
	n.Use(c)
	n.UseHandler(router)
	controller.NewHandler(m, router)
	log.Printf("listen on http: %d", conf.Conf.App.Port)
	errChan <- netHttp.ListenAndServe(fmt.Sprintf(":%d", conf.Conf.App.Port), n)
}
