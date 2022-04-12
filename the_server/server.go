package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int
	MessagerChan    chan string
	OnlineUserMap map[string]*User
	mapLock sync.RWMutex
}

// get a server struct
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
		MessagerChan : make(chan string),
		OnlineUserMap: make(map[string]*User),
	}
	fmt.Println("NerServer...")
	return server
}

// handler
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Handler...")
	user := NewUser(conn,this)
	user.Online()
	isLive := make(chan bool)
	go func(){

		for{
			select{
			case <- isLive:
			case <- time.After(time.Second * 600):
				user.SendMessageToSelf("你过长时间没有说话，被提出了聊天室\n")
				user.Offline()
				// close(user.C)
				conn.Close()
				return
			}
		}
	}()

	for{
		buffer := make([]byte,4096);
		n,err := conn.Read(buffer)
		if n==0{
			user.Offline()
			return
		}
		if err!=nil && n!=0{
			fmt.Println("Handler:Read error",err);
		}
		msg := string(buffer[:n-1]);
		user.DoMessage(msg);
		isLive <- true
		//有消息就往isLive发一条消息
	}


}

func (this *Server)Broadcast(msg string,name string){
	fmt.Println("Broadcast...")
	msgSend := "[ "+ name +" ]"+":"+ msg ;
	this.MessagerChan <- msgSend;
}
func LockAndDo(lock sync.RWMutex,do func()){
	lock.Lock();
	do();
	lock.Unlock();
}
func (this *Server) MessageListener(){
	fmt.Println("MessageListener...")

	for{
		msg := <-this.MessagerChan
		fmt.Println("there is new news...")
		//this.mapLock.Lock();
		//for _,value := range(this.OnlineUserMap){
		//	value.C <- msg
		//}
		//this.mapLock.Unlock();
		LockAndDo(this.mapLock,func(){
			for _,value := range(this.OnlineUserMap){
				value.C <- msg
		}})
	}
}



func (this *Server) Start() {
	fmt.Println("Start...")
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net.Listen error", err)
	}
	defer Listener.Close()
	go this.MessageListener(); // 有消息就发送给chan
	for {
		//accet
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Accept error", err)
			continue
		}

		//处理请求
		go this.Handler(conn)
	}

}
