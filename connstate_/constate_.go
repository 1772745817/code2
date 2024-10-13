package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	servel := &http.Server{
		Addr: ":8080",
		ConnState: func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				fmt.Println("StateNew : 新的连接已建立")
			case http.StateActive:
				fmt.Println("StateActive : 连接已激活,正在处理请求")
			case http.StateIdle:
				fmt.Println("StateIdle : 连接已空闲,等待新请求")
			case http.StateHijacked:
				fmt.Println("StateHijacked : 连接已劫持")
			case http.StateClosed:
				fmt.Println("StateClose : 连接已关闭")
			}
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello!")
	})

	fmt.Println("服务器正在监听端口 : 8080")
	if err := servel.ListenAndServe(); err != nil {
		fmt.Println("服务器发生错误 ", err)
	}

}
