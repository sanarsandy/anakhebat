export default defineNuxtRouteMiddleware((to, from) => {
    const authStore = useAuthStore()

    // Skip server-side check to avoid premature redirects if cookie isn't visible yet
    if (process.server) return

    const token = useCookie('token')

    console.log('Auth Middleware - Path:', to.path)
    console.log('Auth Middleware - Token (Cookie):', token.value)
    console.log('Auth Middleware - Token (Store):', authStore.token)

    if (!token.value && !authStore.token && to.path !== '/login' && to.path !== '/register') {
        console.log('Auth Middleware - Redirecting to login')
        return navigateTo('/login')
    }
})
