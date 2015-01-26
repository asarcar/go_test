package main

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"flag"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Example: Google I/O 2011 Go.
//
// Objective: Web Server: Clients can upload images. Each image
// is then added a moustache.
//
// Data Flow:
// 1a. Client ===============> GET / ==============> Server
// 1b. Client <=========== <html><form>... ========= Server
// 2a. Client =========> POST / (with image) ======> Server
// 2b. Client <====== Redirect /edit?id=ID ========= Server
// 3a. Client =========> GET /edit?id=ID ==========> Server
// 3b. Client <== <html><img src="/img?id=ID"> ===== Server
//                       |<----- Img Def ----->|
// 4a. Client == GET /img?id=ID&x=X&y=Y&s=S&d=D ===> Server
// e.g: localhost:4000/img?id=173604969&x=74000&y=140000&s=2000&d=5
//      This occurs each time the moustache is adjusted
// 4b. Client <== PNG img with moustache drawn  ==== Server
//
// PHASES:
// Upload: 1 and 2
// Edit  : 3
// Img   : 4

var gUploadTemplate, gErrorTemplate, gEditTemplate *template.Template
var gTmpDir string

func main() {
	setGlobals()
	setHandleFuncs()
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}

func setGlobals() {
	htmlDir, tmpDir := parseFlags()

	gTmpDir = tmpDir
	gUploadTemplate = getTemplate(htmlDir, "upload.html")
	gErrorTemplate = getTemplate(htmlDir, "error.html")
	gEditTemplate = getTemplate(htmlDir, "edit.html")
}

func setHandleFuncs() {
	http.HandleFunc("/", errorHandler(upload))
	http.HandleFunc("/view", errorHandler(view))
	http.HandleFunc("/edit", errorHandler(edit))
	http.HandleFunc("/img", errorHandler(img))
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		gUploadTemplate.Execute(w, nil)
		return
	}

	// POST the image file
	f, _, err := r.FormFile("image")
	checkError(err)
	defer f.Close()

	t, err := ioutil.TempFile(gTmpDir, "image-")
	checkError(err)
	defer t.Close()

	_, er := io.Copy(t, f)
	checkError(er)
	http.Redirect(w, r, "/edit?id="+t.Name()[len(gTmpDir+"image-"):], 302)
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	imgFile := gTmpDir + "image-" + r.FormValue("id")
	http.ServeFile(w, r, imgFile)
}

func edit(w http.ResponseWriter, r *http.Request) {
	gEditTemplate.Execute(w, r.FormValue("id"))
}

func img(w http.ResponseWriter, r *http.Request) {
	imgFile := gTmpDir + "image-" + r.FormValue("id")
	f, err := os.Open(imgFile)
	checkError(err)
	m, _, err := image.Decode(f)
	checkError(err)

	v := func(s string) int { // helper closure
		i, _ := strconv.Atoi(r.FormValue(s))
		return i
	}

	m = moustache(m, v("x"), v("y"), v("s"), v("d"))

	w.Header().Set("Content-type", "image/jpeg")
	jpeg.Encode(w, m, nil) // Default JPEG options.
}

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Type Assertion: assert that recover() is not nil
			// and value stored is of type os.Error
			// Since os.Error is an interface this call asserts
			// that the dynamic type of recover implements os.Error
			// This form ensures ok is "true" if type assertion holds
			if e, ok := recover().(error); ok {
				w.WriteHeader(500)
				gErrorTemplate.Execute(w, e)
			}
		}()
		fn(w, r)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseFlags() (htmlDirPath string, tmpDirPath string) {
	flag.Parse()
	dPtr := flag.String("d",
		"/home/asarcar/git/go_test/src/github.com/asarcar/go_test/moustachio/html/",
		"full path to directory where template html files exit\n")
	tPtr := flag.String("t",
		"/home/asarcar/tmp/",
		"full path for temporay directory where image file would be saved\n")

	return *dPtr, *tPtr
}

func getTemplate(dirPath, fileName string) *template.Template {
	f := dirPath + fileName
	t, err := template.ParseFiles(f)
	if err != nil {
		panic(err)
	}

	return t
}

func moustache(m image.Image, x, y, size, droopFactor int) image.Image {
	mrgba := rgba(m) // Create specialized RGBA image from m

	p := raster.NewRGBAPainter(mrgba)
	p.SetColor(color.RGBA{0, 0, 0, 255})

	w, h := m.Bounds().Dx(), m.Bounds().Dy()
	r := raster.NewRasterizer(w, h)

	var (
		mag   = raster.Fix32((10 + size) << 8)
		width = raster.Point{20, 0}.Mul(mag)
		mid   = raster.Point{raster.Fix32(x), raster.Fix32(y)}
		droop = raster.Point{0, raster.Fix32(droopFactor)}.Mul(mag)
		left  = mid.Sub(width).Add(droop)
		right = mid.Add(width).Add(droop)
		bow   = raster.Point{0, 5}.Mul(mag).Sub(droop)
		curlx = raster.Point{10, 0}.Mul(mag)
		curly = raster.Point{0, 2}.Mul(mag)
		risex = raster.Point{2, 0}.Mul(mag)
		risey = raster.Point{0, 5}.Mul(mag)
	)

	r.Start(left)
	r.Add3(
		mid.Sub(curlx).Add(curly),
		mid.Sub(risex).Sub(risey),
		mid,
	)
	r.Add3(
		mid.Add(risex).Sub(risey),
		mid.Add(curlx).Add(curly),
		right,
	)
	r.Add2(
		mid.Add(bow),
		left,
	)
	r.Rasterize(p)

	return mrgba
}

func rgba(m image.Image) *image.RGBA {
	// Fast path: if m is already an RGBA, just return it
	if r, ok := m.(*image.RGBA); ok { // Type Assertion
		return r
	}

	// Create a new image and draw into it
	b := m.Bounds()
	r := image.NewRGBA(b)

	draw.Draw(r, b, m, image.ZP, draw.Over)

	return r
}
