package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "os"
    "time"
    "net/smtp"
    "github.com/ashwanthkumar/slack-go-webhook"
    "github.com/joho/godotenv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func send(text string) {
    fmt.Println("enviando correo")
    auth := smtp.PlainAuth(
        "",
        os.Getenv("EMAIL_USER"),
        os.Getenv("EMAIL_PASSWORD"),
        os.Getenv("EMAIL_HOST"),
    )
    fmt.Println(text)
    err := smtp.SendMail(
        os.Getenv("EMAIL_HOST") + ":" + os.Getenv("EMAIL_PORT"),
        auth,
        os.Getenv("EMAIL_TO"),
        []string{os.Getenv("EMAIL_FROMT")},
        []byte(text),
    )
    fmt.Println("correo enviado")
    check(err)
}

func sendSlack(text string) {
    webhookUrl := os.Getenv("SLACK_WEBHOOK")
    fmt.Println("enviando slack")

    payload := slack.Payload {
      Text: text,
      Username: os.Getenv("SLACK_USERNAME"),
      Channel: os.Getenv("SLACK_CHANNEL"),
    }

    err := slack.Send(webhookUrl, "", payload)

    fmt.Println("slack enviado")
    
    if len(err) > 0 {
      fmt.Printf("error: %s\n", err)
    }
}

func timeFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: homePage4")
    inter := r.URL.Query().Get("inter")

    if inter == "" && inter != "eth0" && inter != "eth2" {
        fmt.Fprintf(w, "False")
        return
    }

    info, err := os.Stat("/tmp/" + inter)

    if err != nil {
        send(err.Error())
        sendSlack(err.Error())
        fmt.Fprintf(w, "False")
        return
    }

    duration := time.Since(info.ModTime())
    fmt.Println(duration.Minutes())

    if duration.Minutes() > 4 {
        send("Error en interfaz " + inter)
        sendSlack("Error en interfaz " + inter)
        fmt.Fprintf(w, "False")
        return
    }

    fmt.Fprintf(w, "True")
    return
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: homePage3")
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
    fmt.Println("Endpoint Hit: homePage2")
    http.HandleFunc("/", homePage)
    http.HandleFunc("/timeFile", timeFile)
    log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	fmt.Println(os.Getenv("SLACK_WEBHOOK"))
	err := godotenv.Load()
	check(err)
    fmt.Println("Endpoint Hit: homePage1")
    handleRequests()
}
