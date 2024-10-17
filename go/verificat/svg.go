package main

import (
	"fmt"
	"log/slog"
)

// TODO: svgStart needs to be broken up more to programmatically set the viewBox height
// i.e., the viewBox should expand as the count of elements expands
const (
	// Changing the /viewBox/ here can have drastic scale effects
	// A width of 200 in the viewBox when the width of the SVG is 400
	// will effectively double the relative scale of all shapes in the viewBox
	svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
    width="400px"
    height="auto"
	viewBox="0 0 200 800"
    version="2.0">`
	svgEnd = `</svg>`
)

// Build the SVG XML for rendering the Almanac data
func BuildSVG(a *Almanac, sc *SVGCfg) string {
	var c int // count of Services processed

	xmlsvg := svgStart
	for _, service := range *a {
		// for each iteration, pass an increasing number to hashBar
		// so that y values are moved down the page
		c += 1
		xmlsvg = xmlsvg + hashBarSVG(c, &service, sc)
		slog.Debug("SVG Added", slog.Int("SVG Count", c), slog.Any("Service", service.Name))
	}
	xmlsvg = xmlsvg + svgEnd
	return xmlsvg
}

// SVGCfg holds starting points for drawing the graphic,
// described (and named) for what they create.
type SVGCfg struct {
	Gutter int // Vertical space above the graphic
	TxtOff int // Vertical offset to the baseline of the text
	Spacer int // Vertical space between each graphic
}

// hashBarSVG displays a colorful representation of
// failed tests versus the total score of 100.
//
// The 'background bar' is revealed by the 'foreground bar' and
// as more tests come up "fail", more of the 'foreground bar' will recede
// and display more of the 'background bar'.
func hashBarSVG(c int, sd *WMService, sc *SVGCfg) string {
	// XML Path is used to draw boxes and fill them.
	// background bar is static at 100, but it's Y starting point changes
	// foreground bar is fully dynamic
	// the score and name are both dynamic
	//	NB: The color /fill/ for the first line of text
	//	remains in hex to allow future score-dependent gradient coloring
	//		score > 90 ::: fill="#f08080"
	//		90 > score > 80 ::: fill="#f00000"
	//		etc.
	//
	// These variable names match the notation for the XML Path
	// v == X coordinate for drawing the box (top) == sd.Score
	// H == Y coordinate for drawing the box (right side)
	// z == X coordinate for completing the box (bottom)
	//
	h := sc.Gutter + (c * sc.Spacer)             // Y coordinate for drawing the box (top)
	y := sc.Gutter + sc.TxtOff + (c * sc.Spacer) //  Y coordinate for drawing the text line
	return fmt.Sprintf(`<path d="M5 %dh100v5H5z" fill="tomato" />
    <path d="M5 %dh%dv10H5z" fill="darkgreen" />
    <text x="10" y="%d" font-family="Helvetica" font-size="14" font-weight="bolder" font-style="oblique" fill="#f08080" fill-opacity=".5" stroke-width=".5" stro-kopacity=".5" stroke-linejoin="round" stroke="coral">%d</text>
	<text x="30" y="%d" font-family="Palatino" font-style="oblique" font-size="6" fill="snow" fill-opacity="1" >%d</text>
    <text x="40" y="%d" font-family="Monaco" font-size="12" fill="turquoise" fill-opacity="1" stroke-width=".5" stroke-linejoin="miter" stroke="darkmagenta">%s</text>`,
		h, h, sd.Score,
		y, sd.Score,
		y, sd.LastID,
		y, sd.Name)
}
