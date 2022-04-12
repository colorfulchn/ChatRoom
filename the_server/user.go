package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct{
	UserConn net.Conn
	server *Server
	C chan string
	name string
}

func NewUser(conn net.Conn,ser *Server) *User{
	user := &User{
		UserConn: conn,
		server: ser,
		name:conn.RemoteAddr().String(),
		C: make(chan string),
	}
	go user.MessageListener();
	return user
}




//发送消息,C一有消息就写到对端
func (this *User) MessageListener(){
	for{
		msg := <-this.C
		fmt.Println("there is new news in Userlistener...")
		this.UserConn.Write([]byte(msg+"\n"));
	}
}
func (this *User) Online(){
	fmt.Println("Online...")
	name := this.name;
	this.server.mapLock.Lock()
	this.server.OnlineUserMap[name] = this;
	this.server.mapLock.Unlock();
	//NewToKnowMsg := "\n欢迎来到balabala-House\n"+
	//	"(1)输入#who，可以知道目前谁在线\n" +
	//	"(2)输入#rename|张三，可以把名字改成 张三\n" +
	//	"(3)您现在的名字为" + this.name +"\n"+
	//	"互联网不是法外之地，请和谐聊天\n\n";
	//this.SendMessageToSelf(NewToKnowMsg);
	this.server.Broadcast("上线了",name)
}

//offline
func (this *User) Offline(){
	fmt.Println("Offline...")
	name := this.name;
	this.server.mapLock.Lock()
	delete(this.server.OnlineUserMap,name);
	this.server.mapLock.Unlock();
	this.server.Broadcast("下线了",name)
}

func (this *User) SendMessageToSelf(msg string){
	this.UserConn.Write([]byte(msg));
}
func (this *User) SendMessageToSomeBody(){
	fmt.Println("SendMessageToSomeBody")
}
func (this *User) DoMessage(msg string){
	if len(msg) <=0{
		this.SendMessageToSelf("你输入了空字符串")
		return
	}
	if msg[0]=='#'{
		if msg=="#who"{ //#who
			this.server.mapLock.Lock();
			for key,_ := range this.server.OnlineUserMap{

				tmp := "[ " + key+ " ]" + "在线！"+" \n"
				this.SendMessageToSelf(tmp)
			}
			this.server.mapLock.Unlock();
		}
		if len(msg)>=8 && msg[1:8]=="rename|"{
			NewName := msg[8:]
			_,ok := this.server.OnlineUserMap[NewName];
			if(len(NewName) >= 30){
				this.SendMessageToSelf("你输入的名字太长了\n")
			}else{
				if ok {
					this.SendMessageToSelf("这个用户名已经被占用了\n")
				}else{
					this.server.mapLock.Lock();
					delete(this.server.OnlineUserMap,this.name)
					this.server.OnlineUserMap[NewName] = this
					this.name=NewName;
					this.server.mapLock.Unlock();
				}
			}
		}
		if len(msg) > 5 && msg[:4] == "#to|"{
			RemoteName := strings.Split(msg, "|")[1]
			if RemoteName == ""{
				this.SendMessageToSelf("您输入的姓名为空\n")
				return
			}
			RemoteUser,ok := this.server.OnlineUserMap[RemoteName]
			if !ok{
				this.SendMessageToSelf("您输入的姓名不存在\n")
				return
			}
			content := strings.Split(msg, "|")[2]
			if content==""{
				this.SendMessageToSelf("您输入的聊天内容为空\n")
				return
			}
			RemoteUser.SendMessageToSelf("---------【密】：[" + this.name + "] 对您说:" +content +"\n")
		}
	}else{
		this.server.Broadcast(msg,this.name)
	}
}