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

	// 	x:ep = *BuildElementProcessPage
	// 	x:ec = *BuildElementContainerPage
	// 	x:es = *BuildElementSubscriptionPage
	// 	x:ed = *BuildElementDnsPage
	// 	x:eh = *BuildElementChannelPage
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
