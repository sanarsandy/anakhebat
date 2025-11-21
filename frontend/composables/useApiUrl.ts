export const useApiUrl = () => {
    const config = useRuntimeConfig()
    if (process.server) {
        return config.apiInternal
    }
    return config.public.apiBase
}
