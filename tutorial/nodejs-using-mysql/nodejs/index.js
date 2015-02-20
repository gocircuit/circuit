var mysql = require('mysql');

var connection = mysql.createConnection({
	host: "localhost:3306", // argument?
	user: "tutorial",
	password : "",
	database : "tutorial"
});

connection.connect();

connection.query('SELECT * from Events',
	function(err, rows, fields) {
		if (!err)
			console.log('The solution is: ', rows);
		else
			console.log('Error while performing query.');
	});

connection.end();
