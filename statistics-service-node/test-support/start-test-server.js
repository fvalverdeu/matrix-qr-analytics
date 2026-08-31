const { createApp } = require('../dist/src/app');

const HOST = '127.0.0.1';

const app = createApp();
const server = app.listen(0, HOST);

let shuttingDown = false;

function shutdown(signal) {
  if (shuttingDown) {
    return;
  }
  shuttingDown = true;

  server.close((error) => {
    if (error) {
      console.error(`ERROR: graceful shutdown failed after ${signal}: ${error.message}`);
      process.exit(1);
      return;
    }

    process.exit(0);
  });
}

server.once('error', (error) => {
  console.error(`ERROR: failed to start test server: ${error.message}`);
  process.exit(1);
});

server.once('listening', () => {
  const address = server.address();

  if (!address || typeof address === 'string' || typeof address.port !== 'number') {
    console.error('ERROR: test server did not expose a numeric TCP port');
    server.close(() => process.exit(1));
    return;
  }

  process.stdout.write(`TEST_SERVER_PORT=${address.port}\n`);
});

process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));
