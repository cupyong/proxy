package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Agent struct {
	Port string `yaml:"port"`
}
func (a *Agent) Start() error {
	l, err := net.Listen("tcp", ":"+a.Port)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(a.Port,"启动代理")
	for {
		client, err := l.Accept()
		fmt.Println(client)
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
	return  nil
}

func NewAgent() Service {
	YamlConf:=new(Agent)
	yamlFile, err := ioutil.ReadFile("./agent.yaml")
	if err != nil {
		log.Printf("yaml open err is : \n %v \n", err)
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, YamlConf)
	return &Agent{
		Port: YamlConf.Port,
	}
}
func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(method, host, address)
	if hostPortURL.Opaque == "443" { //https访问
		address = hostPortURL.Scheme + ":443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n")
	} else {
		server.Write(b[:n])
	}
	//进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
}
