package io

import "bytes"

// BufferCloser is an bytes.buffer implementation of io.WriterCloser
type BufferCloser struct {
	bytes.Buffer
}

// NewBufferCloser returns a new BufferCloser
func NewBufferCloser() *BufferCloser {
	return &BufferCloser{
		Buffer: bytes.Buffer{},
	}
}

// Close closes the BufferCloser
func (bc *BufferCloser) Close() error {
	// No-op close
	return nil
}

// Write writes to the BufferCloser
func (bc *BufferCloser) Write(p []byte) (n int, err error) {
	return bc.Buffer.Write(p)
}
