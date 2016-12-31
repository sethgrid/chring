package chring

import (
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/sethgrid/chring/simpledraw"
)

func (r *Ring) drawChart(w http.ResponseWriter, req *http.Request) {
	m, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	dest := image.NewRGBA(image.Rect(0, 0, 500, 375))
	fontFolder := FindInGOPATH(filepath.Join("resources/font"))
	draw2d.SetFontFolder(fontFolder)
	gc := simpledraw.Draw{draw2dimg.NewGraphicContext(dest)}

	ring := simpledraw.NewCircle(340, 175, 150)

	// TODO parse the legend first to determine the needed canvas size, then draw all the things
	legend := &simpledraw.Legend{}
	legend.Title = "Consistent Hash Ring"
	legend.Caption = "Distribution Visualization"
	legend.Elements = make([]simpledraw.LegendElement, 0)
	gc.DrawCircle(ring)

	// draw the elements in order you want them z-stacked visually, later elements will be on top

	for i, param := range m["key[]"] {
		square := 4
		props := simpledraw.DefaultBasicProperties
		props.Color = simpledraw.Pallate[(i+3)%len(simpledraw.Pallate)]
		gc.DrawOnEdge(ring, hashAngle(r.Hasher(param)), square, 4, props)
		legend.AppendElement(square, param, props)
	}

	for i, n := range r.Nodes {
		circle := 0
		props := simpledraw.DefaultBasicProperties
		props.Color = simpledraw.Pallate[i%len(simpledraw.Pallate)]
		gc.DrawOnEdge(ring, hashAngle(n.HashID), circle, 12, props)
		legend.PrependElement(circle, n.ID, props)
	}

	for i, param := range m["hashid[]"] {
		triangle := 3
		hashID, _ := strconv.Atoi(param)
		props := simpledraw.DefaultBasicProperties
		props.Color = simpledraw.Pallate[(i+5)%len(simpledraw.Pallate)]
		gc.DrawOnEdge(ring, hashAngle(uint32(hashID)), triangle, 10, props)
		hashStr := fmt.Sprintf("hash #%d", hashID)
		legend.AppendElement(triangle, hashStr, props)
	}

	gc.DrawLegend(legend)

	w.Header().Set("Content-Type", "image/png")

	err = png.Encode(w, dest) //Encode writes the Image m to w in PNG format.
	if err != nil {
		fmt.Printf("Error rendering pie chart: %v\n", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := FindInGOPATH(filepath.Join("resources", "index.html"))

	index, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to serve index.html"))
		return
	}
	w.Write(index)
}

// Serve presents a web view into your consistent hash ring
func Serve(r *Ring, addr string) {
	http.HandleFunc("/ring.png", r.drawChart)
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// hashAngle is a helper to find the angle in radians of the hashID in uint32 space
func hashAngle(hashID uint32) float64 {
	return 2 * math.Pi * float64(hashID) / float64(math.MaxUint32)
}

// FindInGOPATH searches through all GOPATHS and attempts to find the given file
// this is useful here because we want to find chring files but we can't know the relative import path
// as the importer could be a subpackage
func FindInGOPATH(filename string) string {
	gopath := os.Getenv("GOPATH")
	paths := strings.Split(gopath, ":")
	for _, path := range paths {
		seek := filepath.Join(path, "src", "github.com", "sethgrid", "chring", filename)
		_, err := os.Stat(seek)
		if err == nil {
			return seek
		}
	}
	// not found in any GOPATH, just return what came in
	return filename
}
