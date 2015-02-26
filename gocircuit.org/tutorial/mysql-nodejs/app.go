package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsApp() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(appBody, nil))
}

const appBody = `
<h1>Starting a MySQL and node.js stack using a circuit app</h1>

        `
