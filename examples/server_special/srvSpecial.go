package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	modbus "github.com/aloncn/gomodbus"
)

func main() {
	srv := modbus.NewTCPServerSpecial()
	if err := srv.AddRemoteServer("192.168.199.148:2404"); err != nil {
		panic(err)
	}
	srv.LogMode(true)
	srv.AddNodes(
		modbus.NewNodeRegister(
			1,
			0, 10, 0, 10,
			0, 10, 0, 10),
		modbus.NewNodeRegister(
			2,
			0, 10, 0, 10,
			0, 10, 0, 10),
		modbus.NewNodeRegister(
			3,
			0, 10, 0, 10,
			0, 10, 0, 10))

	srv.SetOnConnectHandler(func(c *modbus.TCPServerSpecial) error {
		_, err := c.UnderlyingConn().Write([]byte("hello world"))
		return err
	})

	srv.SetConnectionLostHandler(func(c *modbus.TCPServerSpecial) {
		log.Println("connect lost")
	})

	srv.SetKeepAlive(true, time.Second*20, func(c *modbus.TCPServerSpecial) {
		c.UnderlyingConn().Write([]byte("keep alive"))
	})

	if err := srv.Start(); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":6060", nil); err != nil {
		panic(err)
	}
}
