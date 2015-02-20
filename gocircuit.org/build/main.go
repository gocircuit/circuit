package main

import (
	"os"
	"path"

	. "github.com/gocircuit/circuit/gocircuit.org/man"
)

func main() {
	Build("index.html", RenderIndexPage())
	Build("install.html", RenderInstallPage())
	Build("cmd.html", RenderCommandPage())
	Build("history.html", RenderHistoryPage())
	Build("security.html", RenderSecurityPage())
	Build("metaphor.html", RenderMetaphorPage())
	Build("run.html", RenderRunPage())

	Build("element-process.html", RenderElementProcessPage())
	Build("element-container.html", RenderElementContainerPage())
	Build("element-subscription.html", RenderElementSubscriptionPage())
	Build("element-dns.html", RenderElementDnsPage())
	Build("element-channel.html", RenderElementChannelPage())

	// 	x:tut = *
}

func Build(file string, content string) {
	dir, file := path.Split(file)
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
