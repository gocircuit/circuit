package main

import (
	"bytes"
	"text/template"
)

type A map[string]interface{} // template arguments

func Render(source string, aux interface{}) string {
	var w bytes.Buffer
	if err := template.Must(template.New("").Parse(source)).Execute(&w, aux); err != nil {
		panic(2)
	}
	return w.String()
}

func RenderFigurePngSvg(caption, file, width string) string {
	const source = `
	<p><center>
	<figure class="shadowless">
		{{.Body}}
		<div><em>{{.Caption}}</em></div>
	</figure>
	</center></p>
	`
	return Render(source, A{"Caption": caption, "Body": RenderPngSvg(file, width)})
}

func RenderPngSvg(file, width string) string {
	const source = `
		<object data="img/{{.Name}}.svg" type="image/svg+xml" width="{{.Width}}">
		<img src="img/{{.Name}}.png" alt="" />
		</object>`
	return Render(source, A{"Name": file, "Width": width})
}
