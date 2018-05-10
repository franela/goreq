package goreq

import (
	"bytes"
	"io"
)

type compressReader struct {
	source     io.Reader
	compressor io.WriteCloser
	dest       io.Reader
	eof        bool
}

func newCompressReader(source io.Reader, writerGen func(buffer io.Writer) (io.WriteCloser, error)) (*compressReader, error) {
	buffer := bytes.NewBuffer([]byte{})
	writer, err := writerGen(buffer)
	if err != nil {
		return nil, err
	}
	return &compressReader{source: source, compressor: writer, dest: buffer}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	n, err = cr.source.Read(p)
	if n > 0 {
		i, e := cr.compressor.Write(p[:n])
		if e != nil {
			return i, e
		}
	}
	if err == io.EOF && !cr.eof {
		cr.eof = true
		cr.compressor.Close()
	}
	return cr.dest.Read(p)
}

func (cr *compressReader) Close() error {
	if v, ok := cr.source.(io.ReadCloser); !ok {
		if v != nil {
			v.Close()
		}
	}
	cr.compressor.Close()
	return nil
}
