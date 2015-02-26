package mysql_nodejs

import (
	. "github.com/gocircuit/circuit/gocircuit.org/render"
)

func RenderMysqlNodejsBoot() string {
	return RenderHtml("Starting a MySQL and node.js stack using a circuit app", Render(bootBody, nil))
}

const bootBody = `


        `
