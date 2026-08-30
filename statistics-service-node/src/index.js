require('dotenv').config();

const config = require('../config');
const { createApp } = require('./app');

const app = createApp();

app.listen(config.port, () => {
  console.log(`Statistics service listening on port ${config.port}`);
});
