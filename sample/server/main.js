const http = require('http');

const server = http.createServer((req, res) => {
  console.log('New request for', req.url)

  if (req.url === '/hello') {
    return res.end(JSON.stringify({
      message: 'Hello World!',
    }));
  }

  if (req.url === '/healthcheck') {
    res.statusCode = 200;

    return res.end(JSON.stringify({
      status: 'OK',
    }));
  }

  res.statusCode = 404;
  res.end('Not Found');
});

server.listen(3000, () => {
  console.log('Server is running on port 3000');
})
