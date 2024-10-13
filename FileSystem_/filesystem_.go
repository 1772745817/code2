package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 自定义 MYSQL 系统文件
type dbFileSystem struct {
	db *sql.DB
}

// 自定义实现 http.File 接口
type dbFile struct {
	content []byte
	size    int64
	modtime time.Time
	name    string
}

// 实现 http.File的stat()方法
func (f *dbFile) Stat() (os.FileInfo, error) {
	return &dbFileInfo{
		name:    f.name,
		size:    f.size,
		modTime: f.modtime,
	}, nil
}

// 实现 http.File 接口的 Close 方法
func (f *dbFile) Close() error {
	// 没有真正的资源需要释放，所以可以留空
	return nil
}

// 实现 http.File 的 Read() 方法
func (f *dbFile) Read(p []byte) (n int, err error) {
	if len(f.content) == 0 {
		return 0, io.EOF // 如果文件内容为空，返回 EOF
	}

	// 读取文件内容到 p 中
	n = copy(p, f.content)
	f.content = f.content[n:] // 更新剩余内容
	return n, nil
}

type dbFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

// 实现fileInfo接口方法
func (f dbFileInfo) Name() string {
	return f.name
}

func (f dbFileInfo) Size() int64 {
	return f.size
}

func (f dbFileInfo) Mode() os.FileMode {
	return 0444
}
func (f dbFileInfo) ModTime() time.Time {
	return f.modTime
}
func (f dbFileInfo) IsDir() bool {
	return false
}
func (f dbFileInfo) Sys() interface{} {
	return nil
}

func (fs *dbFileSystem) Open(name string) (*dbFile, error) {
	var content []byte
	var size int64
	var modTimeBytes []byte

	query := "SELECT content ,size, modtime FROM file  WHERE name = ?"
	//查询数据
	row := fs.db.QueryRow(query, name)
	//将查询到的数据存储到定义的变量当中
	err := row.Scan(&content, &size, &modTimeBytes)
	if err == sql.ErrNoRows {
		return nil, os.ErrNotExist
	}
	if err != nil {
		return nil, err
	}

	// 将 []byte 转换为 time.Time
	modTime, err := time.Parse("2006-01-02 15:04:05", string(modTimeBytes))
	if err != nil {
		return nil, err
	}
	//创建dbFile实例
	return &dbFile{
		content: content,
		size:    size,
		modtime: modTime,
		name:    name,
	}, nil
}

// 提供文的http处理函数
func dbFileHandler(w http.ResponseWriter, r *http.Request) {
	//从URL中获取参数名
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "fileName is require", http.StatusBadRequest)
		return
	}
	//从dbFileSystem文件系统中打开文件
	file, err := fs.Open(fileName)
	if err != nil {
		//当函数或方法返回错误（即 error 类型）时，你可以通过调用 Error() 方法来获取该错误的详细描述。
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	//设置响应头
	// w.Header().Set("Content-Disposition", "attachment;filename="+fileName)

	//w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", file.size))

	//将文件写入响应
	//io.Copy(w, file)

	body := file.content

	w.Write(body)

}

var fs *dbFileSystem

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/filedb"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("successful connected to mysqldb")
	//创建dbFileSystem实例
	fs = &dbFileSystem{
		db: db,
	}

	http.HandleFunc("/file", dbFileHandler)

	fmt.Println("starting servel at port :8080")
	http.ListenAndServe(":8080", nil)
}
