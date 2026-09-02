const defaultApiBaseUrl = 'http://localhost:8080'

const configuredApiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').trim()

export const apiBaseUrl = (configuredApiBaseUrl || defaultApiBaseUrl).replace(/\/+$/, '')
