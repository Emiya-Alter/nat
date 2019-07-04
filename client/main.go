package main

import (
	"flag"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	local := flag.String("local", "127.0.0.1:80", "本地客户端访问ip端口")
	server := flag.String("server", "88.88.88.88:50050", "服务端数据转发ip端口")
	flag.Parse()
	log.Printf( "start to connect address: %s and address: %s", *local, *server)
	host2host(*local, *server)
}


func host2host(address1, address2 string) {
	for {
		log.Printf("try to connect host: %s and %s", address1, address2)
		var host1, host2 net.Conn
		var err error
		for {
			host1, err = net.Dial("tcp", address1)
			if err == nil {
				log.Printf("connect %s success.", address1)
				break
			} else {
				log.Printf("connect target address %s faild. retry in 3 seconds.", address1)
				time.Sleep(3 * time.Second)
			}
		}
		for {
			host2, err = net.Dial("tcp", address2)
			if err == nil {
				log.Printf( "connect %s success.", address2)
				break
			} else {
				log.Printf("connect target address %s faild. retry in 3 seconds.", address2)
				time.Sleep(3 * time.Second)
			}
		}
		forward(host1, host2)
	}
}


func forward(conn1 net.Conn, conn2 net.Conn) {
	log.Printf("starting transmit. [%s],[%s] <-> [%s],[%s] \n", conn1.LocalAddr().String(), conn1.RemoteAddr().String(), conn2.LocalAddr().String(), conn2.RemoteAddr().String())
	var wg sync.WaitGroup
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	wg.Wait()
}


func connCopy(conn1 net.Conn, conn2 net.Conn, wg *sync.WaitGroup) {
	_, _ = io.Copy(conn1, conn2)
	_ = conn1.Close()
	log.Printf("close the connect at local: %s and remote: %s", conn1.LocalAddr().String(), conn1.RemoteAddr().String())
	wg.Done()
}
