"""
Verifies JWTs issued by the Auth service. Both services share the same
JWT_SECRET, so Catalog verifies signatures locally — no call to Auth
happens per-request. This is the standard stateless-auth pattern in
microservices: trust is established once (shared secret), not per-request.
"""
from fastapi import Depends
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from jose import JWTError, jwt

from ..config import settings
from ..exceptions import UnauthorizedError

bearer_scheme = HTTPBearer()


def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(bearer_scheme),
) -> dict:
    token = credentials.credentials
    try:
        payload = jwt.decode(token, settings.jwt_secret, algorithms=[settings.jwt_algorithm])
    except JWTError:
        raise UnauthorizedError("invalid or expired token")

    # Auth service signs { userId, email } — see auth.controller.js login()
    if "userId" not in payload:
        raise UnauthorizedError("malformed token")

    return payload
