package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderAnchorPage() string {
	figs := A{
		"FigHierarchy": RenderFigurePngSvg(
			"Virtual anchor hierarchy and its mapping to Go <code>Anchor</code> objects.", "hierarchy", "600px"),
	}
	return RenderHtml("Navigating the anchor hierarchy", Render(anchorBody, figs))
}

const anchorBody = `

<h2>Navigating the anchor hierarchy</h2>

{{.FigHierarchy}}


        `
