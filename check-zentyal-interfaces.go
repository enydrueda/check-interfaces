package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "os"
    "time"
    "net/smtp"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func send(text string) {
    auth := smtp.PlainAuth(
        "",
        "<USER>",
        "<PASSWORD>",
        "<HOST>",
    )
    
    err := smtp.SendMail(
        "<HOST:PORT>",
        auth,
        "<FROM>",
        []string{"<RCPT>"},
        []byte(text),
    )

    check(err)
}

func timeFile(w http.ResponseWriter, r *http.Request) {
    inter := r.URL.Query().Get("inter")

    if inter == "" && inter != "eth0" && inter != "eth2" {
        fmt.Fprintf(w, "False")
        return
    }

    info, err := os.Stat("/tmp/" + inter)

    if err != nil {
        send(err.Error())
        fmt.Fprintf(w, "False")
        return
    }

    duration := time.Since(info.ModTime())

    if duration.Minutes() > 4 {
        send("Error en interfaz " + inter)
        fmt.Fprintf(w, "False")
        return
    }

    fmt.Fprintf(w, "True")
    return
}

func homePage(w http.ResponseWriter, r *http.Request) {
    inter := r.URL.Query().Get("inter")

    if inter == "" && inter != "eth0" && inter != "eth2" {
        fmt.Fprintf(w, "False")
        return
    }

    d1 := []byte("ready")
    err := ioutil.WriteFile("/tmp/" + inter, d1, 0644)
    check(err)

    fmt.Fprintf(w, "True")
    return
}

func handleRequests() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/timeFile", timeFile)
    log.Fatal(http.ListenAndServe("<:PORT>", nil))
}

func main() {
    handleRequests()
}
