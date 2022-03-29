package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	s := flag.String("s", "md5defalut", "get encode value ")
	flag.Parse()
	if *s != "md5defalut" {
		log.Println(hex.EncodeToString(md5.New().Sum([]byte(*s + "cjie"))))
		return
	}

	MailDemo.User = mustGetEnv("MAIL_USER")
	MailDemo.Pwd = mustGetEnv("MAIL_PWD")
	MailDemo.From = mustGetEnv("MAIL_NAME")
	MAIL_API_LISTEN := mustGetEnv("MAIL_API_LISTEN")
	SIGN := mustGetEnv("MAIL_SIGN")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rsign := hex.EncodeToString(
			md5.New().Sum([]byte(
				r.URL.Query().Get("sign") + "cjie"),
			),
		)
		if rsign != SIGN {
			w.WriteHeader(403)
			return
		}

		data := &Data{}
		err := json.Unmarshal([]byte(r.URL.Query().Get("data")), data)
		if err != nil {
			w.WriteHeader(500)
			log.Println("jsonunmash:", err)
			return
		}
		err = MailDemo.Set(
			strings.Split(data.To, ","), data.Sub, data.Body,
		).Send()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			log.Println("SendMail:", err)
			return
		}
		w.WriteHeader(200)
	})

	err := http.ListenAndServe(MAIL_API_LISTEN, nil)
	if err != nil {
		log.Println(err)
	}
}

func mustGetEnv(name string) string {
	res, isset := os.LookupEnv(name)
	if !isset {
		log.Println("no set env ", name)
		os.Exit(0)
	}
	return res
}

type Data struct {
	To   string
	Sub  string
	Body string
}
