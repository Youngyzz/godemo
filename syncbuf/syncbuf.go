//youngy 2018/5/7 9:57
package syncbuf

import (
	"bufio"
	"bytes"
	"sync"
)

const (
	NonBlocking = iota
	Blocking
)

type SyncBuf struct {
	bytes.Buffer
	block *int
	lock  sync.Locker
	cond  *sync.Cond
}

func NewSyncBuf(b []byte, blocking *int, l sync.Locker, c *sync.Cond) *SyncBuf {
	buf := &SyncBuf{
		Buffer: *bytes.NewBuffer(b),
		block:  blocking,
		lock:   l,
		cond:   c,
	}

	return buf
}

func (self *SyncBuf) Write(b []byte) (int, error) {
	self.cond.L.Lock()
	self.cond.Signal()
	defer self.cond.L.Unlock()
	return self.Buffer.Write(b)
}

func (self *SyncBuf) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	self.cond.L.Lock()
	for self.Buffer.Len() <= 0 && *self.block == Blocking {
		self.cond.Wait()
	}

	defer self.cond.L.Unlock()
	return self.Buffer.Read(b)
}

type SyncReadWriter struct {
	/*
		设置是否阻塞，非阻塞，立即返回，如果没有消息可读时，会返回io.EOF错误，阻塞，阻塞直到有消息可读。
		这个值用来设置RW阻塞或非阻塞。
	*/
	blocking int
	lock     sync.Locker
	cond     *sync.Cond
	RW       *bufio.ReadWriter
}

func (self *SyncReadWriter) SetBlocking(b int) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.cond.Signal()
	self.blocking = b
}

func NewSyncReadWriter() *SyncReadWriter {
	rw := &SyncReadWriter{
		blocking: Blocking,
		lock:     &sync.Mutex{},
	}
	rw.cond = sync.NewCond(rw.lock)
	b := NewSyncBuf(make([]byte, 0), &rw.blocking, rw.lock, rw.cond)
	br := bufio.NewReader(b)
	bw := bufio.NewWriter(b)
	rw.RW = bufio.NewReadWriter(br, bw)
	return rw
}
