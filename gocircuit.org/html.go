package main

func RenderHtml(title, body string) string {
	return Render(sourceHtml, A{"Header": sourceHeader, "Footer": sourceFooter})
}

const sourceHtml = `<!doctype html><html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<link href="css/main.css" rel="stylesheet" type="text/css" />
		<title>{{.Title}}</title>
	</head>
	<body>
	{{.Header}}
	<div class="page">
	{{.Body}}
	</div>
	{{.Footer}}
	</body>
	</html>`

const sourceFooter = `
	<div class="footer">
	The <a href="http://escher.io">Escher</a> and <a href="http://gocircuit.org">Circuit</a> projects are
	partially supported by the 
	<a href="http://www.darpa.mil/Our_Work/I2O/Programs/XDATA.aspx">DARPA XData Initiative</a>.<br>
	Sponsors and partners are welcome and appreciated. Contact <a href="mailto:p@gocircuit.org">Petar Maymounkov</a> for details.
	</div>`

const sourceHeader = `
	<div class="header">
	<a href="http://gocircuit.org">Circuit</a> Self-managed infrastructure, programmatic monitoring and orchestration
	</div>`
