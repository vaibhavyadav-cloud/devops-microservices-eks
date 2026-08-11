from fastapi import Request
from fastapi.responses import JSONResponse
from .logger import logger


class ProductNotFoundError(Exception):
    def __init__(self, product_id: int):
        self.product_id = product_id


class UnauthorizedError(Exception):
    def __init__(self, message: str = "unauthorized"):
        self.message = message


async def product_not_found_handler(request: Request, exc: ProductNotFoundError):
    return JSONResponse(
        status_code=404,
        content={"error": "not_found", "message": f"product {exc.product_id} not found"},
    )


async def unauthorized_handler(request: Request, exc: UnauthorizedError):
    return JSONResponse(
        status_code=401,
        content={"error": "unauthorized", "message": exc.message},
    )


async def unhandled_exception_handler(request: Request, exc: Exception):
    logger.error(f"unhandled_error path={request.url.path} error={exc}")
    return JSONResponse(
        status_code=500,
        content={"error": "internal_server_error", "message": "something went wrong"},
    )


def register_exception_handlers(app):
    app.add_exception_handler(ProductNotFoundError, product_not_found_handler)
    app.add_exception_handler(UnauthorizedError, unauthorized_handler)
    app.add_exception_handler(Exception, unhandled_exception_handler)
