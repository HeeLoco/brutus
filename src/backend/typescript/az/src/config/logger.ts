import winston from 'winston';
import config from './environment';

const logFormat = winston.format.combine(
  winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
  winston.format.errors({ stack: true }),
  winston.format.printf(({ level, message, timestamp, correlationId, stack, ...meta }) => {
    const metaString = Object.keys(meta).length ? JSON.stringify(meta) : '';
    const correlation = correlationId ? `[${correlationId}]` : '';
    const logMessage = `${timestamp} ${level.toUpperCase()} ${correlation} ${message}`;
    
    return stack 
      ? `${logMessage}\n${stack}${metaString ? ` ${metaString}` : ''}`
      : `${logMessage}${metaString ? ` ${metaString}` : ''}`;
  })
);

const jsonFormat = winston.format.combine(
  winston.format.timestamp(),
  winston.format.errors({ stack: true }),
  winston.format.json()
);

const logger = winston.createLogger({
  level: config.LOG_LEVEL,
  format: config.LOG_FORMAT === 'json' ? jsonFormat : logFormat,
  defaultMeta: { service: 'azure-api-typescript' },
  transports: [
    new winston.transports.Console({
      handleExceptions: true,
      handleRejections: true
    })
  ],
  exitOnError: false
});

// Create a stream object for morgan middleware
export const loggerStream = {
  write: (message: string): void => {
    logger.info(message.trim());
  }
};

export class Logger {
  private correlationId?: string;

  constructor(correlationId?: string) {
    this.correlationId = correlationId;
  }

  private log(level: string, message: string, meta?: any): void {
    logger.log(level, message, { 
      correlationId: this.correlationId, 
      ...meta 
    });
  }

  error(message: string, meta?: any): void {
    this.log('error', message, meta);
  }

  warn(message: string, meta?: any): void {
    this.log('warn', message, meta);
  }

  info(message: string, meta?: any): void {
    this.log('info', message, meta);
  }

  debug(message: string, meta?: any): void {
    this.log('debug', message, meta);
  }
}

export const createLogger = (correlationId?: string): Logger => {
  return new Logger(correlationId);
};

export default logger;