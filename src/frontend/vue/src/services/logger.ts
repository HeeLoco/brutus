/**
 * Centralized logging service for the frontend application
 * Sends logs to backend container logs and provides fallback to console
 * 
 * @example
 * ```typescript
 * import { logger } from '../services/logger'
 * 
 * // Basic logging
 * logger.info('User logged in successfully')
 * logger.error('Failed to load data')
 * 
 * // Logging with context
 * logger.info('Resource created', { name: 'my-resource', id: '123' })
 * 
 * // Error logging with full context
 * try {
 *   await apiCall()
 * } catch (error) {
 *   logger.logError(error, 'API call failed', { endpoint: '/api/data' })
 * }
 * ```
 */

/**
 * Available log levels in order of severity
 */
export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

/**
 * Structure of a log entry that gets sent to the backend
 */
export interface LogEntry {
  /** ISO 8601 timestamp of when the log was created */
  timestamp: string;
  /** Severity level of the log entry */
  level: LogLevel;
  /** Human-readable message describing the event */
  message: string;
  /** Optional additional data providing context about the event */
  context?: Record<string, any>;
  /** Source identifier, always 'frontend' for frontend logs */
  source: 'frontend';
  /** Unique identifier for the current user session */
  sessionId: string;
  /** Browser user agent string */
  userAgent: string;
  /** Current page URL where the log was generated */
  url: string;
}

/**
 * Centralized logging service that handles all frontend logging operations.
 * Implements singleton pattern to ensure consistent logging across the application.
 */
class LoggerService {
  /** Unique identifier for the current browser session */
  private sessionId: string;
  /** Base URL for the backend API where logs are sent */
  private backendUrl: string;
  /** Whether to also log to browser console (always true for debugging) */
  private fallbackToConsole: boolean;
  /** Queue of log entries waiting to be sent to backend */
  private logQueue: LogEntry[] = [];
  /** Flag to prevent concurrent queue flushing operations */
  private isFlushingQueue = false;

  /**
   * Initialize the logger service with configuration from environment variables
   */
  constructor() {
    this.sessionId = this.generateSessionId();
    this.backendUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api/v1';
    this.fallbackToConsole = true; // Always fallback to console for immediate debugging
  }

  /**
   * Generate a unique session identifier for tracking user sessions
   * @returns A unique session ID in the format 'frontend-{timestamp}-{random}'
   */
  private generateSessionId(): string {
    return `frontend-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  /**
   * Create a structured log entry with all required metadata
   * @param level - The severity level of the log entry
   * @param message - Human-readable description of the event
   * @param context - Optional additional data for context
   * @returns A complete LogEntry object ready for transmission
   */
  private createLogEntry(level: LogLevel, message: string, context?: Record<string, any>): LogEntry {
    return {
      timestamp: new Date().toISOString(),
      level,
      message,
      context,
      source: 'frontend',
      sessionId: this.sessionId,
      userAgent: navigator.userAgent,
      url: window.location.href
    };
  }

  /**
   * Send a log entry to the backend API
   * @param logEntry - The log entry to transmit
   * @throws Does not throw - errors are silently handled to prevent logging loops
   */
  private async sendToBackend(logEntry: LogEntry): Promise<void> {
    try {
      const response = await fetch(`${this.backendUrl}/logs`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(logEntry),
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    } catch (error) {
      // Silently fail to avoid infinite logging loops
      // Backend logging failed, but console fallback will still work
    }
  }

  private logToConsole(level: LogLevel, message: string, context?: Record<string, any>): void {
    const timestamp = new Date().toISOString();
    const prefix = `[${timestamp}] [${level.toUpperCase()}] [Session: ${this.sessionId.split('-')[1]}]`;
    
    switch (level) {
      case 'debug':
        console.debug(prefix, message, context || '');
        break;
      case 'info':
        console.info(prefix, message, context || '');
        break;
      case 'warn':
        console.warn(prefix, message, context || '');
        break;
      case 'error':
        console.error(prefix, message, context || '');
        break;
    }
  }

  private async log(level: LogLevel, message: string, context?: Record<string, any>): Promise<void> {
    const logEntry = this.createLogEntry(level, message, context);

    // Always log to console for immediate visibility
    if (this.fallbackToConsole) {
      this.logToConsole(level, message, context);
    }

    // Queue log entry for backend transmission
    this.logQueue.push(logEntry);

    // Attempt to flush queue (non-blocking)
    this.flushQueue();
  }

  private async flushQueue(): Promise<void> {
    if (this.isFlushingQueue || this.logQueue.length === 0) {
      return;
    }

    this.isFlushingQueue = true;

    try {
      // Send logs in batches to avoid overwhelming the backend
      const batchSize = 10;
      while (this.logQueue.length > 0) {
        const batch = this.logQueue.splice(0, batchSize);
        
        // Send each log entry (we could optimize this to send batches)
        for (const logEntry of batch) {
          await this.sendToBackend(logEntry);
        }
      }
    } catch (error) {
      // If batch sending fails, we've already logged to console
      // so user debugging is not affected
    } finally {
      this.isFlushingQueue = false;
    }
  }

  // Public logging methods
  
  /**
   * Log a debug message - detailed diagnostic information
   * @param message - Description of the debug event
   * @param context - Optional additional data for debugging
   */
  debug(message: string, context?: Record<string, any>): void {
    this.log('debug', message, context);
  }

  /**
   * Log an informational message - general application events
   * @param message - Description of the informational event
   * @param context - Optional additional data for context
   */
  info(message: string, context?: Record<string, any>): void {
    this.log('info', message, context);
  }

  /**
   * Log a warning message - potentially problematic situations
   * @param message - Description of the warning
   * @param context - Optional additional data for context
   */
  warn(message: string, context?: Record<string, any>): void {
    this.log('warn', message, context);
  }

  /**
   * Log an error message - error conditions that occurred
   * @param message - Description of the error
   * @param context - Optional additional data for debugging
   */
  error(message: string, context?: Record<string, any>): void {
    this.log('error', message, context);
  }

  /**
   * Convenience method for logging errors with full context and stack traces
   * @param error - The error object or unknown value that was thrown
   * @param message - Optional custom message, defaults to 'An error occurred'
   * @param context - Optional additional context data
   */
  logError(error: Error | unknown, message?: string, context?: Record<string, any>): void {
    const errorMessage = message || 'An error occurred';
    const errorContext = {
      ...context,
      error: {
        name: error instanceof Error ? error.name : 'UnknownError',
        message: error instanceof Error ? error.message : String(error),
        stack: error instanceof Error ? error.stack : undefined,
      }
    };
    
    this.error(errorMessage, errorContext);
  }

  /**
   * Force flush all pending logs to the backend
   * Useful before page unload to ensure logs are not lost
   * @returns Promise that resolves when all logs have been sent
   */
  async flush(): Promise<void> {
    await this.flushQueue();
  }
}

/**
 * Singleton instance of the logger service.
 * Import and use this instance throughout your application for consistent logging.
 * 
 * @example
 * ```typescript
 * import { logger } from '../services/logger'
 * logger.info('Application started')
 * ```
 */
export const logger = new LoggerService();

// Setup page unload handler to flush logs before the page is closed
window.addEventListener('beforeunload', () => {
  logger.flush();
});