var express = require("express");
var mysql = require("mysql");

var app = express();

var pool = mysql.createPool({
	connectionLimit: 100, //important
	host: "localhost",
	port: 3306,
	user: "tutorial",
	password: "",
	database: 'tutorial',
	debug:  false
});

// on HTTP request
app.get("/pop", function(req, resp) {

	// SELECT must be "pop"
	pool.query('SELECT * from Messages', function(err, rows, fields) {
		if (!err) {
			console.log('Popped: ', rows);
			resp.json(rows);
		} else {
			console.log('Error while performing query.');
			resp.json({
				"code": 100,
				"status": "Database problem."
			});
		}
	});

});

app.listen(3000); // argument?
console.log("Listening.");
