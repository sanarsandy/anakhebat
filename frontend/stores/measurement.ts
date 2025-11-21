import { defineStore } from 'pinia'

export const useMeasurementStore = defineStore('measurement', {
    state: () => ({
        measurements: [],
        latestMeasurement: null,
        loading: false,
        error: null
    }),

    getters: {
        hasMeasurements: (state) => state.measurements.length > 0,

        sortedMeasurements: (state) => {
            return [...state.measurements].sort((a, b) =>
                new Date(b.measurement_date) - new Date(a.measurement_date)
            )
        }
    },

    actions: {
        async fetchMeasurements(childId) {
            this.loading = true
            this.error = null

            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                console.log('Fetching measurements for child:', childId)
                const data = await $fetch(`${apiBase}/api/children/${childId}/measurements`, {
                    headers: {
                        'Authorization': `Bearer ${authStore.token}`
                    }
                })

                console.log('Measurements fetched:', data)
                console.log('Measurements count:', Array.isArray(data) ? data.length : 'not an array')
                
                // Ensure data is an array
                this.measurements = Array.isArray(data) ? data : []
                return { success: true, data }
            } catch (e) {
                console.error('Failed to fetch measurements:', e)
                console.error('Error details:', {
                    message: e.message,
                    data: e.data,
                    statusCode: e.statusCode
                })
                this.error = e.message
                this.measurements = []
                return { success: false, error: e.message }
            } finally {
                this.loading = false
            }
        },

        async fetchLatestMeasurement(childId) {
            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                console.log('Fetching latest measurement for child:', childId)
                const data = await $fetch(`${apiBase}/api/children/${childId}/measurements/latest`, {
                    headers: {
                        'Authorization': `Bearer ${authStore.token}`
                    }
                })

                console.log('Latest measurement fetched:', data)
                this.latestMeasurement = data
                return { success: true, data }
            } catch (e) {
                console.error('Failed to fetch latest measurement:', e)
                console.error('Error details:', {
                    message: e.message,
                    data: e.data,
                    statusCode: e.statusCode
                })
                // If 404, that's okay - just means no measurements yet
                if (e.statusCode === 404) {
                    this.latestMeasurement = null
                    return { success: true, data: null }
                }
                return { success: false, error: e.message }
            }
        },

        async addMeasurement(childId, measurementData) {
            this.loading = true
            this.error = null

            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                console.log('Adding measurement for child:', childId)
                console.log('Measurement data:', measurementData)
                console.log('API URL:', `${apiBase}/api/children/${childId}/measurements`)

                const data = await $fetch(`${apiBase}/api/children/${childId}/measurements`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${authStore.token}`,
                        'Content-Type': 'application/json'
                    },
                    body: measurementData
                })

                console.log('Measurement added successfully:', data)

                // Refresh measurements list
                await this.fetchMeasurements(childId)
                await this.fetchLatestMeasurement(childId)

                return { success: true, data }
            } catch (e) {
                console.error('Failed to add measurement:', e)
                console.error('Error details:', {
                    message: e.message,
                    data: e.data,
                    statusCode: e.statusCode,
                    response: e.response
                })

                // Extract error message from various possible error formats
                const errorMessage = e.data?.error || e.data?.message || e.message || 'Gagal menambahkan pengukuran'
                this.error = errorMessage
                return { success: false, error: errorMessage }
            } finally {
                this.loading = false
            }
        },

        async updateMeasurement(childId, measurementId, measurementData) {
            this.loading = true
            this.error = null

            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                console.log('Updating measurement:', measurementId, 'for child:', childId)
                console.log('Update data:', measurementData)

                const data = await $fetch(`${apiBase}/api/children/${childId}/measurements/${measurementId}`, {
                    method: 'PUT',
                    headers: {
                        'Authorization': `Bearer ${authStore.token}`,
                        'Content-Type': 'application/json'
                    },
                    body: measurementData
                })

                console.log('Measurement updated successfully:', data)

                // Refresh measurements list and latest measurement
                await this.fetchMeasurements(childId)
                await this.fetchLatestMeasurement(childId)

                return { success: true, data }
            } catch (e) {
                console.error('Failed to update measurement:', e)
                console.error('Error details:', e.data)

                const errorMessage = e.data?.error || e.data?.message || e.message || 'Gagal memperbarui pengukuran'
                this.error = errorMessage
                return { success: false, error: errorMessage }
            } finally {
                this.loading = false
            }
        },

        async deleteMeasurement(childId, measurementId) {
            this.loading = true
            this.error = null

            try {
                const apiBase = useApiUrl()
                const authStore = useAuthStore()

                console.log('Deleting measurement:', measurementId, 'for child:', childId)

                const response = await $fetch(`${apiBase}/api/children/${childId}/measurements/${measurementId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${authStore.token}`
                    }
                })

                console.log('Delete response:', response)

                // Refresh measurements list
                await this.fetchMeasurements(childId)
                await this.fetchLatestMeasurement(childId)

                return { success: true, data: response }
            } catch (e) {
                console.error('Failed to delete measurement:', e)
                console.error('Error details:', {
                    message: e.message,
                    data: e.data,
                    statusCode: e.statusCode
                })
                
                const errorMessage = e.data?.error || e.data?.message || e.message || 'Gagal menghapus pengukuran'
                this.error = errorMessage
                return { success: false, error: errorMessage }
            } finally {
                this.loading = false
            }
        },

        clearMeasurements() {
            this.measurements = []
            this.latestMeasurement = null
            this.error = null
        }
    }
})
