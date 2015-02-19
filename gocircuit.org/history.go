package main

func RenderHistoryPage() string {
	return RenderHtml("History, links and bibliography", Render(historyBody, nil))
}

const historyBody = `

<h3>Sponsors and credits</h3>

<ul>
<li><a href="http://www.darpa.mil/Our_Work/I2O/Programs/XDATA.aspx">DARPA XDATA</a> initiative, 2012–2014
<li><a href="http://www.data-tactics.com/">Data Tactics Corp</a>, 2012-2014
<li><a href="http://www.l-3com.com/">L3</a>, 2014
<li><a href="http://tumblr.com">Tumblr, Inc.</a>, 2012
</ul>

<h3>Presentations</h3>

<ul>
<li><a href="http://www.darpa.mil">DARPA</a> <a href="http://www.darpa.mil/opencatalog/">Open Catalog</a>, Arlington, VA, 2014
<li><a href="http://confreaks.com/videos/3421-gophercon2014-the-go-circuit-towards-elastic-computation-with-no-failures">GOPHERCON 2014</a>, Denver, CO
<li><a href="http://blog.gocircuit.org/strangeloop-2013">STRANGELOOP 2013</a>, St. Louis, MO
</ul>

<h3>Links</h3>

<ul>
<li><a href="http://gocircuit-org.appspot.com">Old Circuit website, original white papers and design documents</a>
</ul>

        `
