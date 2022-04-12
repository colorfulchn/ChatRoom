package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)
type Client struct{
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int // model of chat
}
var serverIp string
var serverPort int
func NewClient(serverIp string,serverPort int) *Client{
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag : 9,
	}
	conn , err := net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err!=nil{
		fmt.Println("dial error",err)
		return nil
	}
	client.conn = conn
	return client
}
func (this *Client) DealListener(){
	io.Copy(os.Stdout,this.conn)
}
func (this *Client) meau() bool{
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查询当前在线用户")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag<=4{
		this.flag =flag
		return true
	}else{
		fmt.Println(">>>>>>请输入合法的模式数字")
		return false
	}
}
func (this *Client) PublicChat(){
	var msg string
	fmt.Println("PublicChat model")
	fmt.Println(">>>>>>请输入你要输入的聊天内容,quit退出")
	fmt.Scanln(&msg)
	for msg != "quit"{
		if len(msg) != 0{
			sendMsg := msg + "\n"
			_,err := this.conn.Write([]byte(sendMsg))
			if err !=nil{
				fmt.Println("conn Write err:",err)
				break
			}
		}
		msg = ""
		fmt.Scanln(&msg)
	}
}
func (this *Client) PrivateChat(){
	var privateChatUser string
	var msg string
	fmt.Println("PrivateChat model")
	//beginning--2
	this.QueryUsersHelper()
	fmt.Println(">>>>>>请输入私聊对象的用户名，quit退出:")
	fmt.Scanln(&privateChatUser)
	for privateChatUser != "quit"{
		//beginning--1
		fmt.Println(">>>>>>请输入私聊的内容，quit退出:")
		fmt.Scanln(&msg)
		for msg!="quit"{
			if len(msg)!=0{
				sendmsg := "#to|"+ privateChatUser + "|" +msg + "\n"
				_,err := this.conn.Write([]byte(sendmsg))
				if err !=nil{
					fmt.Println("PrivateChat Write error",err)
					break
				}
			}
			msg = ""
			fmt.Println(">>>>>>请输入私聊的内容，quit退出:")
			fmt.Scanln(&msg)
		}
		//ending--1
		this.QueryUsersHelper()
		fmt.Println(">>>>>>请输入私聊对象的用户名，quit退出:")
		fmt.Scanln(&privateChatUser)
	}
	//ending--2
}
func (this *Client) UpdateName(){
	fmt.Println("PrivateChat model")
	fmt.Println(">>>>>>请输入你要更改的账户名")
	var tmpName string
	fmt.Scanln(&tmpName)
	if tmpName =="quit"{
		fmt.Println("这个用户名不合法")
		return
	}
	this.Name = tmpName
	msg := "#rename|" + this.Name + "\n";
	_,err := this.conn.Write([]byte(msg));
	if err !=nil {
		fmt.Println("UpdateName error,the username has been used,",err)
	}
}
func (this *Client) QueryUsersHelper(){
	msg := "#who" + "\n";
	_,err := this.conn.Write([]byte(msg))
	if err !=nil{
		fmt.Println("QueryUsersHelper error",err)
		return
	}
}
func (this *Client) QueryUsers(){
	fmt.Println("QueryUsers model")
	this.QueryUsersHelper()
	return
}
func (this *Client) Run(){
	for this.flag!=0{
		for this.meau()!=true{ }// 输入的不合法 就一直回去输入
		switch this.flag {
		case 1:
			this.PublicChat()
			break;
		case 2:
			this.PrivateChat()
			break;
		case 3:
			this.UpdateName()
			break;
		case 4:
			this.QueryUsers()
			break;
		}
	}
}
func init(){
	flag.StringVar(&serverIp,"ip","124.223.0.129","设置服务器IP地址(默认是124.223.0.129)")
	flag.IntVar(&serverPort,"port",7171,"设置服务器端口（默认是7171）")
}
func main(){
	flag.Parse();
	client := NewClient(serverIp,serverPort)
	if client== nil{
		fmt.Println("......服务器链接失败......")
		return
	}

	go client.DealListener();
	fmt.Println("......服务器链接成功......")
	client.Run()

}