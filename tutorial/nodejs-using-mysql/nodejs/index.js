var mysql = require('mysql');

var connection = mysql.createConnection({
	host: 'localhost', // argument?
	user: 'charlie',
	password : 'brown',
	database : 'charlie-brown'
});

connection.connect();

connection.query('SELECT * from Friends',
	function(err, rows, fields) {
		if (!err)
			console.log('The solution is: ', rows);
		else
			console.log('Error while performing Query.');
	});

connection.end();
