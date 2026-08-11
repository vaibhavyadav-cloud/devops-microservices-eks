// Small hand-rolled validation middleware (no external lib) so the
// request/response contract is explicit and easy to extend later.

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function validateRegister(req, res, next) {
  const { email, password } = req.body || {};
  const errors = [];

  if (!email || !EMAIL_RE.test(email)) errors.push('valid email is required');
  if (!password || password.length < 8) errors.push('password must be at least 8 characters');

  if (errors.length) {
    return res.status(400).json({ error: 'validation_failed', details: errors });
  }
  next();
}

function validateLogin(req, res, next) {
  const { email, password } = req.body || {};
  const errors = [];

  if (!email) errors.push('email is required');
  if (!password) errors.push('password is required');

  if (errors.length) {
    return res.status(400).json({ error: 'validation_failed', details: errors });
  }
  next();
}

module.exports = { validateRegister, validateLogin };
