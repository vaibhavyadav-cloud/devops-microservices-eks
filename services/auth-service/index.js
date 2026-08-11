const express = require('express');
const env = require('./src/config/env');
const logger = require('./src/utils/logger');
const { initDb } = require('./db');
const authRoutes = require('./src/routes/auth.routes');
const { errorHandler } = require('./src/middleware/errorHandler');

const app = express();
app.use(express.json());

// Basic request logging
app.use((req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    logger.info('request', {
      method: req.method,
      path: req.path,
      status: res.statusCode,
      durationMs: Date.now() - start,
    });
  });
  next();
});

// --- Health checks (K8s liveness vs readiness — deliberately different) ---
// Liveness: "is the process alive" — never checks DB, so a slow DB doesn't
// cause K8s to kill+restart a perfectly healthy pod.
app.get('/health/live', (req, res) => {
  res.status(200).json({ status: 'ok' });
});

// Readiness: "can this pod actually serve traffic" — checks DB connectivity.
// If this fails, K8s pulls the pod out of the Service's endpoint list
// (no restart, just stops routing traffic to it) until it passes again.
app.get('/health/ready', async (req, res) => {
  try {
    const { pool } = require('./db');
    await pool.query('SELECT 1');
    res.status(200).json({ status: 'ready' });
  } catch (err) {
    res.status(503).json({ status: 'not_ready', error: err.message });
  }
});

// Keep old /health for backwards compat / simple checks
app.get('/health', (req, res) => res.status(200).json({ status: 'ok', service: 'auth-service' }));

app.use('/', authRoutes);

// Must be registered after all routes
app.use(errorHandler);

let server;

initDb()
  .then(() => {
    server = app.listen(env.port, () => {
      logger.info('service_started', { port: env.port, env: env.nodeEnv });
    });
  })
  .catch((err) => {
    logger.error('db_init_failed_starting_anyway', { message: err.message });
    server = app.listen(env.port, () => {
      logger.info('service_started', { port: env.port, env: env.nodeEnv });
    });
  });

// --- Graceful shutdown ---
// K8s sends SIGTERM before killing a pod (during rollout, scale-down, or node drain).
// Without this, in-flight requests get dropped mid-response.
function shutdown(signal) {
  logger.info('shutdown_signal_received', { signal });
  if (!server) process.exit(0);

  server.close(() => {
    logger.info('server_closed');
    process.exit(0);
  });

  // Force-exit if graceful close hangs (e.g. a stuck connection)
  setTimeout(() => {
    logger.error('forced_shutdown_timeout');
    process.exit(1);
  }, 10000).unref();
}

process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));
