package main

import (
	"os"
	"path"

	"github.com/gocircuit/circuit/gocircuit.org/api"
	"github.com/gocircuit/circuit/gocircuit.org/man"
	"github.com/gocircuit/circuit/gocircuit.org/tutorial"
)

func main() {
	Build("index.html", RenderIndexPage())
	Build("install.html", man.RenderInstallPage())
	Build("cmd.html", man.RenderCommandPage())
	Build("history.html", man.RenderHistoryPage())
	Build("security.html", man.RenderSecurityPage())
	Build("metaphor.html", man.RenderMetaphorPage())
	Build("run.html", man.RenderRunPage())

	Build("element-process.html", man.RenderElementProcessPage())
	Build("element-container.html", man.RenderElementContainerPage())
	Build("element-subscription.html", man.RenderElementSubscriptionPage())
	Build("element-dns.html", man.RenderElementDnsPage())
	Build("element-channel.html", man.RenderElementChannelPage())

	Build("api-process.html", api.RenderProcessPage())

	Build("tutorial-mysql-nodejs.html", tutorial.RenderMysqlNodejsPage())
}

func Build(file string, content string) {
	dir, _ := path.Split(file)
	if len(dir) > 0 {
		if err := os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}
	}
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte(content))
}
