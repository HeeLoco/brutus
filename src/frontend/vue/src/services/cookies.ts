export interface CookieOptions {
  expires?: Date | number // Date object or days from now
  maxAge?: number // seconds
  path?: string
  domain?: string
  secure?: boolean
  httpOnly?: boolean
  sameSite?: 'strict' | 'lax' | 'none'
}

export class CookieService {
  private static isSecureContext(): boolean {
    return window.location.protocol === 'https:' || window.location.hostname === 'localhost'
  }

  static set(name: string, value: string, options: CookieOptions = {}): void {
    let cookieString = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`

    // Set default options for auth cookies
    const defaultOptions: CookieOptions = {
      path: '/',
      secure: this.isSecureContext(),
      sameSite: 'lax',
      maxAge: 86400 // 24 hours in seconds
    }

    const finalOptions = { ...defaultOptions, ...options }

    if (finalOptions.expires) {
      if (typeof finalOptions.expires === 'number') {
        // Convert days to date
        const date = new Date()
        date.setTime(date.getTime() + (finalOptions.expires * 24 * 60 * 60 * 1000))
        cookieString += `; expires=${date.toUTCString()}`
      } else {
        cookieString += `; expires=${finalOptions.expires.toUTCString()}`
      }
    }

    if (finalOptions.maxAge) {
      cookieString += `; max-age=${finalOptions.maxAge}`
    }

    if (finalOptions.path) {
      cookieString += `; path=${finalOptions.path}`
    }

    if (finalOptions.domain) {
      cookieString += `; domain=${finalOptions.domain}`
    }

    if (finalOptions.secure) {
      cookieString += `; secure`
    }

    if (finalOptions.httpOnly) {
      cookieString += `; httponly`
    }

    if (finalOptions.sameSite) {
      cookieString += `; samesite=${finalOptions.sameSite}`
    }

    document.cookie = cookieString
  }

  static get(name: string): string | null {
    const encodedName = encodeURIComponent(name)
    const cookies = document.cookie.split(';')
    
    for (const cookie of cookies) {
      const [cookieName, cookieValue] = cookie.trim().split('=')
      if (cookieName === encodedName) {
        return cookieValue ? decodeURIComponent(cookieValue) : ''
      }
    }
    
    return null
  }

  static remove(name: string, options: Omit<CookieOptions, 'expires' | 'maxAge'> = {}): void {
    this.set(name, '', {
      ...options,
      expires: new Date(0) // Set to past date
    })
  }

  static exists(name: string): boolean {
    return this.get(name) !== null
  }

  // Utility methods for auth-specific cookies
  static setAuthToken(token: string, expiresInSeconds: number = 3600): void {
    // Store actual token in HttpOnly cookie (backend sets this)
    // Frontend only stores auth state for UI
    this.set('auth_state', JSON.stringify({
      isAuthenticated: true,
      expiresAt: Date.now() + (expiresInSeconds * 1000),
      tokenExpiry: expiresInSeconds
    }), {
      secure: true,
      sameSite: 'strict',
      maxAge: expiresInSeconds
    })
  }

  static getAuthToken(): string | null {
    // Frontend cannot access HttpOnly token
    // This method now returns auth state for UI purposes
    const authState = this.get('auth_state')
    if (!authState) return null
    
    try {
      const parsed = JSON.parse(authState)
      return parsed.isAuthenticated && parsed.expiresAt > Date.now() ? 'authenticated' : null
    } catch {
      return null
    }
  }

  static removeAuthToken(): void {
    this.remove('auth_state', { secure: true, sameSite: 'strict' })
    // Backend will remove HttpOnly token via logout endpoint
  }

  static setRefreshToken(token: string, expiresInDays: number = 30): void {
    this.set('refresh_token', token, {
      secure: true,
      sameSite: 'strict',
      httpOnly: false, // MSAL needs to access this
      expires: expiresInDays
    })
  }

  static getRefreshToken(): string | null {
    return this.get('refresh_token')
  }

  static removeRefreshToken(): void {
    this.remove('refresh_token', { secure: true, sameSite: 'strict' })
  }

  // Store user account info
  static setAccountInfo(accountInfo: any): void {
    this.set('account_info', JSON.stringify(accountInfo), {
      secure: true,
      sameSite: 'strict',
      maxAge: 86400 // 24 hours
    })
  }

  static getAccountInfo(): any | null {
    const accountStr = this.get('account_info')
    if (!accountStr) return null
    
    try {
      return JSON.parse(accountStr)
    } catch {
      return null
    }
  }

  static removeAccountInfo(): void {
    this.remove('account_info', { secure: true, sameSite: 'strict' })
  }

  // Clear all auth-related cookies
  static clearAllAuthCookies(): void {
    this.removeAuthToken()
    this.removeRefreshToken()
    this.removeAccountInfo()
    
    // Also clear any MSAL-related cookies
    const allCookies = document.cookie.split(';')
    for (const cookie of allCookies) {
      const cookieName = cookie.split('=')[0].trim()
      if (cookieName.includes('msal') || cookieName.includes('auth') || cookieName.includes('csrf')) {
        this.remove(cookieName)
      }
    }
  }

  // CSRF Token methods
  static setCSRFToken(token: string): void {
    this.set('csrf_token', token, {
      secure: true,
      sameSite: 'strict',
      maxAge: 3600 // 1 hour
    })
  }

  static getCSRFToken(): string | null {
    return this.get('csrf_token')
  }

  static removeCSRFToken(): void {
    this.remove('csrf_token', { secure: true, sameSite: 'strict' })
  }
}