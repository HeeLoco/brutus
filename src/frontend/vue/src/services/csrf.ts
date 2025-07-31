import { CookieService } from './cookies'
import { logger } from './logger'

export class CSRFService {
  private static csrfToken: string | null = null

  // Generate a random CSRF token
  private static generateToken(): string {
    const array = new Uint8Array(32)
    crypto.getRandomValues(array)
    return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('')
  }

  // Initialize CSRF token on app start
  static async initialize(): Promise<void> {
    try {
      // Check if we have a valid token in cookies
      const existingToken = CookieService.getCSRFToken()
      
      if (existingToken) {
        this.csrfToken = existingToken
        logger.info('CSRF token restored from cookie')
        return
      }

      // Generate new token
      await this.refreshToken()
    } catch (error) {
      logger.logError(error, 'Failed to initialize CSRF token')
    }
  }

  // Get current CSRF token, refresh if needed
  static async getToken(): Promise<string> {
    if (!this.csrfToken) {
      await this.refreshToken()
    }
    
    return this.csrfToken || ''
  }

  // Refresh CSRF token
  static async refreshToken(): Promise<void> {
    try {
      // In a real implementation, you might fetch this from the backend
      // For now, we'll generate client-side and use double-submit pattern
      this.csrfToken = this.generateToken()
      CookieService.setCSRFToken(this.csrfToken)
      logger.info('CSRF token refreshed')
    } catch (error) {
      logger.logError(error, 'Failed to refresh CSRF token')
      throw error
    }
  }

  // Clear CSRF token
  static clearToken(): void {
    this.csrfToken = null
    CookieService.removeCSRFToken()
    logger.info('CSRF token cleared')
  }

  // Get headers with CSRF token
  static async getHeaders(): Promise<Record<string, string>> {
    const token = await this.getToken()
    return {
      'X-CSRF-Token': token,
      'X-Requested-With': 'XMLHttpRequest' // Additional CSRF protection
    }
  }

  // Validate response for CSRF errors
  static handleCSRFError(response: Response): boolean {
    if (response.status === 403) {
      // Check if it's a CSRF error
      const csrfError = response.headers.get('X-CSRF-Error')
      if (csrfError) {
        logger.warn('CSRF token validation failed, refreshing token')
        this.refreshToken()
        return true
      }
    }
    return false
  }
}