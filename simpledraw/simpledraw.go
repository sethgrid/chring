package simpledraw

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Draw embeds draw2dimg.GraphicContext and gives us some handy simple drawing tools
type Draw struct {
	*draw2dimg.GraphicContext
}

// Set some decent default colors

// Red is a sensible default RGB color choice
var Red = color.RGBA{244, 86, 66, 255}

// Orange is a sensible default RGB color choice
var Orange = color.RGBA{234, 133, 44, 255}

// Yellow is a sensible default RGB color choice
var Yellow = color.RGBA{239, 235, 117, 255}

// Green is a sensible default RGB color choice
var Green = color.RGBA{62, 214, 89, 255}

// Blue is a sensible default RGB color choice
var Blue = color.RGBA{51, 147, 204, 255}

// Purple is a sensible default RGB color choice
var Purple = color.RGBA{175, 102, 209, 255}

// White is a sensible default RGB color choice
var White = color.RGBA{255, 255, 255, 255}

// Black is a sensible default RGB color choice
var Black = color.RGBA{0, 0, 0, 255}

// Pallate is a collection of preset colors for easy cycling through in list iterations
var Pallate = []color.RGBA{Red, Orange, Yellow, Green, Blue, Purple}

// Rad2Deg is a simple helper for converting from Radians to Degrees, usually used for printing
func Rad2Deg(rad float64) float64 {
	return rad / (2 * math.Pi) * 360
}

// DrawCircle draws the given circle
func (d *Draw) DrawCircle(c Circle) {
	d.ArcTo(c.X, c.Y, c.Radius, c.Radius, 0, 2*math.Pi)
	d.SetFillColor(c.Props.Color)
	d.SetStrokeColor(c.Props.Stroke)
	d.SetLineWidth(c.Props.Weight)
	d.FillStroke()
}

// DrawRegularPolygon draws a regular polygon, but handles special case of also drawing a circle if sides <= 1.
// Could add wrappers for DrawSquare or DrawTriangle if warrented. Not needed currently.
func (d *Draw) DrawRegularPolygon(sides int, x, y, radius float64, props BasicProperties) {
	// special case for making program flow more simple
	if sides <= 1 {
		c := NewCircle(x, y, radius)
		c.Props = props
		d.DrawCircle(c)
		return
	}

	// adj is how much arc we adjust the angle for each iteration
	adj := 2 * math.Pi / float64(sides)

	// angle is our starting angle and will be adjusted by adj for each side
	angle := math.Pi / 2 * float64(sides)

	// keep odd sided polygons with a vertex at the top (90 deg)
	if sides%2 == 1 {
		angle = math.Pi / 2
	}
	// special case, want a square over a diamond
	if sides == 4 {
		angle = math.Pi / 4
	}

	// a regular polygon is easiest to draw if inscribed in a circle
	inscribedCircle := NewCircle(x, y, radius)

	for i := 0; i <= sides; i++ {
		xN, yN := inscribedCircle.PointAtAngle(-angle)
		angle -= adj
		if i == 0 {
			d.MoveTo(xN, yN)
		}
		d.LineTo(xN, yN)
	}

	d.SetFillColor(props.Color)
	d.SetStrokeColor(props.Stroke)
	d.SetLineWidth(props.Weight)
	d.FillStroke()
}

// WriteStringAt writes black text at the given location
func (d *Draw) WriteStringAt(text string, x, y float64) (width float64) {
	d.SetFillColor(Black)
	return d.FillStringAt(text, x, y)
}

// WriteBoldStringAt writes bold black text at the given location
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

// DrawLegend draws out the legend and all its registered elements
func (d *Draw) DrawLegend(l *Legend) {
	var x, y float64 = 25, 25

	d.WriteBoldStringAt(l.Title, x, y)
	_, top, _, bottom := d.GetStringBounds(l.Title)
	y += (bottom - top) + 5

	d.WriteStringAt(l.Caption, x, y)
	_, top, _, bottom = d.GetStringBounds(l.Caption)
	y += (bottom - top) + 15

	for _, e := range l.Elements {
		var radius float64 = 5
		if e.PolygonSides <= 1 { // make circle a bit smaller
			radius = 4
		}
		d.DrawRegularPolygon(e.PolygonSides, x, y-4, radius, e.Props)
		d.WriteStringAt(e.Name, x+10, y)
		_, top, _, bottom = d.GetStringBounds(l.Caption)
		y += (bottom - top) + 5
	}
}

// BasicProperties contains the basic properties that all shapes need
type BasicProperties struct {
	Color, Stroke color.RGBA
	Weight        float64
}

// DefaultBasicProperties provides sensible defaults for shape properties
var DefaultBasicProperties = BasicProperties{
	Color:  White,
	Stroke: Black,
	Weight: 1,
}

// Circle is a simple struct with circle's properties
type Circle struct {
	Props        BasicProperties
	X, Y, Radius float64
}

// NewCircle returns a circle with sensible defaults
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

// PointAtAngle returns the x, y pair along a circle's circumference at the given angle in radians
func (c Circle) PointAtAngle(radian float64) (float64, float64) {
	return math.Cos(radian)*c.Radius + c.X, math.Sin(radian)*c.Radius + c.Y
}

// DrawOnEdge will draw a regular poloygon on the circumference of a given circle
func (d *Draw) DrawOnEdge(c Circle, angle float64, sides int, radius float64, props BasicProperties) {
	x, y := c.PointAtAngle(angle)
	d.DrawRegularPolygon(sides, x, y, radius, props)
}

// Legend is used for displaying a legend, as for a chart
type Legend struct {
	Title    string
	Caption  string
	Elements []LegendElement
}

// LegendElement is used in the Legend struct and is not likely to be used outside the package. It is exported just in case it it needed.
type LegendElement struct {
	Name          string
	PolygonSides  int
	Width, Height float64
	Props         BasicProperties
}

// PrependElement is a helper to cut down on visual clutter when developing. Not async.
func (l *Legend) PrependElement(polygonSides int, name string, props BasicProperties) {
	dest := image.NewRGBA(image.Rect(0, 0, 0, 0)) // some dest is needed
	gc := draw2dimg.NewGraphicContext(dest)
	left, top, right, bottom := gc.GetStringBounds(name)
	l.Elements = append([]LegendElement{{
		PolygonSides: polygonSides,
		Name:         name,
		Props:        props,
		Width:        right - left,
		Height:       bottom - top,
	}}, l.Elements...)
}

// AppendElement is a helper to cut down on visual clutter when developing. Not async.
func (l *Legend) AppendElement(polygonSides int, name string, props BasicProperties) {
	dest := image.NewRGBA(image.Rect(0, 0, 0, 0)) // some dest is needed
	gc := draw2dimg.NewGraphicContext(dest)
	left, top, right, bottom := gc.GetStringBounds(name)
	l.Elements = append(l.Elements, LegendElement{
		PolygonSides: polygonSides,
		Name:         name,
		Props:        props,
		Width:        right - left,
		Height:       bottom - top,
	})
}

// ContentWidth returns the width of the computed legend
func (l *Legend) ContentWidth() float64 {
	var greatest float64
	for _, e := range l.Elements {
		if e.Width > greatest {
			greatest = e.Width
		}
	}
	return greatest + 5
}

// ContentHeight returns the height of the computed legend
func (l *Legend) ContentHeight() float64 {
	var total float64
	var top, bottom float64
	dest := image.NewRGBA(image.Rect(0, 0, 0, 0)) // some dest is needed
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
