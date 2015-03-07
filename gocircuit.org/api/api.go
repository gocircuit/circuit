package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMainPage() string {
	figs := A{
		"FigHierarchy": RenderFigurePngSvg(
			"Virtual anchor hierarchy and its mapping to Go <code>Anchor</code> objects.", "hierarchy", "600px"),
	}
	return RenderHtml("Go client API", Render(mainBody, figs))
}

const mainBody = `

<h2>Go client API</h2>

{{.FigHierarchy}}

        `
