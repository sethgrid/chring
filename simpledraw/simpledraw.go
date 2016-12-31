package simpledraw

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

type Draw struct {
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

var Pallate = []color.RGBA{Red, Orange, Yellow, Green, Blue, Purple}

func Rad2Deg(rad float64) float64 {
	return rad / (2 * math.Pi) * 360
}

func (d *Draw) DrawCircle(c Circle) {
	d.ArcTo(c.X, c.Y, c.Radius, c.Radius, 0, 2*math.Pi)
	d.SetFillColor(c.Props.Color)
	d.SetStrokeColor(c.Props.Stroke)
	d.SetLineWidth(c.Props.Weight)
	d.FillStroke()
}

func (d *Draw) DrawSquare(s Square) {
	x, y := s.X-(s.Width/2), s.Y-(s.Width/2)
	d.MoveTo(x, y)
	d.LineTo(x+s.Width, y)
	d.LineTo(x+s.Width, y+s.Width)
	d.LineTo(x, y+s.Width)
	d.LineTo(x, y)
	d.SetFillColor(s.Props.Color)
	d.SetStrokeColor(s.Props.Stroke)
	d.SetLineWidth(s.Props.Weight)
	d.FillStroke()
}

// DrawTriangle TODO: this could be simplified into DrawPolygon(angle...)
func (d *Draw) DrawTriangle(t Triangle) {
	inscribedCircle := NewCircle(t.X, t.Y, t.Radius)
	x1, y1 := inscribedCircle.PointAtAngle(-math.Pi / 2)
	x2, y2 := inscribedCircle.PointAtAngle(-math.Pi * 7 / 6)
	x3, y3 := inscribedCircle.PointAtAngle(-math.Pi * 11 / 6)
	d.MoveTo(x1, y1)
	d.LineTo(x2, y2)
	d.LineTo(x3, y3)
	d.LineTo(x1, y1)
	d.SetFillColor(t.Props.Color)
	d.SetStrokeColor(t.Props.Stroke)
	d.SetLineWidth(t.Props.Weight)
	d.FillStroke()
}

func (d *Draw) WriteStringAt(text string, x, y float64) (width float64) {
	d.SetFillColor(Black)
	return d.FillStringAt(text, x, y)
}

func (d *Draw) WriteBoldStringAt(text string, x, y float64) (width float64) {
	f := d.GetFontData()
	f.Style = draw2d.FontStyleBold
	d.SetFontData(f)

	defer func() {
		f.Style = draw2d.FontStyleNormal
		d.SetFontData(f)
	}()
	return d.WriteStringAt(text, x, y)
}

func (d *Draw) DrawLegend(l *Legend) {
	var x, y float64 = 25, 25

	d.WriteBoldStringAt(l.Title, x, y)
	_, top, _, bottom := d.GetStringBounds(l.Title)
	y += (bottom - top) + 5

	d.WriteStringAt(l.Caption, x, y)
	_, top, _, bottom = d.GetStringBounds(l.Caption)
	y += (bottom - top) + 15

	for _, e := range l.Elements {
		if e.IsCircle {
			shape := NewCircle(x, y-4, 4)
			shape.Props = e.Props
			d.DrawCircle(shape)
		}
		if e.IsSquare {
			shape := NewSquare(x, y-4, 8)
			shape.Props = e.Props
			d.DrawSquare(shape)
		}
		if e.IsTriangle {
			shape := NewTriangle(x, y-4, 6)
			shape.Props = e.Props
			d.DrawTriangle(shape)
		}

		d.WriteStringAt(e.Name, x+10, y)
		_, top, _, bottom = d.GetStringBounds(l.Caption)
		y += (bottom - top) + 5
	}
}

type BasicProperties struct {
	Color, Stroke color.RGBA
	Weight        float64
}

type Triangle struct {
	Props        BasicProperties
	X, Y, Radius float64
}

type Square struct {
	Props       BasicProperties
	X, Y, Width float64
}

type Circle struct {
	Props        BasicProperties
	X, Y, Radius float64
}

func NewCircle(x, y, radius float64) Circle {
	return Circle{
		X:      x,
		Y:      y,
		Radius: radius,
		Props: BasicProperties{
			Color:  White,
			Stroke: Black,
			Weight: 1,
		},
	}
}

func NewSquare(x, y, width float64) Square {
	return Square{
		X:     x,
		Y:     y,
		Width: width,
		Props: BasicProperties{
			Color:  White,
			Stroke: Black,
			Weight: 1,
		},
	}
}

func NewTriangle(x, y, radius float64) Triangle {
	return Triangle{
		X:      x,
		Y:      y,
		Radius: radius,
		Props: BasicProperties{
			Color:  White,
			Stroke: Black,
			Weight: 1,
		},
	}
}

func (c Circle) PointAtAngle(radian float64) (float64, float64) {
	return math.Cos(radian)*c.Radius + c.X, math.Sin(radian)*c.Radius + c.Y
}

type Legend struct {
	Title    string
	Caption  string
	Elements []LegendElement
}

type LegendElement struct {
	Name          string
	IsTriangle    bool
	IsSquare      bool
	IsCircle      bool
	Width, Height float64
	Props         BasicProperties
}

func (l *Legend) ContentWidth() float64 {
	var greatest float64
	for _, e := range l.Elements {
		if e.Width > greatest {
			greatest = e.Width
		}
	}
	return greatest + 5
}

func (l *Legend) ContentHeight() float64 {
	var total float64
	var top, bottom float64
	dest := image.NewRGBA(image.Rect(0, 0, 500, 500))
	gc := draw2dimg.NewGraphicContext(dest)
	_, top, _, bottom = gc.GetStringBounds(l.Title)
	total += (bottom - top)
	_, top, _, bottom = gc.GetStringBounds(l.Caption)
	total += (bottom - top)

	for _, e := range l.Elements {
		total += e.Height
	}
	return total + (5 * float64(len(l.Elements)))
}
