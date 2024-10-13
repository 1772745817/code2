package main

import (
	"fmt"
	"net"
	"net/http"
)

func handel1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from port 8080!")
}

func handel2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from port 8081")
}
func main() {
	listen1, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	listen2, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("start servel on port 8080")
		//http.HandlerFunc()可以将普通函数转化成实现了Handler接口的类型,它实现了http.Handler接口
		//http.Handle是用于注册http路由的函数
		//目的是将特定路径（或路由）映射到相应的处理逻辑(函数)上。
		//func Handle(pattern string, handler Handler)
		err := http.Serve(listen1, http.HandlerFunc(handel1))
		if err != nil {
			fmt.Println("Erro serving on port 8080")
		}
	}()

	go func() {
		fmt.Println("start servel on port 8081")
		err := http.Serve(listen2, http.HandlerFunc(handel2))
		if err != nil {
			fmt.Println("Erro serving on port 8081")
		}
	}()

	select {}
}
