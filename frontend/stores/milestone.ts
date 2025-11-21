import { defineStore } from 'pinia'

export const useMilestoneStore = defineStore('milestone', {
    state: () => ({
        milestones: [],
        assessments: [],
        summary: null,
        loading: false,
        error: null,
        draftAssessments: {} // LocalStorage sync
    }),

    getters: {
        // Group milestones by pyramid level
        milestonesByLevel: (state) => {
            const grouped = {
                1: { name: 'Sensorik', items: [] },
                2: { name: 'Motorik', items: [] },
                3: { name: 'Persepsi', items: [] },
                4: { name: 'Kognitif', items: [] }
            }

            state.milestones.forEach(m => {
                if (grouped[m.pyramid_level]) {
                    grouped[m.pyramid_level].items.push(m)
                }
            })

            return grouped
        },

        // Group milestones by Denver II domain
        milestonesByDenverDomain: (state) => {
            const grouped = {
                'PS': { name: 'Personal-Social', items: [] },
                'FM': { name: 'Fine Motor-Adaptive', items: [] },
                'L': { name: 'Language', items: [] },
                'GM': { name: 'Gross Motor', items: [] }
            }

            state.milestones.forEach(m => {
                if (m.denver_domain && grouped[m.denver_domain]) {
                    grouped[m.denver_domain].items.push(m)
                }
            })

            return grouped
        },

        // Check if there are unsaved drafts
        hasDrafts: (state) => Object.keys(state.draftAssessments).length > 0
    },

    actions: {
        async fetchMilestones(ageMonths) {
            this.loading = true
            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                const data = await $fetch(`${apiBase}/api/milestones`, {
                    params: { age_months: ageMonths },
                    headers: { 'Authorization': `Bearer ${authStore.token}` }
                })

                this.milestones = data
                return { success: true, data }
            } catch (e) {
                console.error('Failed to fetch milestones:', e)
                this.error = e.message
                return { success: false, error: e.message }
            } finally {
                this.loading = false
            }
        },

        async fetchDenverIIMilestones(ageMonths) {
            this.loading = true
            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                const data = await $fetch(`${apiBase}/api/milestones/denver-ii`, {
                    params: { age_months: ageMonths },
                    headers: { 'Authorization': `Bearer ${authStore.token}` }
                })

                this.milestones = data
                return { success: true, data }
            } catch (e) {
                console.error('Failed to fetch Denver II milestones:', e)
                this.error = e.message
                return { success: false, error: e.message }
            } finally {
                this.loading = false
            }
        },

        async fetchSummary(childId) {
            this.loading = true
            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                const data = await $fetch(`${apiBase}/api/children/${childId}/assessments/summary`, {
                    headers: { 'Authorization': `Bearer ${authStore.token}` }
                })

                this.summary = data
                return { success: true, data }
            } catch (e) {
                console.error('Failed to fetch summary:', e)
                return { success: false, error: e.message }
            } finally {
                this.loading = false
            }
        },

        // Save to local state/storage only
        saveDraft(childId, milestoneId, status, notes = '') {
            if (!this.draftAssessments[childId]) {
                this.draftAssessments[childId] = []
            }

            // Remove existing draft for this milestone if any
            this.draftAssessments[childId] = this.draftAssessments[childId].filter(
                item => item.milestone_id !== milestoneId
            )

            // Add new draft
            this.draftAssessments[childId].push({
                milestone_id: milestoneId,
                status,
                notes
            })

            // Persist to localStorage
            if (process.client) {
                localStorage.setItem('milestone_drafts', JSON.stringify(this.draftAssessments))
            }
        },

        // Sync drafts to backend
        async syncAssessments(childId, assessmentDate) {
            const drafts = this.draftAssessments[childId]
            if (!drafts || drafts.length === 0) {
                return { success: false, error: 'Tidak ada data penilaian untuk disimpan' }
            }

            this.loading = true
            this.error = null
            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                console.log('Syncing assessments:', {
                    childId,
                    assessmentDate,
                    itemCount: drafts.length
                })

                const response = await $fetch(`${apiBase}/api/children/${childId}/assessments/batch`, {
                    method: 'PUT',
                    headers: { 
                        'Authorization': `Bearer ${authStore.token}`,
                        'Content-Type': 'application/json'
                    },
                    body: {
                        assessment_date: assessmentDate,
                        items: drafts
                    }
                })

                console.log('Assessment sync response:', response)

                // Clear drafts after successful sync
                delete this.draftAssessments[childId]
                if (process.client) {
                    localStorage.setItem('milestone_drafts', JSON.stringify(this.draftAssessments))
                }

                // Refresh summary
                await this.fetchSummary(childId)

                return { success: true, data: response }
            } catch (e) {
                console.error('Failed to sync assessments:', e)
                console.error('Error details:', {
                    message: e.message,
                    data: e.data,
                    statusCode: e.statusCode,
                    response: e.response
                })
                
                const errorMessage = e.data?.error || e.data?.message || e.message || 'Gagal menyimpan penilaian'
                this.error = errorMessage
                return { success: false, error: errorMessage }
            } finally {
                this.loading = false
            }
        },

        loadDraftsFromStorage() {
            if (process.client) {
                const stored = localStorage.getItem('milestone_drafts')
                if (stored) {
                    this.draftAssessments = JSON.parse(stored)
                }
            }
        },

        clearState() {
            this.milestones = []
            this.assessments = []
            this.summary = null
            this.draftAssessments = {}
            this.error = null
            if (process.client) {
                localStorage.removeItem('milestone_drafts')
            }
        }
    }
})
