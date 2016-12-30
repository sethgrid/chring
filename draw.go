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

	var x, y float64
	ring := simpledraw.NewCircle(340, 175, 150)

	// TODO parse the legend first to determine the needed canvas size, then draw all the things
	legend := &simpledraw.Legend{}
	legend.Title = "Consistent Hash Ring"
	legend.Caption = "Distribution Visualization"
	legend.Elements = make([]simpledraw.LegendElement, 0)

	gc.DrawCircle(ring)
	for i, n := range r.Nodes {
		x, y = ring.PointAtAngle(hashAngle(n.HashID))
		c := simpledraw.NewCircle(x, y, 10)
		c.Props.Color = simpledraw.Pallate[i%len(simpledraw.Pallate)]
		gc.DrawCircle(c)
		left, top, right, bottom := gc.GetStringBounds(n.ID)
		legend.Elements = append(legend.Elements, simpledraw.LegendElement{
			IsCircle: true,
			Name:     n.ID,
			Props:    c.Props,
			Width:    right - left,
			Height:   bottom - top,
		})
	}

	for i, param := range m["key[]"] {
		hashID := r.Hasher(param)
		x, y = ring.PointAtAngle(hashAngle(hashID))
		s := simpledraw.NewSquare(x, y, 8)
		s.Props.Color = simpledraw.Pallate[i%len(simpledraw.Pallate)+3]
		gc.DrawSquare(s)
		left, top, right, bottom := gc.GetStringBounds(param)
		legend.Elements = append(legend.Elements, simpledraw.LegendElement{
			IsSquare: true,
			Name:     param,
			Props:    s.Props,
			Width:    right - left,
			Height:   bottom - top,
		})
	}

	for i, param := range m["hashid[]"] {
		hashID, _ := strconv.Atoi(param)
		x, y = ring.PointAtAngle(hashAngle(uint32(hashID)))
		t := simpledraw.NewTriangle(x, y, 8)
		t.Props.Color = simpledraw.Pallate[i%len(simpledraw.Pallate)+5]
		gc.DrawTriangle(t)
		hashStr := fmt.Sprintf("hash #%d", hashID)
		left, top, right, bottom := gc.GetStringBounds(hashStr)
		legend.Elements = append(legend.Elements, simpledraw.LegendElement{
			IsTriangle: true,
			Name:       hashStr,
			Props:      t.Props,
			Width:      right - left,
			Height:     bottom - top,
		})
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
