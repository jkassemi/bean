package bean

import (
  "fmt"
  "testing"
  "net/http"
  "net/http/httptest"
)

func TestTestGetRequest(t *testing.T){
  dummy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    if r.FormValue("a") != "Hello" {
      t.Error("Form value issue: a != 'Hello'")
    }

    if r.Header.Get("b") != "World" {
      t.Error("Header issue: b != 'World'")
    }

    fmt.Fprintf(w, "Hello World")
  }))

  defer dummy.Close()

  data := map[string] string {
    "a": "Hello",
  }

  headers := map[string] string {
    "b": "World",
  }

  TestGetRequest(dummy.URL, data, headers, t)
}

func TestTestPostRequest(t *testing.T){
  dummy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    if r.FormValue("a") != "Hello" {
      t.Error("Form value issue: a != 'Hello'")
    }

    if r.Header.Get("b") != "World" {
      t.Error("Header issue: b != 'World'")
    }
  }))

  defer dummy.Close()

  data := map[string] string {
    "a": "Hello",
  }

  headers := map[string] string {
    "b": "World",
  }

  TestPostRequest(dummy.URL, data, headers, t)
}

func TestAssertRedirectedTo(t *testing.T){
  dummy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Location", "/hello_world")
  }))

  defer dummy.Close()

  tr := TestGetRequest(dummy.URL, nil, nil, t)
  tr.AssertRedirectedTo(fmt.Sprintf("%s/hello_world", dummy.URL), t)
}

func TestAssertContains(t *testing.T){
  dummy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello World")
  }))

  defer dummy.Close()

  tr := TestGetRequest(dummy.URL, nil, nil, t)
  tr.AssertContains("Hello World", t)
}

func TestAssertSelector(t *testing.T){
  dummy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, `<div class="content"><span id="test">Hello World</span></div>`)
  }))

  defer dummy.Close()

  tr := TestGetRequest(dummy.URL, nil, nil, t)
  tr.AssertSelector("div.content span#test", t)
}

func TestAssertCode(t *testing.T){
  dummy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusNotFound)
  }))

  defer dummy.Close()

  tr := TestGetRequest(dummy.URL, nil, nil, t)
  tr.AssertCode(http.StatusNotFound, t)
}
