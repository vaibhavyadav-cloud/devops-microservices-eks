const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const { pool } = require('../../db');
const env = require('../config/env');
const logger = require('../utils/logger');
const { AppError, asyncHandler } = require('../middleware/errorHandler');

const register = asyncHandler(async (req, res) => {
  const { email, password } = req.body;
  const hash = await bcrypt.hash(password, 10);

  const result = await pool.query(
    'INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at',
    [email, hash]
  );

  logger.info('user_registered', { userId: result.rows[0].id });
  res.status(201).json(result.rows[0]);
});

const login = asyncHandler(async (req, res) => {
  const { email, password } = req.body;

  const result = await pool.query('SELECT * FROM users WHERE email = $1', [email]);
  const user = result.rows[0];

  // Same error for "no user" and "wrong password" — don't leak which one it was,
  // that's a user-enumeration vector.
  if (!user) throw new AppError(401, 'invalid_credentials', 'invalid email or password');

  const valid = await bcrypt.compare(password, user.password_hash);
  if (!valid) throw new AppError(401, 'invalid_credentials', 'invalid email or password');

  const token = jwt.sign({ userId: user.id, email: user.email }, env.jwtSecret, {
    expiresIn: env.jwtExpiresIn,
  });

  logger.info('user_logged_in', { userId: user.id });
  res.json({ token, expiresIn: env.jwtExpiresIn });
});

const verify = (req, res) => {
  const authHeader = req.headers.authorization || '';
  const token = authHeader.replace('Bearer ', '');

  if (!token) throw new AppError(401, 'missing_token', 'authorization token required');

  try {
    const payload = jwt.verify(token, env.jwtSecret);
    res.json({ valid: true, payload });
  } catch {
    throw new AppError(401, 'invalid_token', 'token is invalid or expired');
  }
};

module.exports = { register, login, verify };
