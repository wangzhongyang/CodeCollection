package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	r := gin.Default()
	r.POST("/test", MiddleFunc, DoPost)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func MiddleFunc(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalf("read body failed at Before,err:%s", err.Error())
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	log.Println("read body: ", string(body))

	r := NewLogResponseWriter(c.Writer)
	c.Writer = r
	
	c.Next()

	log.Println("response body:", r.buf.String())
}

func DoPost(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalf("read body failed at Before,err:%s", err.Error())
	}
	var req Body
	if err = json.Unmarshal(body, &req); err != nil {
		log.Fatalf("json unmarshal failed, body:%s, err:%s", string(body), err.Error())
	}
	c.JSON(http.StatusOK, fmt.Sprintf("this is response, get body id:%s", req.Id))
}

type LogResponseWriter struct {
	buf    *bytes.Buffer
	writer gin.ResponseWriter
}

func (r *LogResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.writer.Hijack()
}

func (r *LogResponseWriter) Flush() {
	r.writer.Flush()
}

func (r *LogResponseWriter) CloseNotify() <-chan bool {
	return r.writer.CloseNotify()
}

func (r *LogResponseWriter) Status() int {
	return r.writer.Status()
}

func (r *LogResponseWriter) Size() int {
	return r.writer.Size()
}

func (r *LogResponseWriter) WriteString(s string) (int, error) {
	return r.writer.WriteString(s)
}

func (r *LogResponseWriter) Written() bool {
	return r.writer.Written()
}

func (r *LogResponseWriter) WriteHeaderNow() {
	r.writer.WriteHeaderNow()
}

func (r *LogResponseWriter) Pusher() http.Pusher {
	return r.writer.Pusher()
}

func NewLogResponseWriter(r gin.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{
		buf:    new(bytes.Buffer),
		writer: r,
	}
}

func (r *LogResponseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *LogResponseWriter) Write(b []byte) (int, error) {
	r.buf = bytes.NewBuffer(b)
	return r.writer.Write(b)
}

func (r *LogResponseWriter) WriteHeader(statusCode int) {
	r.writer.WriteHeader(statusCode)
}

type Body struct {
	Id string `json:"id"`
}
