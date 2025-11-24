<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Sign in to your account
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <input type="hidden" name="remember" value="true" />
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email-address" class="sr-only">Email address</label>
            <input id="email-address" name="email" type="email" autocomplete="email" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm" placeholder="Email address" v-model="email" />
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <input id="password" name="password" type="password" autocomplete="current-password" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm" placeholder="Password" v-model="password" />
          </div>
        </div>

        <div v-if="error" class="text-red-500 text-sm text-center">
          {{ error }}
        </div>

        <div>
          <button type="submit" :disabled="loading" class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="loading">Signing in...</span>
            <span v-else>Sign in</span>
          </button>
        </div>

        <!-- Divider -->
        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-300"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-gray-50 text-gray-500">Or continue with</span>
          </div>
        </div>

        <!-- Google Sign-In Button -->
        <div>
          <button 
            type="button"
            @click="handleGoogleLogin"
            :disabled="loading"
            class="w-full flex items-center justify-center gap-3 py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            <span>Sign in with Google</span>
          </button>
        </div>

        <div class="text-center">
          <NuxtLink to="/register" class="font-medium text-indigo-600 hover:text-indigo-500">
            Don't have an account? Register
          </NuxtLink>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: false
})

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const authStore = useAuthStore()
const config = useRuntimeConfig()
const router = useRouter()

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  
  try {
    const data = await $fetch(`${config.public.apiBase}/api/auth/login`, {
      method: 'POST',
      body: {
        email: email.value,
        password: password.value
      }
    })

    console.log('Login success:', data)

    if (data.token) {
      // Clear all stores and localStorage before setting new token
      const childStore = useChildStore()
      const measurementStore = useMeasurementStore()
      const milestoneStore = useMilestoneStore()
      
      // Clear all states
      childStore.clearState()
      measurementStore.clearMeasurements()
      milestoneStore.clearState()
      
      // Set new auth data
      authStore.setToken(data.token)
      authStore.setUser(data.user)
      
      // Navigate to dashboard - fetchChildren will be called there
      router.push('/dashboard')
    }
  } catch (e) {
    console.error('Login error:', e)
    error.value = e.data?.error || e.message || 'Login failed'
  } finally {
    loading.value = false
  }
}

const handleGoogleLogin = async () => {
  error.value = ''
  loading.value = true
  
  try {
    // Get Google OAuth URL from backend
    const response = await $fetch(`${config.public.apiBase}/api/auth/google`, {
      method: 'GET'
    })

    if (response.auth_url) {
      // Redirect to Google OAuth
      window.location.href = response.auth_url
    } else {
      error.value = 'Failed to get Google OAuth URL'
      loading.value = false
    }
  } catch (e) {
    console.error('Google login error:', e)
    error.value = e.data?.error || e.message || 'Google login failed'
    loading.value = false
  }
}
</script>
