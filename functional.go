// Functional testing helpers that allow simple requests
// and assertions based on those requests.

package bean

import (
  "fmt"
  "errors"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "testing"
  "strings"
  "code.google.com/p/go-html-transform/html/transform"
)

type TestResponse struct {
  Url         string
  Request     *http.Request
  Response    *http.Response
}

func encodedParams(data map[string] string) string {
  params := &url.Values{}

  for k, v := range data {
    params.Set(k, v)
  }

  return params.Encode()
}

func doRequest(method string, request_url string, request_body io.Reader, headers map[string] string, t *testing.T) (*TestResponse){
  req, _ := http.NewRequest(method, request_url, request_body)

  if headers != nil {
    for k, v := range headers {
      req.Header.Set(k, v)
    }
  }

  if method == "POST" {
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  }

  client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
    return errors.New("stop") 
  }}

  var response *http.Response

  response, e := client.Do(req)

  if e != nil {
    t.Errorf("Bad response: %s", e.Error())
  }

  return &TestResponse{Url: request_url, Request: req, Response: response}
}

// Issue a GET request to the server and retrieve a TestResponse for
// later assertions.
func TestGetRequest(test_url string, data map[string] string, headers map[string] string, t *testing.T) (*TestResponse){
  var request_url string

  if data != nil {
    params := encodedParams(data)
    request_url = fmt.Sprintf("%s?%s", test_url, params)
  } else {
    request_url = test_url
  }

  return doRequest("GET", request_url, nil, headers, t)
}

// Issue a POST request to the server and retrieve a TestResponse for
// later assertions.
func TestPostRequest(request_url string, data map[string] string, headers map[string] string, t *testing.T) (*TestResponse){
  var params *strings.Reader

  if data != nil {
    params = strings.NewReader(encodedParams(data))
  }else{
    params = nil
  }

  return doRequest("POST", request_url, params, headers, t)
}

// Asserts that the TestResponse requests a client redirection 
// to the specified resource.
func (r *TestResponse) AssertRedirectedTo(url string, t *testing.T) {
  if l, e := r.Response.Location(); e != nil {
    t.Errorf("Response not a redirect: %s", r.Response.Status)
  } else if l.String() != url {
    t.Errorf("Redirect not matched: %s does not equal %s", l.String(), url)
  }
}

// Asserts that the TestResponse contains the specified string.
func (r *TestResponse) AssertContains(text string, t *testing.T){
  body, e := ioutil.ReadAll(r.Response.Body)

  if e != nil {
    t.Errorf("Response not readable: %s", r.Url)
  }

  if !strings.Contains(text, string(body)) {
    t.Errorf("Response body does not contain %s", text)
  }
}

// Asserts that the TestResponse contains a CSS-style selector. 
// This method supports the SelectorQuery style selectors from 
// go-html-transform, but accepts them as a single string with
// space delimiters
func (r *TestResponse) AssertSelector(selector string, t *testing.T){
  body, e := ioutil.ReadAll(r.Response.Body)

  if e != nil {
    t.Errorf("Response not readable: %s", r.Url)
  }

  doc, e := transform.NewDoc(string(body))

  if e != nil {
    t.Errorf("Could not parse response: %s", r.Url)
  }

  selector_parts := strings.Split(selector, " ")

  q := make(transform.SelectorQuery, len(selector_parts))

  for i, str := range selector_parts {
    selector := transform.NewSelector(str)

    if selector == nil {
      t.Errorf("Problem with selector: %s", str)
    }

    q[i] = selector
  }

  if len(q.Apply(doc)) == 0 {
    t.Errorf("Selector not found: %s", selector)
  }
}

// Assert the response is of the specified status code
func (r *TestResponse) AssertCode(code int, t *testing.T){
  if r.Response.StatusCode != code {
    t.Errorf("Invalid response code. Expected %d, Received %d", code, r.Response.StatusCode)
  }
}
