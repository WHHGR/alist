package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alist-org/alist/v3/cmd/args"
	_ "github.com/alist-org/alist/v3/drivers"
	"github.com/alist-org/alist/v3/internal/bootstrap"
	"github.com/alist-org/alist/v3/internal/bootstrap/data"
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	flag.StringVar(&args.Config, "conf", "data/config.json", "config file")
	flag.BoolVar(&args.Debug, "debug", false, "start with debug mode")
	flag.BoolVar(&args.Version, "version", false, "print version info")
	flag.BoolVar(&args.Password, "password", false, "print current password")
	flag.BoolVar(&args.NoPrefix, "no-prefix", false, "disable env prefix")
	flag.BoolVar(&args.Dev, "dev", false, "start with dev mode")
	flag.Parse()
}

func Init() {
	if args.Version {
		fmt.Printf("Built At: %s\nGo Version: %s\nAuthor: %s\nCommit ID: %s\nVersion: %s\nWebVersion: %s\n",
			conf.BuiltAt, conf.GoVersion, conf.GitAuthor, conf.GitCommit, conf.Version, conf.WebVersion)
		os.Exit(0)
	}
	bootstrap.InitConfig()
	bootstrap.Log()
	bootstrap.InitDB()
	data.InitData()
	bootstrap.InitAria2()
}
func main() {
	Init()
	if !args.Debug && !args.Dev {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
	server.Init(r)
	base := fmt.Sprintf("%s:%d", conf.Conf.Address, conf.Conf.Port)
	log.Infof("start server @ %s", base)
	var err error
	if conf.Conf.Scheme.Https {
		err = r.RunTLS(base, conf.Conf.Scheme.CertFile, conf.Conf.Scheme.KeyFile)
	} else {
		err = r.Run(base)
	}
	if err != nil {
		log.Errorf("failed to start: %s", err.Error())
	}
}
