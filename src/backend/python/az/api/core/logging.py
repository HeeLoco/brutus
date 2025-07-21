import logging
import sys
import uuid
from contextvars import ContextVar
from typing import Optional

from api.core.config import get_settings

# Context variable to store correlation ID
correlation_id: ContextVar[Optional[str]] = ContextVar('correlation_id', default=None)


class CorrelationFilter(logging.Filter):
    """Add correlation ID to log records."""
    
    def filter(self, record):
        record.correlation_id = correlation_id.get() or "no-correlation-id"
        return True


def setup_logging():
    """Setup application logging configuration."""
    settings = get_settings()
    
    # Create formatter
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - [%(correlation_id)s] - %(message)s"
    )
    
    # Setup handler
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(formatter)
    handler.addFilter(CorrelationFilter())
    
    # Setup root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(getattr(logging, settings.log_level.upper()))
    root_logger.addHandler(handler)
    
    # Setup specific loggers
    logging.getLogger("uvicorn.access").handlers = []
    logging.getLogger("azure").setLevel(logging.WARNING)


def get_logger(name: str) -> logging.Logger:
    """Get logger instance with correlation ID support."""
    return logging.getLogger(name)


def generate_correlation_id() -> str:
    """Generate a new correlation ID."""
    return str(uuid.uuid4())[:8]