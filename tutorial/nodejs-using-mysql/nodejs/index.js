// process command-line arguments
var program = require('commander');

program
	.version('0.0.1')
	.usage('--mysql_host <host> --mysql_port <port> --api_host <host> --api_port <port>')
	.option('--mysql_host <host>', 'MySQL server host')
	.option('--mysql_port <port>', 'MySQL server port', parseInt)
	.option('--api_host <host>', 'API host')
	.option('--api_port <port>', 'API port', parseInt)
	.parse(process.argv);

// web and mysql frameworks
var express = require("express");
var mysql = require("mysql");

var app = express();

var pool = mysql.createPool({
	connectionLimit: 100, //important
	host: program.mysql_host,
	port: program.mysql_port,
	user: "tutorial",
	password: "",
	database: 'tutorial',
	debug:  false
});


// DELETE request

// GET request
app.get("/pop", function(req, resp) {

	// SELECT must be "pop"
	pool.query('SELECT * from Messages', function(err, rows, fields) {
		if (!err) {
			console.log('/pop ok');
			resp.json(rows);
		} else {
			console.log('/pop error');
			resp.json({
				"code": 100,
				"status": "Database problem."
			});
		}
	});

});

app.listen(program.api_port, program.api_host);
console.log("Listening.");
