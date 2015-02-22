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

// PUT request
app.get("/put/:name/:value", function(req, resp) {
	var name = req.params.name;
	var value = req.params.value;
	pool.query('REPLACE INTO NameValue SET name = ' + mysql.escape(name) + ', value = ' + mysql.escape(value),
		function(err, rows, fields) {
			if (!err) {
				console.log('put ok');
				resp.json({status: "ok"});
			} else {
				console.log('put error');
				resp.json({status: "error"});
			}
		}
	);
});

// DELETE request
app.get("/del/:name", function(req, resp) {
	var name = req.params.name;
	pool.query('DELETE FROM NameValue WHERE name = ' + mysql.escape(name), function(err, rows, fields) {
		if (!err) {
			console.log('del ok');
			resp.json({status: "ok"});
		} else {
			console.log('del error');
			resp.json({status: "error"});
		}
	});
});

// GET request
app.get("/get/:name", function(req, resp) {
	var name = req.params.name;
	pool.query('SELECT * from NameValue WHERE name = ' + mysql.escape(name),
		function(err, rows, fields) {
			if (!err) {
				console.log('get ok');
				resp.json({status: "ok", result: rows});
			} else {
				console.log('get error');
				resp.json({status: "error"});
			}
		}
	);
});

app.listen(program.api_port, program.api_host);
console.log("Listening.");
