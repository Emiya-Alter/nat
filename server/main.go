package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	cp := flag.Int("cp", 50050, "内网客户端连接端口")
	sp := flag.Int("sp", 50060, "其他客户端访问端口")
	flag.Parse()
	if *cp < 0 || *cp > 65535 || *sp < 0 || *sp > 65535 {
		log.Fatalln("please input port in range [0, 65535]")
	}
	log.Printf("listening port %d , %d", *cp, *sp)
	listenPort(*cp, *sp)

}


func listenServer(address string)  (server net.Listener) {
	log.Printf("try to start server on: %s", address)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("listen %s faild.", address)
	}
	log.Printf("start listen at address: %s", address)
	return
}


func listenPort(cp , sp int){
	listen1 := listenServer("0.0.0.0:" + strconv.Itoa(cp))
	listen2 := listenServer("0.0.0.0:" + strconv.Itoa(sp))
	log.Printf("listening port: %d and port: %d , waiting for client", cp, sp)
	for {
		conn1 := acceptConn(listen1)
		conn2 := acceptConn(listen2)
		if conn1 == nil || conn2 == nil {
			log.Println("accept client faild. retry in 3 seconds.")
			time.Sleep(3 * time.Second)
			continue
		}
		forward(conn1, conn2)
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


func acceptConn(listener net.Listener) (conn net.Conn) {
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("accept connect %s faild. %s", conn.RemoteAddr().String(), err.Error())
		return nil
	}
	log.Printf("accept a new client. remote address: %s local address: %s ", conn.RemoteAddr().String(), conn.LocalAddr().String())
	return
}
