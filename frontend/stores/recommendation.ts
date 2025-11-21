import { defineStore } from 'pinia'

export interface StimulationContent {
  id: string
  milestone_id?: string
  category: string
  title: string
  description?: string
  content_type: 'video' | 'article'
  url: string
  thumbnail_url?: string
  age_min_months?: number
  age_max_months?: number
}

export interface Recommendation {
  content: StimulationContent
  reason: string
  priority: 'high' | 'medium' | 'low'
  related_milestone?: any
}

export interface RecommendationsResponse {
  child_id: string
  age_months: number
  recommendations: Recommendation[]
}

export const useRecommendationStore = defineStore('recommendation', {
  state: () => ({
    recommendations: [] as Recommendation[],
    loading: false,
    error: null as string | null,
    ageMonths: 0,
  }),
  
  getters: {
    hasRecommendations: (state) => state.recommendations.length > 0,
    
    highPriorityRecommendations: (state) => 
      state.recommendations.filter(r => r.priority === 'high'),
    
    mediumPriorityRecommendations: (state) => 
      state.recommendations.filter(r => r.priority === 'medium'),
    
    lowPriorityRecommendations: (state) => 
      state.recommendations.filter(r => r.priority === 'low'),
    
    recommendationsByCategory: (state) => {
      const grouped: Record<string, Recommendation[]> = {}
      state.recommendations.forEach(rec => {
        const cat = rec.content.category
        if (!grouped[cat]) {
          grouped[cat] = []
        }
        grouped[cat].push(rec)
      })
      return grouped
    },
  },
  
  actions: {
    async fetchRecommendations(childId: string) {
      if (!childId) {
        this.recommendations = []
        return
      }
      
      this.loading = true
      this.error = null
      
      try {
        const authStore = useAuthStore()
        const apiBase = useApiUrl()
        
        const data = await $fetch<RecommendationsResponse>(
          `${apiBase}/api/children/${childId}/recommendations`,
          {
            headers: {
              'Authorization': `Bearer ${authStore.token}`
            }
          }
        )
        
        this.recommendations = data.recommendations || []
        this.ageMonths = data.age_months || 0
      } catch (e: any) {
        console.error('Failed to fetch recommendations:', e)
        this.error = e.data?.error || e.message || 'Failed to fetch recommendations'
        this.recommendations = []
      } finally {
        this.loading = false
      }
    },
    
    clearRecommendations() {
      this.recommendations = []
      this.error = null
      this.ageMonths = 0
    },
  },
})

