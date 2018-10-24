package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println(time.Now().Format("Mon_Jan_2_15_04_05_2006"))
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/save", handler)
	http.HandleFunc("/jq.js", jq)

	log.Println("Listening... http://127.0.0.1:13000")
	http.ListenAndServe(":13000", nil)
}

func jq(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("static/jquery.js")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "text/javascript")
	w.Write(body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error read body: ", err.Error())
	}
	defer r.Body.Close()
	bbb := strings.TrimPrefix(string(body), "data:image/png;base64,")
	unbased, err := base64.StdEncoding.DecodeString(bbb)
	if err != nil {
		panic("Cannot decode b64")
	}

	b64r := bytes.NewReader(unbased)
	im, err := png.Decode(b64r)
	if err != nil {
		panic("Bad png")
	}

	fmt.Println("debug: saving:", time.Now().Format("Mon_Jan_2_15_04_05_2006")+".png")

	f, err := os.OpenFile(time.Now().Format("Mon_Jan_2_15_04_05_2006")+".png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(f, im)
	f.Close()
}
