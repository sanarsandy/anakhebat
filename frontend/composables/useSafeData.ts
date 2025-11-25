/**
 * Composable untuk safe data access dan validation
 * Mengurangi duplikasi kode dan meningkatkan konsistensi
 */

/**
 * Validasi apakah value adalah number yang valid
 */
export const isValidNumber = (value: any): value is number => {
  return typeof value === 'number' && !isNaN(value) && isFinite(value)
}

/**
 * Safe access untuk zscore dengan fallback
 */
export const safeZScore = (zscore: any, fallback: number = 0): number => {
  return isValidNumber(zscore) ? zscore : fallback
}

/**
 * Format zscore dengan validasi
 */
export const formatZScore = (zscore: any, decimals: number = 2): string => {
  if (!isValidNumber(zscore)) return '-'
  return zscore.toFixed(decimals)
}

/**
 * Validasi object dan nested properties
 */
export const hasNestedProperty = (obj: any, ...path: string[]): boolean => {
  if (!obj || typeof obj !== 'object') return false
  
  let current = obj
  for (const key of path) {
    if (current == null || typeof current !== 'object' || !(key in current)) {
      return false
    }
    current = current[key]
  }
  return true
}

/**
 * Safe get nested property dengan fallback
 */
export const safeGet = <T>(obj: any, path: string[], fallback: T): T => {
  if (!hasNestedProperty(obj, ...path)) return fallback
  
  let current = obj
  for (const key of path) {
    current = current[key]
  }
  return current as T
}

/**
 * Validasi string yang tidak kosong
 */
export const isValidString = (value: any): value is string => {
  return typeof value === 'string' && value.trim().length > 0
}

/**
 * Safe format date dengan error handling
 */
export const safeFormatDate = (dateString: any, locale: string = 'id-ID'): string => {
  if (!dateString) return ''
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return String(dateString)
    
    return date.toLocaleDateString(locale, { 
      day: 'numeric', 
      month: 'long', 
      year: 'numeric' 
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return String(dateString || '')
  }
}

/**
 * Filter array dengan validasi item
 */
export const filterValidItems = <T>(
  items: T[] | null | undefined,
  validator: (item: T) => boolean
): T[] => {
  if (!Array.isArray(items)) return []
  return items.filter(validator)
}

