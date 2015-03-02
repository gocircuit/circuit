package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderRun() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(runBody, nil))
}

const runBody = `
<h1>Starting a MySQL and node.js stack using a circuit app</h1>

        `
