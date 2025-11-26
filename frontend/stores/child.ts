import { defineStore } from 'pinia'

export const useChildStore = defineStore('child', {
    state: () => ({
        children: [] as any[],
        selectedChild: null as any | null,
    }),
    getters: {
        hasChildren: (state) => state.children.length > 0,
    },
    actions: {
        async fetchChildren() {
            const authStore = useAuthStore()
            const apiBase = useApiUrl()

            if (!authStore.token) {
                // No token, clear state
                this.clearState()
                return
            }

            try {
                const data = await $fetch(`${apiBase}/api/children`, {
                    headers: {
                        Authorization: `Bearer ${authStore.token}`
                    }
                })

                if (data) {
                    this.children = data as any[]
                    
                    // Reset selectedChild first
                    this.selectedChild = null
                    
                    // First, try to restore selected child from localStorage
                    // IMPORTANT: Only restore if the child exists in the current user's children list
                    if (process.client && this.children.length > 0) {
                        const savedChildId = localStorage.getItem('selectedChildId')
                        if (savedChildId) {
                            const savedChild = this.children.find(c => c.id === savedChildId)
                            if (savedChild) {
                                // Valid child found, restore it
                                this.selectedChild = savedChild
                                return // Exit early if we found and restored the saved child
                            } else {
                                // Saved child ID doesn't belong to current user, clear it
                                localStorage.removeItem('selectedChildId')
                            }
                        }
                    }
                    
                    // Auto-select first child only if no saved selection and no current selection
                    if (this.children.length > 0 && !this.selectedChild) {
                        this.selectedChild = this.children[0]
                        if (process.client) {
                            localStorage.setItem('selectedChildId', this.children[0].id)
                        }
                    } else if (this.children.length === 0) {
                        // No children, clear selection
                        this.selectedChild = null
                        if (process.client) {
                            localStorage.removeItem('selectedChildId')
                        }
                    }
                } else {
                    // No data returned, clear state
                    this.children = []
                    this.selectedChild = null
                    if (process.client) {
                        localStorage.removeItem('selectedChildId')
                    }
                }
            } catch (e) {
                console.error('Failed to fetch children', e)
                // On error, clear state
                this.children = []
                this.selectedChild = null
                if (process.client) {
                    localStorage.removeItem('selectedChildId')
                }
            }
        },

        selectChild(childId: string) {
            const child = this.children.find(c => c.id === childId)
            if (child) {
                this.selectedChild = child
                if (process.client) {
                    localStorage.setItem('selectedChildId', childId)
                }
            }
        },

        async addChild(childData) {
            const authStore = useAuthStore()
            const apiBase = useApiUrl()

            try {
                const data = await $fetch(`${apiBase}/api/children`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${authStore.token}`,
                        'Content-Type': 'application/json'
                    },
                    body: childData
                })

                await this.fetchChildren()
                return { success: true, data }
            } catch (e) {
                console.error('Failed to add child:', e)
                return { success: false, error: e.data?.error || e.message || 'Failed to add child' }
            }
        },

        async updateChild(childId: string, childData: any) {
            const authStore = useAuthStore()
            const apiBase = useApiUrl()

            try {
                await $fetch(`${apiBase}/api/children/${childId}`, {
                    method: 'PUT',
                    headers: {
                        Authorization: `Bearer ${authStore.token}`,
                        'Content-Type': 'application/json'
                    },
                    body: childData
                })

                await this.fetchChildren()
                return { success: true }
            } catch (e) {
                console.error('Failed to update child', e)
                return { success: false, error: e }
            }
        },

        async deleteChild(childId: string) {
            const authStore = useAuthStore()
            const apiBase = useApiUrl()

            try {
                await $fetch(`${apiBase}/api/children/${childId}`, {
                    method: 'DELETE',
                    headers: {
                        Authorization: `Bearer ${authStore.token}`
                    }
                })

                // If the deleted child was selected, clear the selection
                if (this.selectedChild?.id === childId) {
                    this.selectedChild = null
                    if (process.client) {
                        localStorage.removeItem('selectedChildId')
                    }
                }

                await this.fetchChildren()
                return { success: true }
            } catch (e) {
                console.error('Failed to delete child', e)
                return { success: false, error: e }
            }
        },

        initialize() {
            // This method is called before fetchChildren, so we just ensure localStorage is checked
            // The actual restoration happens in fetchChildren after children are loaded
            // This method is kept for backward compatibility but the main logic is in fetchChildren
        },

        clearState() {
            this.children = []
            this.selectedChild = null
            if (process.client) {
                localStorage.removeItem('selectedChildId')
            }
        }
    },
})
