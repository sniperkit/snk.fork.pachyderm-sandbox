package example

import(
	"io"
)

type BufferWriter struct {
	Content []byte
}

func NewBufferWriter() *BufferWriter {
	return &BufferWriter{}
}
func (bw *BufferWriter) Write(p []byte) (n int, err error) {
	bw.Content = append(bw.Content, p...)
	return len(p), nil
}

type CacheReader struct {
	content []byte
	index int
}

func NewCacheReader(content []byte) *CacheReader {
	return &CacheReader{
		content: content,
		index: 0,
	}
}

func (cr *CacheReader) Read(p []byte) (n int, err error) {
	if len(p) < ( len(cr.content) - cr.index ) {
		p = cr.content[cr.index:len(p)-1]
		cr.index = len(p)

		return len(p), nil
	}

	bufferSize := len(cr.content) - cr.index
	p = append(p, cr.content[cr.index:]...)

	return bufferSize, io.EOF
}
