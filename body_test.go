package luddite

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

type sample struct {
	XMLName   xml.Name  `json:"-" xml:"sample"`
	Id        int       `json:"id" xml:"id"`
	Name      string    `json:"name" xml:"name"`
	Flag      bool      `json:"flag" xml:"flag"`
	Data      []byte    `json:"data" xml:"data"`
	Timestamp time.Time `json:"timestamp" xml:"timestamp"`
}

const (
	sampleId       = 1234
	sampleName     = "dave"
	sampleData     = "Hello world"
	sampleJsonBody = "{\"id\":1234,\"name\":\"dave\",\"flag\":true,\"data\":\"SGVsbG8gd29ybGQ=\",\"timestamp\":\"2015-03-18T14:30:00Z\"}"
	sampleXmlBody  = "<sample><id>1234</id><name>dave</name><flag>true</flag><data>Hello world</data><timestamp>2015-03-18T14:30:00Z</timestamp></sample>"
)

var (
	sampleTimestamp = time.Date(2015, 3, 18, 14, 30, 0, 0, time.UTC)
)

type sampleResource struct {
	NotImplementedResource
}

func (r *sampleResource) New() interface{} {
	return &sample{}
}

func (r *sampleResource) Id(value interface{}) string {
	return strconv.Itoa(value.(*sample).Id)
}

func TestReadJson(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", strings.NewReader(sampleJsonBody))
	req.Header[HeaderContentType] = []string{ContentTypeJson}

	v, err := readRequest(req, &sampleResource{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	s := v.(*sample)
	if s.Id != sampleId {
		t.Error("JSON int deserialization failed")
	}
	if s.Name != sampleName {
		t.Error("JSON string deserialization failed")
	}
	if !s.Flag {
		t.Error("JSON bool deserialization failed")
	}
	if !bytes.Equal(s.Data, []byte(sampleData)) {
		t.Error("JSON binary deserialization failed")
	}
	if s.Timestamp != sampleTimestamp {
		t.Error("JSON date deserialization failed")
	}
}

func TestWriteJson(t *testing.T) {
	s := &sample{
		Id:        sampleId,
		Name:      sampleName,
		Flag:      true,
		Data:      []byte(sampleData),
		Timestamp: sampleTimestamp,
	}

	rw := httptest.NewRecorder()
	rw.Header().Add(HeaderContentType, ContentTypeJson)

	if err := writeResponse(rw, http.StatusOK, s); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if rw.Code != http.StatusOK {
		t.Error("status code never written")
	}

	if rw.Body != nil {
		if body := string(rw.Body.String()); body != sampleJsonBody {
			t.Errorf("JSON serialization failed, got: %s, expected: %s\n", body, sampleJsonBody)
		}
	} else {
		t.Error("body never written")
	}
}

func TestReadXml(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", strings.NewReader(sampleXmlBody))
	req.Header[HeaderContentType] = []string{ContentTypeXml}

	v, err := readRequest(req, &sampleResource{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	s := v.(*sample)
	if s.Id != sampleId {
		t.Error("XML int deserialization failed")
	}
	if s.Name != sampleName {
		t.Error("XML string deserialization failed")
	}
	if !s.Flag {
		t.Error("XML bool deserialization failed")
	}
	if !bytes.Equal(s.Data, []byte(sampleData)) {
		t.Error("XML binary deserialization failed")
	}
	if s.Timestamp != sampleTimestamp {
		t.Error("XML date deserialization failed")
	}
}

func TestWriteXml(t *testing.T) {
	s := &sample{
		Id:        sampleId,
		Name:      sampleName,
		Flag:      true,
		Data:      []byte(sampleData),
		Timestamp: sampleTimestamp,
	}

	rw := httptest.NewRecorder()
	rw.Header().Add(HeaderContentType, ContentTypeXml)

	if err := writeResponse(rw, http.StatusOK, s); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if rw.Code != http.StatusOK {
		t.Error("status code never written")
	}

	if rw.Body != nil {
		if body := string(rw.Body.String()); body != sampleXmlBody {
			t.Errorf("XML serialization failed, got: %s, expected: %s\n", body, sampleXmlBody)
		}
	} else {
		t.Error("body never written")
	}
}

func TestWriteHtml(t *testing.T) {
	// Write []byte
	rw := httptest.NewRecorder()
	rw.Header().Add(HeaderContentType, ContentTypeHtml)

	if err := writeResponse(rw, http.StatusOK, []byte(sampleData)); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if rw.Code != http.StatusOK {
		t.Error("status code never written")
	}

	if rw.Body != nil {
		if body := string(rw.Body.String()); body != sampleData {
			t.Errorf("HTML body write failed, got: %s, expected: %s\n", body, sampleData)
		}
	} else {
		t.Error("body never written")
	}

	// Write string
	rw = httptest.NewRecorder()
	rw.Header().Add(HeaderContentType, ContentTypeHtml)

	if err := writeResponse(rw, http.StatusOK, sampleData); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if rw.Code != http.StatusOK {
		t.Error("status code never written")
	}

	if rw.Body != nil {
		if body := string(rw.Body.String()); body != sampleData {
			t.Errorf("HTML body write failed, got: %s, expected: %s\n", body, sampleData)
		}
	} else {
		t.Error("body never written")
	}

	// Write other type
	s := &sample{
		Id:        sampleId,
		Name:      sampleName,
		Flag:      true,
		Data:      []byte(sampleData),
		Timestamp: sampleTimestamp,
	}

	rw = httptest.NewRecorder()
	rw.Header().Add(HeaderContentType, ContentTypeHtml)

	if err := writeResponse(rw, http.StatusOK, s); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if rw.Code != http.StatusOK {
		t.Error("status code never written")
	}
}
