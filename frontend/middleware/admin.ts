export default defineNuxtRouteMiddleware(async (to, from) => {
    const authStore = useAuthStore()

    // Skip server-side check to avoid premature redirects
    if (process.server) return

    // Allow access to admin login page
    if (to.path === '/admin/login') {
        return
    }

    // Check if user is authenticated
    const token = useCookie('token')
    if (!token.value && !authStore.token) {
        return navigateTo('/admin/login')
    }

    // If user data is not loaded, try to fetch it
    if (!authStore.user && token.value) {
        try {
            const { apiBase } = useApiUrl()
            const userData = await $fetch(`${apiBase}/api/user/profile`, {
                headers: {
                    Authorization: `Bearer ${token.value}`
                }
            })
            if (userData && userData.user) {
                authStore.setUser(userData.user)
            }
        } catch (e) {
            // If fetch fails, redirect to admin login
            return navigateTo('/admin/login')
        }
    }

    // Check if user is admin
    if (!authStore.isAdmin) {
        // If not admin, redirect to regular dashboard
        return navigateTo('/dashboard')
    }
})

