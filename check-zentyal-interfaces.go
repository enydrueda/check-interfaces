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
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func sendMail(text string) {
    auth := smtp.PlainAuth(
        "",
        os.Getenv("EMAIL_USER"),
        os.Getenv("EMAIL_PASSWORD"),
        os.Getenv("EMAIL_HOST"),
    )

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

func interfaceInEnv(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func timeFile(w http.ResponseWriter, r *http.Request) {
    inter := r.URL.Query().Get("inter")

	if !interfaceInEnv(inter, strings.Split(os.Getenv("INTERFACES"), ",")) {
        return
	}

    info, err := os.Stat("/tmp/" + inter)

    if err != nil {
        sendMail(err.Error())
        sendSlack(err.Error())
        return
    }

    duration := time.Since(info.ModTime())

    if duration.Minutes() > 4 {
        sendMail("Error en interfaz " + inter)
        sendSlack("Error en interfaz " + inter)
        return
    }

    return
}

func homePage(w http.ResponseWriter, r *http.Request) {
    inter := r.URL.Query().Get("inter")

    if !interfaceInEnv(inter, strings.Split(os.Getenv("INTERFACES"), ",")) {
        return
	}

    d1 := []byte("ready")
    err := ioutil.WriteFile("/tmp/" + inter, d1, 0644)
    check(err)

    return
}

func handleRequests() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/timeFile", timeFile)
    log.Fatal(http.ListenAndServe(":" + os.Getenv("SERVER_PORT"), nil))
}

func main() {
	err := godotenv.Load()
	check(err)

    fmt.Println("Run")
    handleRequests()
}