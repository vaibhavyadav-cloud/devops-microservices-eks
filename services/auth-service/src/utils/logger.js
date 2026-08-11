// Minimal structured JSON logger. In real production you'd reach for pino/winston,
// but the point here is the *shape* — structured JSON logs are what your
// EKS log pipeline (Fluent Bit -> CloudWatch/Loki) expects, not plain strings.

function log(level, message, meta = {}) {
  const entry = {
    timestamp: new Date().toISOString(),
    level,
    service: 'auth-service',
    message,
    ...meta,
  };
  console.log(JSON.stringify(entry));
}

module.exports = {
  info: (message, meta) => log('info', message, meta),
  warn: (message, meta) => log('warn', message, meta),
  error: (message, meta) => log('error', message, meta),
};
