package main

import (
	"bufio"
	"errors"
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
	"strings"

	"github.com/golang/freetype/raster"
	"google.golang.org/api/drive/v2"
)

// Example: Google I/O 2011 Go.
//
// Objective: Web Server: Clients can upload images. Each image
// is then added a moustache.
//
// Data Flow:
// 1a. Client ===============> GET / ==============> Moustachio
// 1b. Client <=========== <html><form>... ========= Moustachio
// 2a. Client =========> POST / (with image) ======> Moustachio
// 2b. Client <====== Redirect /edit?id=ID ========= Moustachio
// 3a. Client =========> GET /edit?id=ID ==========> Moustachio
// 3b. Client <== <html><img src="/img?id=ID"> ===== Moustachio
//                       |<----- Img Def ----->|
// 4a. Client == GET /img?id=ID&x=X&y=Y&s=S&d=D ===> Moustachio
// e.g: localhost:4000/img?id=173604969&x=74000&y=140000&s=2000&d=5
//      This occurs each time the moustache is adjusted
// 4b. Client <==== IMG with moustache drawn  ====== Moustachio
//      client triggers share when happy with image
// 5.  Client = GET /share?id=ID&x=X&y=Y&s=S&d=D  => Moustachio
// 6a. Client <== Redirect to OAuth svc with img === Moustachio
// 6b. Client <== User auths & grants access to Moustachio ============> OAuthSvc
// 6c. Client <== Redirect to Moustachio: /post?code=CODE&state=IMAGE == OAuthSvc
// 7a. Client === GET /post?code=CODE&state=IMAGE => Moustachio
// 7w.           Moustachio ================ CODE =====================> OAuthSvc
// 7x.           Moustachio <============== Token ====================== OAuthSvc
// 7y.           Moustachio ============ POST Buzz Token ==============> GoogleBuzz
// 7b. Client <=<html> Your image has been shared!=> Moustachio=200 OK=> GoogleBuzz
//
// PHASES:
// -------
// 1 & 2 : Upload: 1 and 2
// 3     : Edit
// 4     : Image
// 5     : Share
// 6     : Authenticate
// 7     : Post

const (
	kListenAddr     = "localhost:4000"
	kUploadFileName = "upload.html"
	kErrorFileName  = "error.html"
	kEditFileName   = "edit.html"
	kPostFileName   = "post.html"
	kMsgFileName    = "msg.html"
)

var gUploadTemplate, gErrorTemplate *template.Template
var gEditTemplate, gPostTemplate, gMsgTemplate *template.Template
var gTmpDir string

func main() {
	setGlobals()
	setHandleFuncs()
	log.Fatal(http.ListenAndServe(kListenAddr, nil))
}

func setGlobals() {
	htmlDir, tmpDir := parseFlags()

	gTmpDir = tmpDir
	gUploadTemplate = getTemplate(htmlDir, kUploadFileName)
	gErrorTemplate = getTemplate(htmlDir, kErrorFileName)
	gEditTemplate = getTemplate(htmlDir, kEditFileName)
	gPostTemplate = getTemplate(htmlDir, kPostFileName)
	gMsgTemplate = getTemplate(htmlDir, kMsgFileName)
}

func setHandleFuncs() {
	http.HandleFunc("/", errorHandler(upload))
	http.HandleFunc("/edit", errorHandler(edit))
	http.HandleFunc("/img", errorHandler(img))
	http.HandleFunc("/post", errorHandler(post))
	http.HandleFunc("/share", errorHandler(share))
	http.HandleFunc("/view", errorHandler(view))
	http.HandleFunc("/fetch", errorHandler(fetch))
	http.HandleFunc("/display", errorHandler(display))
}

// 1a. Client ===============> GET / ==============> Moustachio
// 1b. Client <=========== <html><form>... ========= Moustachio
// 2a. Client =========> POST / (with image) ======> Moustachio
// 2b. Client <====== Redirect /edit?id=ID ========= Moustachio
func upload(w http.ResponseWriter, r *http.Request) {
	// Executes response to GET query via upload.html
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
	// Redirect for image editing
	http.Redirect(w, r,
		"/edit?id="+t.Name()[len(gTmpDir+"image-"):],
		http.StatusFound) // 302
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	imgFile := gTmpDir + "image-" + r.FormValue("id")
	http.ServeFile(w, r, imgFile)
}

// 3a. Client =========> GET /edit?id=ID ==========> Moustachio
// 3b. Client <== <html><img src="/img?id=ID"> ===== Moustachio
func edit(w http.ResponseWriter, r *http.Request) {
	gEditTemplate.Execute(w, r.FormValue("id"))
}

//	|<----- Img Def ----->|
//
// 4a. Client == GET /img?id=ID&x=X&y=Y&s=S&d=D ===> Moustachio
// e.g: localhost:4000/img?id=173604969&x=74000&y=140000&s=2000&d=5
//
//	This occurs each time the moustache is adjusted
//
// 4b. Client <==== IMG with moustache drawn  ====== Moustachio
func img(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/jpeg")

	iP := ImgParams{
		r.FormValue("id"), r.FormValue("x"), r.FormValue("y"),
		r.FormValue("s"), r.FormValue("d"),
	}

	err := imgWrite(w, iP)
	checkError(err)
}

// AuthURL: Target of initial request. Handles active session lookup,
//
//	authenticating the user and user consent. Result include
//	access tokens, refresh tokens, and authorization codes.
var config = &oauth.Config{
	ClientId:     "311569668069-qith1651t1jgh9e7qck9n9mhv7td37ug.apps.googleusercontent.com",
	ClientSecret: "PiSJX4iqm0COBYXr7DiVRW3j",
	Scope:        "https://www.googleapis.com/auth/drive",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

// Redirect user to Google's OAuth service
// 5.  Client = GET /share?id=ID&x=X&y=Y&s=S&d=D  => Moustachio
// 6a. Client <== Redirect to OAuth svc with img === Moustachio
// 6b. Client <== User auths & grants access to Moustachio ============> OAuthSvc
// 6c. Client <== Redirect to Moustachio: /post?code=CODE&state=IMAGE == OAuthSvc
func share(w http.ResponseWriter, r *http.Request) {
	config.RedirectURL = "http://localhost:4000/post"
	// Raw Query is our Img Def: id=ID&x=X&y=Y&s=S&d=D
	url := config.AuthCodeURL(r.URL.RawQuery)
	http.Redirect(w, r, url, http.StatusFound) // 302
}

type FileParams struct {
	ID, Title, Description, MimeType string
}

// post handler: accepts authentication code and image state from the OAuth
// service, exchanges the code for an OAut authentication token, and
// posts to Buzz
//
// 7a. Client === GET /post?code=CODE&state=IMAGE => Moustachio
// 7w.           Moustachio ================ CODE =====================> OAuthSvc
// 7x.           Moustachio <============== Token ====================== OAuthSvc
// 7y.           Moustachio ============ POST Buzz Token ==============> GoogleBuzz
// 7b. Client <=<html> Your image has been shared!=> Moustachio=200 OK=> GoogleBuzz
func post(w http.ResponseWriter, r *http.Request) {
	var t *oauth.Transport = &oauth.Transport{Config: config}

	// 7w. Parses code from GET query and submits to OAuthSvc
	code := r.FormValue("code")
	_, err := t.Exchange(code)
	checkError(err)

	imgId := r.FormValue("state")
	imgFileName := "img-" + imgId
	imgFullFileName := gTmpDir + imgFileName
	writeImageFile(imgId, imgFullFileName)

	// Create a new authorized Drive client
	var auth_client *http.Client = t.Client()
	svc, err := drive.New(auth_client)
	checkError(err)

	// Define the metadata for the file we are uploading
	f := &drive.File{
		Title:       imgFileName,
		Description: "Moustachio of image: " + imgFileName,
	}
	// Upload the image file we created to GDrive
	m, err := os.Open(imgFullFileName)
	checkError(err)
	defer m.Close()

	dF, err := svc.Files.Insert(f).Media(m).Do()
	checkError(err)

	// 7b. Displays Msg "Your image has been shared!"
	gPostTemplate.Execute(w, FileParams{
		ID:          dF.Id,
		Title:       dF.Title,
		Description: dF.Description,
		MimeType:    dF.MimeType})
}

// * User authenticates with GDrive
// * Redirect to display routine with CODE
func fetch(w http.ResponseWriter, r *http.Request) {
	config.RedirectURL = "http://localhost:4000/display"
	// Raw Query is our Img Def: id=GDriveFileID
	url := config.AuthCodeURL(r.URL.RawQuery)
	http.Redirect(w, r, url, http.StatusFound) // 302
}

// * Use CODE to get access token
// * Use Token to read file
// * Display file to User
func display(w http.ResponseWriter, r *http.Request) {
	var t *oauth.Transport = &oauth.Transport{Config: config}

	code := r.FormValue("code")
	_, err := t.Exchange(code)
	checkError(err)

	idQuery := r.FormValue("state")
	fileId, err := func(qstr string) (string, error) {
		if !strings.HasPrefix(qstr, "id=") {
			return "", errors.New("Malformed Query: " + idQuery + " no 'id=' prefix")
		}
		return strings.TrimPrefix(qstr, "id="), nil
	}(idQuery)

	// Get the current file to validate its existence and presence of downloadURL
	svc, err := drive.New(t.Client())
	dF, err := svc.Files.Get(fileId).Do()
	checkError(err)

	if dF.DownloadUrl == "" {
		w.WriteHeader(http.StatusNoContent)
		gMsgTemplate.Execute(w, "File "+dF.Title+": Content not downloadbable")
		return
	}

	// Fetch the body content of the File and display the image
	req, err := http.NewRequest("GET", dF.DownloadUrl, nil)
	checkError(err)

	resp, err := t.RoundTrip(req)
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "image")
	_, er := io.Copy(w, resp.Body)
	checkError(er)
}

func writeImageFile(imgId, imgFullFileName string) {
	imgFile, err := os.Create(imgFullFileName)
	checkError(err)
	defer imgFile.Close()

	strs := strings.Split(imgId, "&") // %26 == '&'
	// Parse the id, x, y, s, and d string from imgId
	fV := func(varstr string) (str string, err error) {
		for _, str := range strs {
			if !strings.HasPrefix(str, varstr) {
				continue
			}
			return strings.TrimPrefix(str, varstr), nil
		}
		return "", errors.New("String: " + varstr + " not found in " + imgId)
	}

	idStr, err := fV("id=")
	checkError(err)
	xStr, err := fV("x=") // %3D == '='
	checkError(err)
	yStr, err := fV("y=")
	checkError(err)
	sStr, err := fV("s=")
	checkError(err)
	dStr, err := fV("d=")
	checkError(err)

	wImgFile := bufio.NewWriter(imgFile)
	defer wImgFile.Flush()
	er := imgWrite(wImgFile, ImgParams{idStr, xStr, yStr, sStr, dStr})
	checkError(er)
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
				w.WriteHeader(http.StatusInternalServerError)
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
	dPtr := flag.String("d",
		"/home/asarcar/git/go_test/src/github.com/asarcar/go_test/moustachio/html/",
		"full path to directory where template html files exist\n")
	tPtr := flag.String("t",
		"/home/asarcar/tmp/",
		"full path for temporay directory where image file would be saved\n")
	flag.Parse()

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

type ImgParams struct {
	Id, X, Y, S, D string
}

func imgWrite(w io.Writer, iP ImgParams) error {
	v := func(s string) int { // helper closure
		i, err := strconv.Atoi(s)
		checkError(err)
		return i
	}

	imgFile := gTmpDir + "image-" + iP.Id
	f, err := os.Open(imgFile)
	checkError(err)
	m, _, err := image.Decode(f)
	checkError(err)

	// draw when the user clicks somewhere inside the picture
	x, y := v(iP.X), v(iP.Y)
	if x > 0 || y > 0 {
		m = moustache(m, x, y, v(iP.S), v(iP.D))
	}

	return jpeg.Encode(w, m, nil) // Default JPEG options.
}

func moustache(m image.Image, x, y, size, droopFactor int) image.Image {
	mrgba := rgba(m) // Create specialized RGBA image from m

	p := raster.NewRGBAPainter(mrgba)
	p.SetColor(color.RGBA{0, 0, 0, 255})

	w, h := m.Bounds().Dx(), m.Bounds().Dy()
	r := raster.NewRasterizer(w, h)

	var (
		mag   = raster.Fix32((10 + size) << 8)
		width = pt(20, 0).Mul(mag)
		mid   = pt(x, y)
		droop = pt(0, droopFactor).Mul(mag)
		left  = mid.Sub(width).Add(droop)
		right = mid.Add(width).Add(droop)
		bow   = pt(0, 5).Mul(mag).Sub(droop)
		curlx = pt(10, 0).Mul(mag)
		curly = pt(0, 2).Mul(mag)
		risex = pt(2, 0).Mul(mag)
		risey = pt(0, 5).Mul(mag)
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

// pt returns the raster.Point corresponding to the pixel position (x, y).
func pt(x, y int) raster.Point {
	return raster.Point{X: raster.Fix32(x << 8), Y: raster.Fix32(y << 8)}
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
