<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <div v-if="loading" class="space-y-4">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="text-gray-600">Completing sign in...</p>
      </div>
      <div v-else-if="error" class="space-y-4">
        <div class="text-red-500 text-5xl mb-4">⚠️</div>
        <h2 class="text-xl font-bold text-gray-900">Sign in failed</h2>
        <p class="text-gray-600">{{ error }}</p>
        <NuxtLink to="/login" class="inline-block mt-4 px-6 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700">
          Back to Login
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: false
})

const loading = ref(true)
const error = ref('')
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

onMounted(async () => {
  const token = route.query.token
  const userData = route.query.user
  const state = route.query.state

  if (!token) {
    error.value = 'No authentication token received'
    loading.value = false
    return
  }

  try {
    // Clear all stores before setting new token
    const childStore = useChildStore()
    const measurementStore = useMeasurementStore()
    const milestoneStore = useMilestoneStore()
    
    // Clear all states
    childStore.clearState()
    measurementStore.clearMeasurements()
    milestoneStore.clearState()
    
    // Set token
    authStore.setToken(token)
    
    // Set user data if provided
    if (userData) {
      try {
        const user = JSON.parse(decodeURIComponent(userData))
        authStore.setUser(user)
      } catch (e) {
        console.warn('Failed to parse user data:', e)
      }
    }
    
    // Navigate to dashboard
    router.push('/dashboard')
  } catch (e) {
    console.error('Callback error:', e)
    error.value = e.data?.error || e.message || 'Failed to complete sign in'
    loading.value = false
  }
})
</script>

