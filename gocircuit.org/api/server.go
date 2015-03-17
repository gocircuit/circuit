package api

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderServerPage() string {
	return RenderHtml("Using server", Render(serverBody, nil))
}

const serverBody = `

<h2>Using servers</h2>

<p>??

        `
