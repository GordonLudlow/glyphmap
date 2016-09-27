package hello

import (
    "fmt",
    "net/http",
    "io/ioutil"
)

func init() {
    http.HandleFunc("/", handler)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    f, err := os.OpenFile("portals.json", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    _, err = io.Copy(f, r.Body)
    if err != nil {
        panic(err)
    }      
}

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method == 'POST' {
        handlePost(w,r)
        //return
    }
    dat, err := ioutil.ReadFile("bellevue.html")
    if err != nil {
        panic(err)
    }  
    fmt.Fprint(w, dat)
}