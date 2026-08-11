from contextlib import asynccontextmanager
from fastapi import FastAPI
from sqlalchemy import text

from app.database import init_db, engine
from app.routers import products
from app.exceptions import register_exception_handlers
from app.logger import logger
from app.config import settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup: create tables if they don't exist yet.
    # (Real production would use Alembic migrations instead of create_all —
    # noting that here deliberately since it's a common interview question.)
    try:
        init_db()
        logger.info("db_initialized")
    except Exception as e:
        logger.error(f"db_init_failed_starting_anyway error={e}")
    yield
    # Shutdown: dispose the connection pool cleanly.
    engine.dispose()
    logger.info("engine_disposed")


app = FastAPI(title="catalog-service", lifespan=lifespan)
register_exception_handlers(app)
app.include_router(products.router)


# --- Health checks (same liveness/readiness split as auth-service) ---
@app.get("/health/live")
def live():
    return {"status": "ok"}


@app.get("/health/ready")
def ready():
    try:
        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))
        return {"status": "ready"}
    except Exception as e:
        from fastapi.responses import JSONResponse
        return JSONResponse(status_code=503, content={"status": "not_ready", "error": str(e)})


@app.get("/health")
def health():
    return {"status": "ok", "service": "catalog-service"}
