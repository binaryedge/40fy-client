package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var (
	token          = "token"
	jobID          = "1234"
	cmdWithJobID   = []string{"-token=" + token, "-job-id=" + jobID}
	cmd            = []string{"-token=" + token}
	serverResponse = []byte("test")
)

func TestCmdWithJobID(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(serverResponse)
		if h := r.Header.Get("X-Token"); h != token {
			t.Fatal("Token is different", h, " != ", token)
		}
		if jobid := r.URL.Query().Get("job_id"); jobID != jobid {
			t.Fatal("URL doesnt contain jobID ", r.URL)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	buffer := bytes.NewBuffer([]byte{})
	c := StreamCommand{http.Client{}, server.URL, buffer, ""}

	if status := c.Run(cmdWithJobID); status != 0 {
		t.Fatal("Status not 0", status, " != ", 0)
	}
	if !reflect.DeepEqual(buffer.Bytes(), serverResponse) {
		t.Fatal("Server Response is different ", buffer.Bytes(), " != ", serverResponse)
	}

}

func TestCmdWithoutJobID(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(serverResponse)
		if h := r.Header.Get("X-Token"); len(h) > 0 && h != token {
			t.Fatal("Token is different", h, " != ", token)
		}
		t.Log("url ", r.URL)
		if jobid := r.URL.Query().Get("job_id"); jobid != "" {
			t.Fatal("URL doesnt contain jobID ", jobid)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	buffer := bytes.NewBuffer([]byte{})
	c := StreamCommand{http.Client{}, server.URL, buffer, ""}

	if status := c.Run(cmd); status != 0 {
		t.Fatal("Status not 0 ", status, " != ", 0)
	}
	if !reflect.DeepEqual(buffer.Bytes(), serverResponse) {
		t.Fatal("Server Response is different ", buffer.Bytes(), " != ", serverResponse)
	}
}

func TestCmdWithoutToken(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(serverResponse)
		if h := r.Header.Get("X-Token"); len(h) > 0 && h != token {
			t.Fatal("Token is different", h, " != ", token)
		}
		t.Log("url ", r.URL)
		if jobid := r.URL.Query().Get("job_id"); jobid != "" {
			t.Fatal("URL doesnt contain jobID ", jobid)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	buffer := bytes.NewBuffer([]byte{})
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(file.Name())

	file.Write([]byte(token))

	c := StreamCommand{http.Client{}, server.URL, buffer, file.Name()}
	if status := c.Run([]string{}); status != 0 {
		t.Fatal("Status should be 0 ", status, " != ", 0)
	}
}
