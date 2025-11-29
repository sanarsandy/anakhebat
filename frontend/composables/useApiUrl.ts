export const useApiUrl = () => {
    const config = useRuntimeConfig()
    if (process.server) {
        return config.apiInternal || 'http://localhost:8080'
    }
    return config.public.apiBase || 'http://localhost:8080'
}
