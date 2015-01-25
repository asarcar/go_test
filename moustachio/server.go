package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
// 2b. Client <====== Redirect /view?id=ID ========= Server
// 3a. Client =========> GET /view?id=ID ==========> Server
// 3b. Clinet <========== Image Data =============== Server

var gUploadTemplate, gErrorTemplate *template.Template
var gHtmlDir, gTmpDir string

func main() {
	gHtmlDir, gTmpDir = parseFlags()
	gUploadTemplate = getTemplate(gHtmlDir, "upload.html")
	gErrorTemplate = getTemplate(gHtmlDir, "error.html")

	http.HandleFunc("/", errorHandler(upload))
	http.HandleFunc("/view", errorHandler(view))
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
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
	http.Redirect(w, r, "/view?id="+t.Name()[len(gTmpDir+"image-"):], 302)
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	imgFileName := gTmpDir + "image-" + r.FormValue("id")
	http.ServeFile(w, r, imgFileName)
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
				// http.Error(w, e.Error(), 500)
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
