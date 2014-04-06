package goreq

import (
	"bytes"
	"io"
	"mime/multipart"
)

// MultipartForm wraps functionality from the stdlib mime/multipart
// You can supply any io.Reader as the source, and goreq will send your
// request as a multipart form. For convenience, you can set the
// "FileName" attribute and that will be added to the field header

// Simple usage to upload a file looks like this
//    file, _ := os.Open("/tmp/myfile")
//    body := goreq.MultipartForm{
//        Field:    "file",
//        FileName: "myfile",
//        Source:   file,
//    }
//    resp, _ := goreq.Request{
//        Uri:          "http://some.place/uristuff",
//        Method:       "POST",
//        Body:         body,
//    }.Do()
type MultipartForm struct {
	Params   map[string]string
	Field    string
	Source   io.Reader
	buffer   bytes.Buffer
	FileName string
}

func (m *MultipartForm) ToReader() (_ io.Reader, err error) {
	writer := multipart.NewWriter(&m.buffer)

	for k, v := range m.Params { // fill in the extra form params
		writer.WriteField(k, v)
	}

	var multipartWriter io.Writer

	if m.FileName != "" {
		multipartWriter, err = writer.CreateFormFile(m.Field, m.FileName)
	} else {
		multipartWriter, err = writer.CreateFormField(m.Field)
	}

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(multipartWriter, m.Source)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return &m.buffer, nil
}
