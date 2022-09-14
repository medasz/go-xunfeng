package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

var Cfg *Config

type Config struct {
	Mode      string `ini:"mode"`
	SecretKey string
	Server    Server  `ini:"server"`
	User      User    `ini:"user"`
	MongoDb   MongoDb `ini:"mongodb"`
}

type Server struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

func (s Server) HostPort() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type User struct {
	Account  string `ini:"account"`
	Password string `ini:"password"`
}

type MongoDb struct {
	Host       string `ini:"host"`
	Port       int    `ini:"port"`
	DbUsername string `ini:"dbusername"`
	DbPassword string `ini:"dbpassword"`
	DbName     string `ini:"dbname"`
}

func init() {
	Cfg = new(Config)
	err := ini.MapTo(&Cfg, "config.ini")
	if err != nil {
		panic(err)
	}
	Cfg.SecretKey = os.Getenv("SECRET_KEY")
	if Cfg.SecretKey == "" {

	}
}
