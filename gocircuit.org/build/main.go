package main

import (
	"os"
	"path"

	"github.com/hoijui/circuit/gocircuit.org/api"
	"github.com/hoijui/circuit/gocircuit.org/man"
	mysql_nodejs "github.com/hoijui/circuit/gocircuit.org/tutorial/mysql-nodejs"
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
	Build("element-server.html", man.RenderElementServerPage())
	Build("element-channel.html", man.RenderElementChannelPage())

	Build("api.html", api.RenderMainPage())
	Build("api-connect.html", api.RenderConnectPage())
	Build("api-anchor.html", api.RenderAnchorPage())
	Build("api-process.html", api.RenderProcessPage())
	Build("api-container.html", api.RenderContainerPage())
	Build("api-subscription.html", api.RenderSubscriptionPage())
	Build("api-name.html", api.RenderNamePage())
	Build("api-server.html", api.RenderServerPage())
	Build("api-channel.html", api.RenderChannelPage())

	Build("tutorial-mysql-nodejs-overview.html", mysql_nodejs.RenderOverview())
	Build("tutorial-mysql-nodejs-image.html", mysql_nodejs.RenderImage())
	Build("tutorial-mysql-nodejs-boot.html", mysql_nodejs.RenderBoot()) //
	Build("tutorial-mysql-nodejs-app.html", mysql_nodejs.RenderApp())
	Build("tutorial-mysql-nodejs-run.html", mysql_nodejs.RenderRun())
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
