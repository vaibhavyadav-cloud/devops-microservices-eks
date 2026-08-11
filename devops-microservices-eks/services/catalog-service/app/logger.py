import logging
import sys
import json
from datetime import datetime, timezone


class JsonFormatter(logging.Formatter):
    """Same structured-JSON idea as auth-service's logger.js — keeps log
    shape consistent across services so your log pipeline (Loki/CloudWatch)
    can parse all 4 services the same way regardless of language."""

    def format(self, record: logging.LogRecord) -> str:
        entry = {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "level": record.levelname.lower(),
            "service": "catalog-service",
            "message": record.getMessage(),
        }
        if record.exc_info:
            entry["exception"] = self.formatException(record.exc_info)
        return json.dumps(entry)


def get_logger(name: str = "catalog-service") -> logging.Logger:
    logger = logging.getLogger(name)
    if not logger.handlers:
        handler = logging.StreamHandler(sys.stdout)
        handler.setFormatter(JsonFormatter())
        logger.addHandler(handler)
        logger.setLevel(logging.INFO)
    return logger


logger = get_logger()
