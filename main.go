package main

import (
	"./services"
	"flag"
)

var (
	service *services.ServiceItem
)
func main() {
	t := *flag.String("t", "agent", "")
	flag.Parse()
	switch t {
	case "agent":
		services.Regist(t, services.NewAgent())
	}
	service.Run(t)
}
