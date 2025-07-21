import time
from typing import Callable

from fastapi import Request, Response
from starlette.middleware.base import BaseHTTPMiddleware

from api.core.logging import correlation_id, generate_correlation_id, get_logger

logger = get_logger(__name__)


class CorrelationMiddleware(BaseHTTPMiddleware):
    """Middleware to handle request correlation IDs and logging."""
    
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        """Process request with correlation ID and logging."""
        
        # Generate or extract correlation ID
        corr_id = request.headers.get("X-Correlation-ID") or generate_correlation_id()
        correlation_id.set(corr_id)
        
        # Log request start
        start_time = time.time()
        logger.info(f"Request started: {request.method} {request.url.path}")
        
        try:
            # Process request
            response = await call_next(request)
            
            # Calculate processing time
            process_time = time.time() - start_time
            
            # Add correlation ID to response headers
            response.headers["X-Correlation-ID"] = corr_id
            response.headers["X-Process-Time"] = str(process_time)
            
            # Log response
            logger.info(
                f"Request completed: {request.method} {request.url.path} "
                f"- Status: {response.status_code} - Time: {process_time:.3f}s"
            )
            
            return response
            
        except Exception as e:
            # Calculate processing time for errors too
            process_time = time.time() - start_time
            
            logger.error(
                f"Request failed: {request.method} {request.url.path} "
                f"- Error: {str(e)} - Time: {process_time:.3f}s"
            )
            raise