/**
 * Utility functions untuk validation
 * Centralized validation logic
 */

/**
 * Validasi zscore untuk status badge
 */
export const getZScoreStatus = (zscore: number | null | undefined): {
  level: 'critical' | 'warning' | 'normal' | 'invalid'
  label: string
} => {
  if (typeof zscore !== 'number' || isNaN(zscore)) {
    return { level: 'invalid', label: 'Invalid' }
  }
  
  if (zscore < -2 || zscore > 2) {
    return { level: 'critical', label: 'Critical' }
  }
  
  if (zscore < -1 || zscore > 1) {
    return { level: 'warning', label: 'Warning' }
  }
  
  return { level: 'normal', label: 'Normal' }
}

/**
 * Get CSS class untuk zscore badge
 */
export const getZScoreBadgeClass = (zscore: number | null | undefined): string => {
  const status = getZScoreStatus(zscore)
  
  const classes = {
    critical: 'bg-red-100 text-red-700',
    warning: 'bg-amber-100 text-amber-700',
    normal: 'bg-emerald-100 text-emerald-700',
    invalid: 'bg-gray-100 text-gray-700'
  }
  
  return classes[status.level]
}

/**
 * Get CSS class untuk zscore text
 */
export const getZScoreTextClass = (zscore: number | null | undefined): string => {
  const status = getZScoreStatus(zscore)
  
  const classes = {
    critical: 'bg-red-50 text-red-700 border border-red-200',
    warning: 'bg-amber-50 text-amber-700 border border-amber-200',
    normal: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
    invalid: 'bg-gray-50 text-gray-700 border border-gray-200'
  }
  
  return classes[status.level]
}

/**
 * Get border color berdasarkan zscore
 */
export const getZScoreBorderColor = (
  weightZ: number | null | undefined,
  heightZ: number | null | undefined
): string => {
  const weightStatus = getZScoreStatus(weightZ)
  const heightStatus = getZScoreStatus(heightZ)
  
  // Jika salah satu critical, return red
  if (weightStatus.level === 'critical' || heightStatus.level === 'critical') {
    return 'border-red-500'
  }
  
  // Jika salah satu warning, return amber
  if (weightStatus.level === 'warning' || heightStatus.level === 'warning') {
    return 'border-amber-500'
  }
  
  // Jika salah satu invalid, return gray
  if (weightStatus.level === 'invalid' || heightStatus.level === 'invalid') {
    return 'border-gray-500'
  }
  
  return 'border-emerald-500'
}

