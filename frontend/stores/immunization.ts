import { defineStore } from 'pinia'

export interface ImmunizationSchedule {
  id: string
  name: string
  name_id?: string
  description?: string
  age_min_days?: number
  age_optimal_days?: number
  age_max_days?: number
  age_min_months?: number
  age_optimal_months?: number
  age_max_months?: number
  dose_number: number
  total_doses?: number
  interval_from_previous_days?: number
  category: string
  priority: string
  is_required: boolean
  notes?: string
  source: string
}

export interface ChildImmunization {
  id: string
  child_id: string
  immunization_schedule_id: string
  given_date: string
  given_at_age_days?: number
  given_at_age_months?: number
  location?: string
  healthcare_facility?: string
  doctor_name?: string
  vaccine_batch_number?: string
  notes?: string
  is_on_schedule?: boolean
  is_catch_up: boolean
  schedule?: ImmunizationSchedule
}

export interface ImmunizationStatus {
  schedule: ImmunizationSchedule
  status: 'pending' | 'completed' | 'overdue' | 'upcoming'
  due_date?: string
  due_age_months?: number
  days_until_due?: number
  days_overdue?: number
  record?: ChildImmunization
}

export interface ImmunizationSummary {
  total: number
  completed: number
  pending: number
  overdue: number
  upcoming: number
}

export interface ImmunizationScheduleResponse {
  child_id: string
  age_months: number
  age_days: number
  immunizations: ImmunizationStatus[]
  summary: ImmunizationSummary
}

export const useImmunizationStore = defineStore('immunization', {
  state: () => ({
    schedule: [] as ImmunizationStatus[],
    summary: null as ImmunizationSummary | null,
    ageMonths: 0,
    ageDays: 0,
    loading: false,
    error: null as string | null,
  }),
  
  getters: {
    hasSchedule: (state) => state.schedule.length > 0,
    
    completedImmunizations: (state) => 
      state.schedule.filter(i => i.status === 'completed'),
    
    pendingImmunizations: (state) => 
      state.schedule.filter(i => i.status === 'pending'),
    
    overdueImmunizations: (state) => 
      state.schedule.filter(i => i.status === 'overdue'),
    
    upcomingImmunizations: (state) => 
      state.schedule.filter(i => i.status === 'upcoming'),
    
    immunizationsByCategory: (state) => {
      const grouped: Record<string, ImmunizationStatus[]> = {}
      state.schedule.forEach(imm => {
        const cat = imm.schedule.category || 'lainnya'
        if (!grouped[cat]) {
          grouped[cat] = []
        }
        grouped[cat].push(imm)
      })
      return grouped
    },
    
    immunizationsByPriority: (state) => {
      const grouped: Record<string, ImmunizationStatus[]> = {
        high: [],
        medium: [],
        low: []
      }
      state.schedule.forEach(imm => {
        const priority = imm.schedule.priority || 'medium'
        if (grouped[priority]) {
          grouped[priority].push(imm)
        }
      })
      return grouped
    },
  },
  
  actions: {
    async fetchSchedule(childId: string) {
      if (!childId) {
        this.schedule = []
        this.summary = null
        return
      }
      
      this.loading = true
      this.error = null
      
      try {
        const authStore = useAuthStore()
        const apiBase = useApiUrl()
        
        const data = await $fetch<ImmunizationScheduleResponse>(
          `${apiBase}/api/children/${childId}/immunizations`,
          {
            headers: {
              'Authorization': `Bearer ${authStore.token}`
            }
          }
        )
        
        this.schedule = data.immunizations || []
        this.summary = data.summary || null
        this.ageMonths = data.age_months || 0
        this.ageDays = data.age_days || 0
      } catch (e: any) {
        console.error('Failed to fetch immunization schedule:', e)
        this.error = e.data?.error || e.message || 'Failed to fetch immunization schedule'
        this.schedule = []
        this.summary = null
      } finally {
        this.loading = false
      }
    },
    
    async recordImmunization(childId: string, record: {
      immunization_schedule_id: string
      given_date: string
      location?: string
      healthcare_facility?: string
      doctor_name?: string
      vaccine_batch_number?: string
      notes?: string
    }) {
      this.loading = true
      this.error = null
      
      try {
        const authStore = useAuthStore()
        const apiBase = useApiUrl()
        
        await $fetch(
          `${apiBase}/api/children/${childId}/immunizations`,
          {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${authStore.token}`,
              'Content-Type': 'application/json'
            },
            body: record
          }
        )
        
        // Refresh schedule after recording
        await this.fetchSchedule(childId)
        
        return { success: true }
      } catch (e: any) {
        console.error('Failed to record immunization:', e)
        this.error = e.data?.error || e.message || 'Failed to record immunization'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },
    
    clearSchedule() {
      this.schedule = []
      this.summary = null
      this.ageMonths = 0
      this.ageDays = 0
      this.error = null
    },
  },
})

