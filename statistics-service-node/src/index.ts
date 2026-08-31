import dotenv from 'dotenv';

import config from '../config';
import { createApp } from './app';

dotenv.config();

const app = createApp();

app.listen(config.port, () => {
  console.log(`Statistics service listening on port ${config.port}`);
});
