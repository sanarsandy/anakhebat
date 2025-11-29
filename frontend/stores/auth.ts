import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
    const token = useCookie<string | null>('token', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: false
    })
    const user = useCookie<any | null>('user', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: false
    })

    const isAuthenticated = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    function setToken(newToken: string) {
        token.value = newToken
        if (process.client) {
            // Fallback: Manually set cookie to ensure persistence
            document.cookie = `token=${newToken}; path=/; max-age=${60 * 60 * 24 * 3}; SameSite=Lax`
        }
    }

    function setUser(newUser: any) {
        user.value = newUser
        if (process.client) {
            document.cookie = `user=${JSON.stringify(newUser)}; path=/; max-age=${60 * 60 * 24 * 3}; SameSite=Lax`
        }
    }

    function initialize() {
        if (process.client && !token.value) {
            const cookieToken = document.cookie.split('; ').find(row => row.startsWith('token='))?.split('=')[1]
            const cookieUser = document.cookie.split('; ').find(row => row.startsWith('user='))?.split('=')[1]

            if (cookieToken) {
                token.value = cookieToken
            }
            if (cookieUser) {
                try {
                    user.value = JSON.parse(decodeURIComponent(cookieUser))
                } catch (e) {
                    console.error('Failed to parse user cookie', e)
                }
            }
        }
    }

    function logout() {
        // Clear all stores before logout
        const childStore = useChildStore()
        const measurementStore = useMeasurementStore()
        const milestoneStore = useMilestoneStore()
        
        // Clear all store states
        childStore.clearState()
        measurementStore.clearMeasurements()
        milestoneStore.clearState()
        
        // Clear auth state
        token.value = null
        user.value = null
        if (process.client) {
            document.cookie = 'token=; path=/; max-age=0'
            document.cookie = 'user=; path=/; max-age=0'
            // Clear localStorage for selected child
            localStorage.removeItem('selectedChildId')
        }
        navigateTo('/login')
    }

    return {
        token,
        user,
        isAuthenticated,
        isAdmin,
        setToken,
        setUser,
        initialize,
        logout
    }
})
