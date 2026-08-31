import assert from 'node:assert/strict';
import test from 'node:test';

import { createApp } from './app';

interface ServerHarness {
  baseUrl: string;
  close: () => Promise<void>;
}

function createServerHarness(): Promise<ServerHarness> {
  const app = createApp();
  const server = app.listen(0, '127.0.0.1');

  return new Promise((resolve, reject) => {
    server.once('error', reject);
    server.once('listening', () => {
      const address = server.address();
      if (!address || typeof address === 'string') {
        reject(new Error('Expected numeric port from ephemeral test server'));
        return;
      }

      const baseUrl = `http://127.0.0.1:${address.port}`;

      async function close() {
        await new Promise<void>((closeResolve, closeReject) => {
          server.close((error) => {
            if (error) {
              closeReject(error);
              return;
            }
            closeResolve();
          });
        });
      }

      resolve({ baseUrl, close });
    });
  });
}

test('GET /health returns current health contract', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/health`);
    const body = await response.json();

    assert.equal(response.status, 200);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      status: 'ok',
      service: 'statistics-service-node',
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics returns statistics for valid matrices', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        q: [
          [1, 2],
          [3, 4],
        ],
        r: [
          [5, 6],
          [7, 8],
        ],
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 200);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      max: 8,
      min: 1,
      average: 4.5,
      sum: 36,
      hasDiagonalMatrix: false,
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics returns INVALID_REQUEST for malformed JSON', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: '{"q":[[1]],"r":',
    });

    const body = await response.json();

    assert.equal(response.status, 400);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      error: {
        code: 'INVALID_REQUEST',
        message: 'Request body must be valid JSON',
      },
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics returns INVALID_MATRICES for missing q', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        r: [[1]],
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 400);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      error: {
        code: 'INVALID_MATRICES',
        message: 'q is required',
      },
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics returns INVALID_MATRICES for missing r', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        q: [[1]],
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 400);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      error: {
        code: 'INVALID_MATRICES',
        message: 'r is required',
      },
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics returns INVALID_MATRICES for invalid q', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        q: [],
        r: [[1]],
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 400);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      error: {
        code: 'INVALID_MATRICES',
        message: 'q must be a valid matrix',
      },
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics returns INVALID_MATRICES for invalid r', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        q: [[1]],
        r: [],
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 400);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      error: {
        code: 'INVALID_MATRICES',
        message: 'r must be a valid matrix',
      },
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics accepts valid rectangular matrices', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        q: [[1, 2]],
        r: [[3, 4]],
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 200);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      max: 4,
      min: 1,
      average: 2.5,
      sum: 10,
      hasDiagonalMatrix: false,
    });
  } finally {
    await harness.close();
  }
});

test('POST /api/v1/statistics ignores additional request properties', async () => {
  const harness = await createServerHarness();

  try {
    const response = await fetch(`${harness.baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        q: [
          [1, 2],
          [3, 4],
        ],
        r: [
          [5, 6],
          [7, 8],
        ],
        extraField: 'ignored',
      }),
    });

    const body = await response.json();

    assert.equal(response.status, 200);
    assert.ok(response.headers.get('content-type')?.includes('application/json'));
    assert.deepEqual(body, {
      max: 8,
      min: 1,
      average: 4.5,
      sum: 36,
      hasDiagonalMatrix: false,
    });
  } finally {
    await harness.close();
  }
});
