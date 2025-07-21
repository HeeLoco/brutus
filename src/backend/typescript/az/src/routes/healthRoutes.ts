import { Router } from 'express';
import { HealthController } from '@/controllers/healthController';
import { asyncHandler } from '@/utils/asyncHandler';

const router = Router();
const healthController = new HealthController();

/**
 * @route   GET /health
 * @desc    Health check endpoint
 * @access  Public
 */
router.get('/health', asyncHandler(healthController.checkHealth.bind(healthController)));

export default router;