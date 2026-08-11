const logger = require('../utils/logger');

// Custom error class so controllers can throw with an intended HTTP status
// instead of every route hand-rolling try/catch + res.status(...).json(...).
class AppError extends Error {
  constructor(statusCode, code, message) {
    super(message);
    this.statusCode = statusCode;
    this.code = code;
  }
}

// Wraps async route handlers so a rejected promise reaches the error
// middleware below instead of crashing the process / hanging the request.
function asyncHandler(fn) {
  return (req, res, next) => Promise.resolve(fn(req, res, next)).catch(next);
}

// Must be registered LAST, after all routes — Express identifies error
// middleware by its 4-argument signature (err, req, res, next).
function errorHandler(err, req, res, next) { // eslint-disable-line no-unused-vars
  if (err.code === '23505') {
    // Postgres unique_violation
    return res.status(409).json({ error: 'conflict', message: 'resource already exists' });
  }

  if (err instanceof AppError) {
    return res.status(err.statusCode).json({ error: err.code, message: err.message });
  }

  logger.error('unhandled_error', { message: err.message, stack: err.stack, path: req.path });
  res.status(500).json({ error: 'internal_server_error', message: 'something went wrong' });
}

module.exports = { AppError, asyncHandler, errorHandler };
