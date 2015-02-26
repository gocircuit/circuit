package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsOverview() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(overviewBody, nil))
}

const overviewBody = `
<h1>Starting a MySQL and node.js stack using a circuit app</h1>

        `
