package loggermw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHeader(t *testing.T) {

	rec := httptest.NewRecorder()
	crw := customResponseWriter{
		ResponseWriter: rec,
		status:         0,
		size:           0,
	}

	crw.WriteHeader(http.StatusOK)
	if crw.status != http.StatusOK {
		t.Errorf("Expected: %d, Got: %d", http.StatusOK, crw.status)
	}
}

func TestWrite(t *testing.T) {

	rec := httptest.NewRecorder()
	crw := customResponseWriter{
		ResponseWriter: rec,
		status:         0,
		size:           0,
	}

	buffer := []byte{'a', 'b', 'c'}

	crw.Write(buffer)
	if crw.size != len(buffer) {
		t.Errorf("Expected: %d, Got: %d", len(buffer), crw.size)
	}
}
