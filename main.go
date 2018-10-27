package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var basePath string

func main() {
	var err error
	basePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("老铁有点错误，请告诉philo： ", err.Error())
	}
	fmt.Println("dir: ", basePath)
	fmt.Println(time.Now().Format("Mon_Jan_2_15_04_05_2006"))
	fs := http.FileServer(http.Dir(basePath + "/static"))
	http.Handle("/", fs)
	http.HandleFunc("/save", handler)
	http.HandleFunc("/saveData", handlerData)
	http.HandleFunc("/jq.js", jq)

	log.Println("Listening... http://127.0.0.1:13000")
	http.ListenAndServe(":13000", nil)
}

type Position struct {
	X int
	Y int
}

func jq(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile(basePath + "/static/jquery.js")
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

	f, err := os.OpenFile(time.Now().Format(basePath+"/png/Mon_Jan_2_15_04_05_2006")+".png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(f, im)
	f.Close()
}

func handlerData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error read body: ", err.Error())
	}
	defer r.Body.Close()
	var data []int
	json.Unmarshal(body, &data)
	f, err := os.OpenFile(time.Now().Format(basePath+"/csv/Mon_Jan_2_15_04_05_2006")+".csv", os.O_WRONLY|os.O_CREATE, 0777)
	f.Write([]byte("X,Y\n"))
	var ps []Position
	for idx := range data {
		if idx%2 == 0 {
			ps = append(ps, Position{
				X: data[idx],
				Y: data[idx+1],
			})
			f.Write([]byte(fmt.Sprintf("%d,%d\n", data[idx], data[idx+1])))
		}
	}
	f.Close()
	outputSVG(ps)
	w.Write([]byte("OK"))
}

func outputSVG(ps []Position) {
	if len(ps) < 2 {
		fmt.Println("not enough data to draw a picture")
		return
	}
	pathD := ""
	for idx, p := range ps {
		if idx == 0 {
			pathD += fmt.Sprintf("M %d,%d\n", p.X, p.Y)
		} else {
			pathD += fmt.Sprintf("Q %f,%f %d,%d\n", float64(ps[idx-1].X+p.X)/2.0, float64(ps[idx-1].Y+p.Y)/2.0, p.X, p.Y)
		}
	}
	svgBody := `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" height="500" width="1105">
  <path fill="none" stroke="red"
        d="` + pathD + `
           " />
  <!-- Mark relevant points -->
  <g stroke="black" stroke-width="3" fill="black">
    <circle id="pointA" cx="` + strconv.Itoa(ps[0].X) + `" cy="` + strconv.Itoa(ps[0].Y) + `" r="3" />
    <circle id="pointB" cx="` + strconv.Itoa(ps[len(ps)-1].X) + `" cy="` + strconv.Itoa(ps[len(ps)-1].Y) + `" r="3" />
  </g>
  <!-- Label the points -->
  <g font-size="30" font-family="sans-serif" fill="black" stroke="none" text-anchor="middle">
    <text x="` + strconv.Itoa(ps[0].X) + `" y="` + strconv.Itoa(ps[0].Y) + `" dx="-30">Start</text>
    <text x="` + strconv.Itoa(ps[len(ps)-1].X) + `" y="` + strconv.Itoa(ps[len(ps)-1].Y) + `" dx="30">End</text>
  </g>
</svg>
`
	f, err := os.OpenFile(time.Now().Format(basePath+"/svg/Mon_Jan_2_15_04_05_2006")+".svg", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("error save svg: ", err.Error())
	}
	f.Write([]byte(svgBody))
	f.Close()
}
