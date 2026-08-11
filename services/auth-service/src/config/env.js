// Centralized env config + validation.
// Fails fast at startup if required vars are missing — in K8s this means
// the pod crashes immediately instead of limping along with undefined values,
// which makes misconfiguration obvious in `kubectl logs` right away.

const required = ['DB_HOST', 'DB_PORT', 'DB_USER', 'DB_PASSWORD', 'DB_NAME', 'JWT_SECRET'];

function loadEnv() {
  const missing = required.filter((key) => !process.env[key]);

  if (missing.length > 0) {
    // In dev, fall back to defaults so `npm run dev` works without a .env file.
    // In production (NODE_ENV=production), refuse to start with missing secrets.
    if (process.env.NODE_ENV === 'production') {
      throw new Error(`Missing required environment variables: ${missing.join(', ')}`);
    }
    console.warn(`[env] Missing vars, using dev defaults: ${missing.join(', ')}`);
  }

  return {
    port: parseInt(process.env.PORT || '3000', 10),
    nodeEnv: process.env.NODE_ENV || 'development',
    db: {
      host: process.env.DB_HOST || 'localhost',
      port: parseInt(process.env.DB_PORT || '5432', 10),
      user: process.env.DB_USER || 'postgres',
      password: process.env.DB_PASSWORD || 'postgres',
      name: process.env.DB_NAME || 'authdb',
    },
    jwtSecret: process.env.JWT_SECRET || 'dev-secret-change-me',
    jwtExpiresIn: process.env.JWT_EXPIRES_IN || '1h',
  };
}

module.exports = loadEnv();
