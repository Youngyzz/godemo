// youngy 2018/5/7 9::57
package syncbuf

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSyncReadHttpRequest(t *testing.T) {
	rw := NewSyncReadWriter()

	//write http request one by one for line.
	go func(*SyncReadWriter) {
		req := []string{"GET /sample.jsp HTTP/1.1\r\n",
			"Accept:image/gif,image/jpeg\r\n",
			"Accept-Language:zh-cn\r\n",
			"Content-Length:30\r\n",
			"Connection:Keep-Alive\r\n",
			"Host:localhost\r\n",
			"User-Agent:Mozila/4.0(compatible;MSIE5.01;Window NT5.0)\r\n",
			"Accept-Encoding:gzip,deflate\r\n\r\n",
			"username=jinqiao&password=1234\r\n",
		}
		//write one by one for line.
		for _, r := range req {
			rw.RW.Write([]byte(r))
			time.Sleep(time.Millisecond * 10)
			rw.RW.Flush()
		}
	}(rw)

	//read http request for http
	go func(*SyncReadWriter) {
		//blocking read for whole http requst.
		req, err := http.ReadRequest(rw.RW.Reader)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(req)
			b := make([]byte, 100)
			if req.Body != nil {
				req.Body.Read(b)
			}
			fmt.Println(string(b))
		}
	}(rw)

	//wait for 3 seconds and then set rw nonblocking,
	//this may be caused io.EOF error when there was no data to read.
	func(*SyncReadWriter) {
		time.Sleep(time.Second * 3)
		rw.SetBlocking(NonBlocking)
	}(rw)

}
