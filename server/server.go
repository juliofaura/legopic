// Webpage management for Casheth

package server

import (
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/context"

	"github.com/juliofaura/legopic/process"
)

///////////////////////////////////////////////////
// Constants (some in fact defined as global variables) and types
///////////////////////////////////////////////////

const (
	WEB_PATH           = "./web/"
	Color0             = 0
	Color1             = 80
	Color2             = 160
	Color3             = 250
	Resolution_default = 15
	Thres1_default     = 120
	Thres2_default     = 150
	Thres3_default     = 180
)

var (
	WEBPORT string = "8050"
)

var templates = template.Must(template.ParseFiles(
	WEB_PATH + "index.html",
))

func StartWeb() {
	http.Handle("/", http.HandlerFunc(HandleRoot))
	http.Handle("/upload", http.HandlerFunc(HandleUpload))
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir(WEB_PATH+"pics"))))
	go func() {
		addr := flag.String("addr", ":"+WEBPORT, "http service address")
		err := http.ListenAndServe(*addr, context.ClearHandler(http.DefaultServeMux))
		if err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func HandleRoot(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var opt jpeg.Options

	var resolution int
	resolutionA, ok := req.Form["resolution"]
	if ok {
		resolution64, err := strconv.ParseInt(resolutionA[0], 10, 64)
		if err != nil {
			resolution = Resolution_default
		} else {
			resolution = int(resolution64)
		}
	} else {
		resolution = Resolution_default
	}
	log.Println("Resolution set to", resolution)

	var thres1, thres2, thres3 int64
	var err error
	thres1A, ok := req.Form["thres1"]
	if ok {
		thres1, err = strconv.ParseInt(thres1A[0], 10, 64)
		if err != nil {
			thres1 = Thres1_default
		}
	} else {
		thres1 = Thres1_default
	}
	thres2A, ok := req.Form["thres2"]
	if ok {
		thres2, err = strconv.ParseInt(thres2A[0], 10, 64)
		if err != nil {
			thres2 = Thres2_default
		}
	} else {
		thres2 = Thres2_default
	}
	thres3A, ok := req.Form["thres3"]
	if ok {
		thres3, err = strconv.ParseInt(thres3A[0], 10, 64)
		if err != nil {
			thres3 = Thres3_default
		}
	} else {
		thres3 = Thres3_default
	}
	log.Println("Resolution set to", resolution)

	var photoID string
	photoIDA, ok := req.Form["photoID"]
	if ok {
		photoID = photoIDA[0]
	} else {
		photoID = "default"
	}

	imgfile, err := os.Open(WEB_PATH + "pics/" + photoID + "_to_process.jpg")
	if err != nil {
		log.Fatal("pics/"+photoID+"_to_process.jpg", "file not found!")
	}
	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	//fmt.Println(img.At(10, 10))
	//bounds := img.Bounds()
	//fmt.Println(bounds)
	//canvas := image.NewAlpha(bounds)
	// is this image opaque
	//op := canvas.Opaque()
	//fmt.Println(op)
	log.Println("Image decoded")

	img_bw := process.TurnToBW(img)
	log.Println("B&W image created")

	/*
		out_bw, err := os.Create(WEB_PATH + "pics/bw.jpg")
		if err != nil {
			log.Fatalln(err)
		}
		opt.Quality = 80
		err = jpeg.Encode(out_bw, img_bw, &opt) // put quality to 80%
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("B&W image written to disk")
	*/

	out_proc, err := os.Create(WEB_PATH + "pics/processed.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	img_processed := image.NewGray(img.Bounds())
	log.Println("New image created")

	/*
		Rs := make([][]int64, 1+img.Bounds().Dx()/resolution)
		Gs := make([][]int64, 1+img.Bounds().Dx()/resolution)
		Bs := make([][]int64, 1+img.Bounds().Dx()/resolution)
		As := make([][]int64, 1+img.Bounds().Dx()/resolution)
		for x := 0; x <= img.Bounds().Dx()/resolution; x++ {
			Rs[x] = make([]int64, 1+img.Bounds().Dy()/resolution)
			Gs[x] = make([]int64, 1+img.Bounds().Dy()/resolution)
			Bs[x] = make([]int64, 1+img.Bounds().Dy()/resolution)
			As[x] = make([]int64, 1+img.Bounds().Dy()/resolution)
		}
		log.Println("pixels matrix created")
	*/

	Is := make([][]int64, 1+img.Bounds().Dx()/resolution)
	for x := 0; x <= img.Bounds().Dx()/resolution; x++ {
		Is[x] = make([]int64, 1+img.Bounds().Dy()/resolution)
	}
	log.Println("pixels matrix created")

	for x := 0; x < img_bw.Bounds().Dx(); x++ {
		for y := 0; y < img_bw.Bounds().Dy(); y++ {
			/*
				r, g, b, a := img_bw.At(x, y).RGBA()
				Rs[x/resolution][y/resolution] += int64(r)
				Gs[x/resolution][y/resolution] += int64(g)
				Bs[x/resolution][y/resolution] += int64(b)
				As[x/resolution][y/resolution] += int64(a)
			*/
			Is[x/resolution][y/resolution] += int64(img_bw.GrayAt(x, y).Y)
		}
	}
	for x := 0; x < img.Bounds().Dx()/resolution; x++ {
		for y := 0; y < img.Bounds().Dy()/resolution; y++ {
			/*
				Rs[x][y] /= int64(resolution)
				Rs[x][y] /= int64(resolution)
				Gs[x][y] /= int64(resolution)
				Gs[x][y] /= int64(resolution)
				Bs[x][y] /= int64(resolution)
				Bs[x][y] /= int64(resolution)
				As[x][y] /= int64(resolution)
				As[x][y] /= int64(resolution)
			*/
			Is[x][y] /= int64(resolution)
			Is[x][y] /= int64(resolution)
			switch {
			case Is[x][y] < thres1:
				Is[x][y] = Color0
			case Is[x][y] < thres2:
				Is[x][y] = Color1
			case Is[x][y] < thres3:
				Is[x][y] = Color2
			default:
				Is[x][y] = Color3
			}
		}
	}
	log.Println("Color components accumulated and averaged")

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			/*
				var thisPix color.RGBA
				thisPix.R = uint8(Rs[x/resolution][y/resolution])
				thisPix.G = uint8(Gs[x/resolution][y/resolution])
				thisPix.B = uint8(Bs[x/resolution][y/resolution])
				thisPix.A = uint8(As[x/resolution][y/resolution])
			*/
			var thisPix color.Gray
			thisPix.Y = uint8(Is[x/resolution][y/resolution])
			img_processed.Set(x, y, thisPix)
		}
	}
	log.Println("New image processed")

	opt.Quality = 80
	err = jpeg.Encode(out_proc, img_processed, &opt) // put quality to 80%
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("New image written to disk")

	passdata := map[string]interface{}{
		"sampledata": "sample data",
		"resolution": resolution,
		"thres1":     thres1,
		"thres2":     thres2,
		"thres3":     thres3,
		"dimensions": fmt.Sprintf("%v x %v, or %.2f m x %.2f m",
			img_processed.Bounds().Dx()/resolution,
			img_processed.Bounds().Dy()/resolution,
			0.008*float32(img_processed.Bounds().Dx()/resolution),
			0.008*float32(img_processed.Bounds().Dy()/resolution),
		),
	}
	templates.ExecuteTemplate(w, "index.html", passdata)
}

func HandleUpload(w http.ResponseWriter, req *http.Request) {
	log.Println("Hello from Upload")
	req.ParseMultipartForm(32 << 20)
	//file, handler, err := req.FormFile("uploadfile")
	file, _, err := req.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	//fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(WEB_PATH+"pics/to_process.jpg", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	log.Println("File uploaded")
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
