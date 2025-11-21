/**
 * Composable untuk menghitung corrected age untuk bayi prematur
 * Corrected age digunakan untuk bayi prematur sampai usia 24 bulan chronological
 */

export const useCorrectedAge = () => {
  /**
   * Menghitung umur kronologis dalam bulan
   */
  const calculateChronologicalAgeInMonths = (dob: string, measurementDate?: string): number => {
    const birthDate = new Date(dob)
    const targetDate = measurementDate ? new Date(measurementDate) : new Date()
    
    const years = targetDate.getFullYear() - birthDate.getFullYear()
    const months = targetDate.getMonth() - birthDate.getMonth()
    const totalMonths = years * 12 + months
    
    // Adjust jika hari belum tercapai di bulan ini
    if (targetDate.getDate() < birthDate.getDate()) {
      return Math.max(0, totalMonths - 1)
    }
    
    return Math.max(0, totalMonths)
  }

  /**
   * Menghitung corrected age untuk bayi prematur
   * Corrected age = Chronological age - (40 weeks - gestational age in weeks)
   * Hanya berlaku sampai 24 bulan chronological age (730 days)
   */
  const calculateCorrectedAge = (
    dob: string,
    isPremature: boolean,
    gestationalAgeWeeks: number | null | undefined,
    measurementDate?: string
  ): {
    chronologicalMonths: number
    correctedMonths: number | null
    useCorrected: boolean
    chronologicalDisplay: string
    correctedDisplay: string | null
  } => {
    const birthDate = new Date(dob)
    const targetDate = measurementDate ? new Date(measurementDate) : new Date()
    
    // Calculate chronological age in months
    const chronologicalMonths = calculateChronologicalAgeInMonths(dob, measurementDate)
    
    // Calculate chronological age in days for 24-month check
    const daysDiff = Math.floor((targetDate.getTime() - birthDate.getTime()) / (1000 * 60 * 60 * 24))
    
    // Jika tidak prematur atau gestational age tidak ada, return chronological only
    if (!isPremature || !gestationalAgeWeeks || gestationalAgeWeeks <= 0) {
      return {
        chronologicalMonths,
        correctedMonths: null,
        useCorrected: false,
        chronologicalDisplay: formatAge(chronologicalMonths),
        correctedDisplay: null
      }
    }
    
    // Jika chronological age >= 24 months (730 days), tidak pakai corrected age
    if (daysDiff >= 730) {
      return {
        chronologicalMonths,
        correctedMonths: chronologicalMonths,
        useCorrected: false,
        chronologicalDisplay: formatAge(chronologicalMonths),
        correctedDisplay: null
      }
    }
    
    // Calculate weeks premature
    // Full term is 40 weeks, so weeks premature = 40 - gestational_age
    const weeksPremature = 40 - gestationalAgeWeeks
    if (weeksPremature <= 0) {
      return {
        chronologicalMonths,
        correctedMonths: chronologicalMonths,
        useCorrected: false,
        chronologicalDisplay: formatAge(chronologicalMonths),
        correctedDisplay: null
      }
    }
    
    // Convert weeks to days (1 week = 7 days)
    const daysToSubtract = weeksPremature * 7
    
    // Calculate corrected age in days
    const correctedDays = daysDiff - daysToSubtract
    
    // Calculate corrected age in months (approximate)
    // Average days per month = 30.44
    const correctedMonths = Math.max(0, Math.floor(correctedDays / 30.44))
    
    return {
      chronologicalMonths,
      correctedMonths,
      useCorrected: true,
      chronologicalDisplay: formatAge(chronologicalMonths),
      correctedDisplay: formatAge(correctedMonths)
    }
  }

  /**
   * Format age dalam format yang mudah dibaca
   */
  const formatAge = (months: number): string => {
    const years = Math.floor(months / 12)
    const remainingMonths = months % 12
    
    if (years > 0) {
      return `${years} tahun ${remainingMonths} bulan`
    }
    return `${months} bulan`
  }

  return {
    calculateCorrectedAge,
    calculateChronologicalAgeInMonths,
    formatAge
  }
}

