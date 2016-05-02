package goreq

import "io"

type compressReader struct {
	source     io.Reader
	compressor io.WriteCloser
	dest       io.Reader
	eof        bool
}

func newCompressReader(source io.Reader, compressor io.WriteCloser, dest io.Reader) *compressReader {
	return &compressReader{source: source, compressor: compressor, dest: dest}
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	buf := make([]byte, len(p), cap(p))
	n, err = cr.source.Read(buf)
	if n > 0 {
		i, e := cr.compressor.Write(buf[:n])
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
