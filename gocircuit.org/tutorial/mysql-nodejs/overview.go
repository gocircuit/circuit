package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsOverview() string {
	figs := A{
		"FigOverview": RenderFigurePngSvg(
			`A cloud service comprised of a Node.js HTTP RESTful public API for a key/value store, 
			backed by MySQL database.`,
			"tutorial/mysql-nodejs",
			"600px",
		),
	}
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(overviewBody, figs))
}

const overviewBody = `
<h1>Overview</h1>

{{.FigOverview}}

        `
