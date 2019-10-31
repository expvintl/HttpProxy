package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"
)

func main() {
	server,err:=net.Listen("tcp",":8080")
	fmt.Println("HTTP Proxy Running on :8080....")
	if err!=nil{
		panic(err)
	}
	for {
		//Listen
		cli, err := server.Accept()
		if err != nil {
			fmt.Println("[ERROR] Socket:", err)
			return
		}
		go handler(cli)
	}
}
func handler(cli net.Conn){
	if cli==nil{return}
	var b [1024]byte
	n, _ := cli.Read(b[:])
	//Method Host Address Port
	var method,host,addr,port string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	port=host[strings.Index(host,":")+1:]
	//if port is 443
		if port == "443" {
			addr = host
		} else{
		//Determine if htt(p://) is included
			if strings.Contains(host,"p://"){
				if url2,err:=url.Parse(host);err!=nil{
					fmt.Println("[ERROR]",err)
					return
				}else{
				//if port is empty default:80
					if url2.Port()=="" {
						addr = url2.Host + ":80"
					}else{
						addr = url2.Host
					}
				}
			}else {
				//IP Handle
					addr = host
				}
			}
			start:=time.Now()
			client, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println("[\033[31mERROR\033[0m] Client:", err)
				return
			}
			if method == "CONNECT" {
				fmt.Fprint(cli, "HTTP/1.1 200 Connection Established\r\n\r\n")
			} else {
				client.Write(b[:n])
			}
			fmt.Printf("[\033[32mINFO\033[0m] Method:%s,Host:%s \033[32m%s\033[0m\n",method,host,time.Since(start))
			go io.Copy(client, cli)
			io.Copy(cli, client)
			cli.Close()
		}