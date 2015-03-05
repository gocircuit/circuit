package render

func RenderHtml(title, body string) string {
	return Render(sourceHtml, A{
		"Title":      title,
		"Body":       body,
		"PathToRoot": "",
		"Header":     sourceHeader,
		"Footer":     sourceFooter,
	})
}

func RenderHtml2(url []string, title, body string) string {
	return Render(sourceHtml, A{
		"Title":      title,
		"Body":       body,
		"PathToRoot": PathToRoot(url),
		"Header":     sourceHeader,
		"Footer":     sourceFooter,
	})
}

// PathToRoot returns "", "../", "../../", etc.
func PathToRoot(url []string) string {
	r := ""
	for range url {
		r += "../"
	}
	return r
}

const sourceHtml = `<!doctype html><html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<link href="{{.PathToRoot}}css/main.css" rel="stylesheet" type="text/css" />
		<title>{{.Title}}</title>
		<link href='http://fonts.googleapis.com/css?family=Lato&subset=latin,latin-ext' rel='stylesheet' type='text/css'>
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
