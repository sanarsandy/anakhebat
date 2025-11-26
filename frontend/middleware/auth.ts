export default defineNuxtRouteMiddleware((to, from) => {
    const authStore = useAuthStore()

    // Skip server-side check to avoid premature redirects if cookie isn't visible yet
    if (process.server) return

    const token = useCookie('token')

    if (!token.value && !authStore.token && to.path !== '/login' && to.path !== '/register') {
        return navigateTo('/login')
    }
})
