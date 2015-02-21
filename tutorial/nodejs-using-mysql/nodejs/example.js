// var http = require('http');
// http.createServer(function (req, res) {
// 	res.writeHead(200, {'Content-Type': 'text/plain'});
// 	res.end('Hello World\n');
// }).listen(1337, '127.0.0.1');
// console.log('Server running at http://127.0.0.1:1337/');

var program = require('commander');

program
	.version('0.0.1')
	.usage('--mysql_host <host> --mysql_port <port> --api <addr>')
	.option('--mysql_host <host>', 'MySQL server host')
	.option('--mysql_port <port>', 'MySQL server port', parseInt)
	.option('-a, --api <addr>', 'API listen address')
	.parse(process.argv);

console.log("mh", program.mysql_host)
console.log("mp", program.mysql_port)
console.log("a", program.api)
