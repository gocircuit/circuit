package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderRun() string {
	return RenderHtml("Run the app on the cluster", Render(runBody, nil))
}

const runBody = `
<h1>Run the app on the cluster</h1>



        `
