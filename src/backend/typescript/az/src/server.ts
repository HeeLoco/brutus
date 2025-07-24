import App from './app';
import config from '@/config/environment';
import logger from '@/config/logger';

const app = new App().getApp();

// Graceful shutdown handler
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully');
  process.exit(0);
});

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully');
  process.exit(0);
});

// Unhandled promise rejection handler
process.on('unhandledRejection', (reason, promise) => {
  logger.error('Unhandled Rejection at:', { promise, reason });
  process.exit(1);
});

// Uncaught exception handler
process.on('uncaughtException', (error) => {
  logger.error('Uncaught Exception:', { error: error.message, stack: error.stack });
  process.exit(1);
});

// Start server
const server = app.listen(config.PORT, config.HOST, () => {
  logger.info(`🚀 Azure API Server (TypeScript) started`, {
    host: config.HOST,
    port: config.PORT,
    environment: config.NODE_ENV,
    apiPrefix: config.API_PREFIX,
    corsOrigins: config.CORS_ORIGINS,
    logLevel: config.LOG_LEVEL
  });
  
  logger.info(`📋 API endpoints:`, {
    root: `http://${config.HOST}:${config.PORT}/`,
    health: `http://${config.HOST}:${config.PORT}/health`,
    api: `http://${config.HOST}:${config.PORT}${config.API_PREFIX}`,
    resourceGroups: `http://${config.HOST}:${config.PORT}${config.API_PREFIX}/resource-groups`
  });
});

export default server;