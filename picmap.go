package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io/ioutil"
	"log"
	"os"
	"flag"
	"html/template"
	"net/http"
)

var addr = flag.String("addr", ":1718", "http service address")

var templ = template.Must(template.New("qr").Parse(templateStr))


func main() {
	
	

	fname := "sample1.jpg"

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	camModel, _ := x.Get(exif.Model)
	date, _ := x.Get(exif.DateTimeOriginal)
	fmt.Println(camModel.StringVal())
	fmt.Println(date.StringVal())

	focal, _ := x.Get(exif.FocalLength)
	numer, denom := focal.Rat2(0) // retrieve first (only) rat. value
	fmt.Printf("%v/%v", numer, denom)
	fmt.Println()

	gps, _ := x.Get(exif.GPSLatitude)
	fmt.Println(gps)
	//fmt.Println(x.String())

	files, err := ioutil.ReadDir("./")

	for _, f := range files {
            fmt.Println(f.Name())
    }

    flag.Parse()
    http.Handle("/", http.HandlerFunc(showPics))
    err = http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
    
}

func showPics(w http.ResponseWriter, req *http.Request) {
    templ.Execute(w, req.FormValue("s"))
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`