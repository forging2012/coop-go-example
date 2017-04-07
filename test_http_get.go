package main

import (
	"fmt"
	"coop-go"
	"kendynet-go/kendynet"
	"kendynet-go/packet"	
	"net/http"
	"os"
)


func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage ./test_http_get ip:port\n")
		return
	}

	service := os.Args[1]

	tcpServer,err := kendynet.NewTcpServer(service,packet.NewRawUnPacker(1024))

	if nil != err {
		fmt.Printf(err.Error())
		return
	}

	var p *coop.CoopScheduler

	p = coop.NewCoopScheduler(func (e interface{}){

		session := e.(*kendynet.StreamConn)

		var resp *http.Response
		var err error

		//执行同步调用
		p.Call(func () {
			resp,err = http.Get("http://www.01happy.com/demo/accept.php?id=1")
		})

		var buff []byte
		var lens int

		if nil != err {
			lens = len(err.Error())
			buff = make([]byte,lens+1)
			copy(buff[:],err.Error()[:lens])
			buff[lens] = byte('\n')
		}else{
			lens = len(resp.Status)
			buff = make([]byte,lens+1)			
			copy(buff[:],resp.Status[:lens])
			buff[lens] = byte('\n')
			resp.Body.Close()
		}

		session.Send(packet.NewBufferByBytes(buff,uint32(len(buff))))
		session.Close(1) 

	})

	go func(){
		p.Start()
	}()

	fmt.Printf("server start on: %s\n",service)

	
	tcpServer.Start( func (session *kendynet.StreamConn,msg *packet.ByteBuffer,errno error){	
		if msg == nil{
			fmt.Printf("error:%s\n",errno)
			session.Close()
			return
		}else{
			p.PostEvent(session)
		}	
	})
}

