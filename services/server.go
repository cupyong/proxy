package services
import (
	"log"
)
type Service interface {
	Start() (err error)
}
type ServiceItem struct {
	S    Service
	Args interface{}
	Name string
}
var servicesMap = map[string]*ServiceItem{}

func Regist(name string, s Service) {
	servicesMap[name] = &ServiceItem{
		S:    s,
		Name: name,
	}
}

func (s *ServiceItem) Run(name string) {
	service, ok := servicesMap[name]
	if ok{
		err := service.S.Start()
		if err != nil {
			log.Fatalf("%s servcie fail, ERR: %s", name, err)
		}
	}
}
