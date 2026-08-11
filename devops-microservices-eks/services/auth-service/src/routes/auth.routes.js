const express = require('express');
const { register, login, verify } = require('../controllers/auth.controller');
const { validateRegister, validateLogin } = require('../middleware/validate');
const { asyncHandler } = require('../middleware/errorHandler');

const router = express.Router();

router.post('/register', validateRegister, register);
router.post('/login', validateLogin, login);
router.get('/verify', asyncHandler(verify));

module.exports = router;
