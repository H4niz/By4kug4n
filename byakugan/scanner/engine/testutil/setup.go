package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "github.com/haniz/byakugan/scanner/proto"
)

type TestServer struct {
	Server *httptest.Server
	Tasks  chan *pb.ScanTask
}

func SetupTestServer(t *testing.T) *TestServer {
	tasks := make(chan *pb.ScanTask, 100)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users/login":
			handleLoginEndpoint(w, r)
		}
	}))

	return &TestServer{
		Server: server,
		Tasks:  tasks,
	}
}

func handleLoginEndpoint(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	// Simulate SQL injection vulnerability
	if strings.Contains(req["username"], "'") {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"SQL syntax error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
