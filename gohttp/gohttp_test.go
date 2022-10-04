package main

import (
    "fmt"
    "testing"
    "net/http"
    "net/http/httptest"
    "encoding/json"
    "bytes"
)

func TestSortRelevanceScore(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "?sortKey=relevanceScore", nil)
    response := httptest.NewRecorder()

    s := Server { }
    s.gohttp(response, request)

    got := "unordered"

    buf := new(bytes.Buffer)
    buf.ReadFrom(response.Body)

    var data JsonItems
    if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
        panic(err)
    }

    inc := 0.0
    for item := range data.Data {
        fmt.Printf("%f-", data.Data[item].RelevanceScore)
	if inc <= data.Data[item].RelevanceScore {
            inc = data.Data[item].RelevanceScore
	} else {
            inc = -1
	    break
	}
    }
    fmt.Printf("\n")
    if inc != -1 {
        got = "ordered"
    }

    want := "ordered"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestSortViews(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "?sortKey=views", nil)
    response := httptest.NewRecorder()

    s := Server { }
    s.gohttp(response, request)

    got := "unordered"

    buf := new(bytes.Buffer)
    buf.ReadFrom(response.Body)

    var data JsonItems
    if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
        panic(err)
    }

    inc := 0
    for item := range data.Data {
        fmt.Printf("%d-", data.Data[item].Views)
	if inc <= data.Data[item].Views {
            inc = data.Data[item].Views
	} else {
            inc = -1
	    break
	}
    }
    fmt.Printf("\n")
    if inc != -1 {
        got = "ordered"
    }

    want := "ordered"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}


func TestLimitLow(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "?limit=1", nil)
    response := httptest.NewRecorder()

    s := Server { }
    s.gohttp(response, request)


    buf := new(bytes.Buffer)
    buf.ReadFrom(response.Body)

    var data JsonItems
    if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
        panic(err)
    }

    got := "toolow"
    fmt.Printf("Length: %d\n", len(data.Data))
    if len(data.Data) > 1 {
        got = "inrange"
    }

    want := "inrange"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestLimit(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "?limit=10", nil)
    response := httptest.NewRecorder()

    s := Server { }
    s.gohttp(response, request)


    buf := new(bytes.Buffer)
    buf.ReadFrom(response.Body)

    var data JsonItems
    if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
        panic(err)
    }

    got := "toohigh"
    fmt.Printf("Length: %d\n", len(data.Data))
    if len(data.Data) <= 10 {
        got = "inrange"
    }

    want := "inrange"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestNoLimit(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "", nil)
    response := httptest.NewRecorder()

    s := Server { }
    s.gohttp(response, request)


    buf := new(bytes.Buffer)
    buf.ReadFrom(response.Body)

    var data JsonItems
    if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
        panic(err)
    }

    got := "toohighorlow"
    fmt.Printf("Length: %d\n", len(data.Data))
    if len(data.Data) <= 200 && len(data.Data) > 1 {
        got = "inrange"
    }

    want := "inrange"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
