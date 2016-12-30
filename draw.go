package chring

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/llgcode/draw2d/draw2dimg"
)

type simpleDraw struct {
	*draw2dimg.GraphicContext
}

var Red = color.RGBA{244, 86, 66, 255}
var Orange = color.RGBA{234, 133, 44, 255}
var Yellow = color.RGBA{239, 235, 117, 255}
var Green = color.RGBA{62, 214, 89, 255}
var Blue = color.RGBA{51, 147, 204, 255}
var Purple = color.RGBA{175, 102, 209, 255}
var White = color.RGBA{255, 255, 255, 255}
var Black = color.RGBA{0, 0, 0, 255}

var pallate = []color.RGBA{Red, Orange, Yellow, Green, Blue, Purple}

func (r *Ring) drawChart(w http.ResponseWriter, req *http.Request) {
	m, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	dest := image.NewRGBA(image.Rect(0, 0, 500, 500))
	gc := simpleDraw{draw2dimg.NewGraphicContext(dest)}

	var x, y float64
	ring := newCircle(250, 250, 150)
	gc.drawCircle(ring)

	for i, n := range r.Nodes {
		x, y = ring.pointAtAngle(hashAngle(n.HashID))
		c := newCircle(x, y, 10)
		c.props.color = pallate[i%len(pallate)]
		gc.drawCircle(c)
	}

	for i, param := range m["key[]"] {
		hashID := r.Hasher(param)
		x, y = ring.pointAtAngle(hashAngle(hashID))
		s := newSquare(x, y, 8)
		s.props.color = pallate[i%len(pallate)+3]
		gc.drawSquare(s)
	}

	for i, param := range m["hashid[]"] {
		hashID, _ := strconv.Atoi(param)
		x, y = ring.pointAtAngle(hashAngle(uint32(hashID)))
		t := newTriangle(x, y, 8)
		t.props.color = pallate[i%len(pallate)+5]
		gc.drawTriangle(t)
	}

	w.Header().Set("Content-Type", "image/png")

	err = png.Encode(w, dest) //Encode writes the Image m to w in PNG format.
	if err != nil {
		fmt.Printf("Error rendering pie chart: %v\n", err)
	}
}

func Rad2Deg(rad float64) float64 {
	return rad / (2 * math.Pi) * 360
}

func hashAngle(hashID uint32) float64 {
	return 2 * math.Pi * float64(hashID) / float64(math.MaxUint32)
}

func (d *simpleDraw) drawCircle(c circle) {
	d.ArcTo(c.x, c.y, c.radius, c.radius, 0, 2*math.Pi)
	d.SetFillColor(c.props.color)
	d.SetStrokeColor(c.props.stroke)
	d.SetLineWidth(c.props.weight)
	d.FillStroke()
}

func (d *simpleDraw) drawSquare(s square) {
	x, y := s.x-(s.width/2), s.y-(s.width/2)
	d.MoveTo(x, y)
	d.LineTo(x+s.width, y)
	d.LineTo(x+s.width, y+s.width)
	d.LineTo(x, y+s.width)
	d.LineTo(x, y)
	d.SetFillColor(s.props.color)
	d.SetStrokeColor(s.props.stroke)
	d.SetLineWidth(s.props.weight)
	d.FillStroke()
}

func (d *simpleDraw) drawTriangle(t triangle) {
	inscribedCircle := newCircle(t.x, t.y, t.radius)
	x1, y1 := inscribedCircle.pointAtAngle(-math.Pi / 2)
	x2, y2 := inscribedCircle.pointAtAngle(-math.Pi * 7 / 6)
	x3, y3 := inscribedCircle.pointAtAngle(-math.Pi * 11 / 6)
	d.MoveTo(x1, y1)
	d.LineTo(x2, y2)
	d.LineTo(x3, y3)
	d.LineTo(x1, y1)
	d.SetFillColor(t.props.color)
	d.SetStrokeColor(t.props.stroke)
	d.SetLineWidth(t.props.weight)
	d.FillStroke()
}

type basicProperties struct {
	color, stroke color.RGBA
	weight        float64
}

type triangle struct {
	props        basicProperties
	x, y, radius float64
}

type square struct {
	props       basicProperties
	x, y, width float64
}

type circle struct {
	props        basicProperties
	x, y, radius float64
}

func newCircle(x, y, radius float64) circle {
	return circle{
		x:      x,
		y:      y,
		radius: radius,
		props: basicProperties{
			color:  White,
			stroke: Black,
			weight: 1,
		},
	}
}

func newSquare(x, y, width float64) square {
	return square{
		x:     x,
		y:     y,
		width: width,
		props: basicProperties{
			color:  White,
			stroke: Black,
			weight: 1,
		},
	}
}

func newTriangle(x, y, radius float64) triangle {
	return triangle{
		x:      x,
		y:      y,
		radius: radius,
		props: basicProperties{
			color:  White,
			stroke: Black,
			weight: 1,
		},
	}
}

func (c circle) pointAtAngle(radian float64) (float64, float64) {
	return math.Cos(radian)*c.radius + c.x, math.Sin(radian)*c.radius + c.y
}

func Serve(r *Ring) {
	http.HandleFunc("/", r.drawChart)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
