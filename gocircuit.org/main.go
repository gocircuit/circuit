package main

import (
	"os"
	"path"
)

func main() {
	Build("cmd.html", RenderCommandPage())
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

// 	x:1 = *BuildIndexPage
// 	x:2 = *BuildInstallPage
// 	x:3 = *BuildRunPage
// 	x:4 = *BuildMetaphorPage
// 	x:5 = *BuildCommandPage
// 	x:6 = *BuildSecurityPage
// 	x:7 = *BuildHistoryPage

// 	x:ep = *BuildElementProcessPage
// 	x:ec = *BuildElementContainerPage
// 	x:es = *BuildElementSubscriptionPage
// 	x:ed = *BuildElementDnsPage
// 	x:eh = *BuildElementChannelPage

// 	x:tut = *
// }
