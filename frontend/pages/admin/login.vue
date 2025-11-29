<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <!-- Logo -->
        <div class="flex justify-center mb-6">
          <img 
            v-if="logoUrl"
            :src="logoUrl" 
            alt="Jurnal Si Kecil" 
            class="h-16 w-auto object-contain"
          />
        </div>
        <h2 class="mt-6 text-center text-3xl font-bold text-jurnal-charcoal-800">
          Admin Login
        </h2>
        <p class="mt-2 text-center text-sm text-jurnal-charcoal-600">
          Masuk ke Admin Panel Jurnal Si Kecil
        </p>
      </div>

      <!-- Login Form -->
      <div class="bg-white p-8 rounded-soft-lg shadow-md border border-gray-200">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input 
              id="email" 
              v-model="email" 
              type="email" 
              placeholder="admin@jurnalsikecil.com" 
              required 
              autocomplete="email"
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-soft focus:outline-none focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500 focus:z-10 sm:text-sm" 
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input 
              id="password" 
              v-model="password" 
              type="password" 
              placeholder="Masukkan password" 
              required 
              autocomplete="current-password"
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-soft focus:outline-none focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500 focus:z-10 sm:text-sm" 
            />
          </div>

          <!-- Error Message -->
          <div v-if="error" class="text-red-500 text-sm">
            <div class="bg-red-50 border border-red-200 rounded-soft p-3">
              <p class="font-medium">{{ error }}</p>
            </div>
          </div>

          <button 
            type="submit"
            :disabled="loading || !email || !password"
            class="w-full flex justify-center items-center gap-2 py-2 px-4 border border-transparent text-sm font-medium rounded-soft text-white bg-jurnal-teal-600 hover:bg-jurnal-teal-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-jurnal-teal-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memproses...
            </span>
            <span v-else>Masuk</span>
          </button>
        </form>

        <!-- Back to User Login -->
        <div class="mt-6 text-center">
          <NuxtLink 
            to="/login" 
            class="text-sm text-jurnal-teal-600 hover:text-jurnal-teal-500 font-medium"
          >
            ← Kembali ke Login User
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import logoImage from '~/assets/images/logo.png'

definePageMeta({
  layout: false
})

const logoUrl = ref<string | null>(null)
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const authStore = useAuthStore()
const config = useRuntimeConfig()
const router = useRouter()

onMounted(() => {
  // Set logo URL only on client to avoid hydration mismatch
  logoUrl.value = logoImage
})

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    const apiBase = config.public.apiBase || 'http://localhost:8080'
    const data = await $fetch(`${apiBase}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: {
        email: email.value.trim(),
        password: password.value
      }
    })

    if (data && data.token && data.user) {
      // Check if user is admin
      if (data.user.role !== 'admin') {
        error.value = 'Akses ditolak. Hanya admin yang dapat mengakses halaman ini.'
        loading.value = false
        password.value = ''
        return
      }

      // Clear all stores before setting new token
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

      // Navigate to admin dashboard
      router.push('/admin/dashboard')
    } else {
      error.value = 'Login gagal. Silakan coba lagi.'
      password.value = ''
      loading.value = false
    }
  } catch (e: any) {
    console.error('Admin login error:', e)
    
    // Extract error message
    let errorMsg = 'Email atau password salah'
    
    if (e.data) {
      if (typeof e.data === 'string') {
        errorMsg = e.data
      } else if (e.data.error) {
        errorMsg = e.data.error
      } else if (e.data.message) {
        errorMsg = e.data.message
      }
    } else if (e.message) {
      errorMsg = e.message
    } else if (e.statusMessage) {
      errorMsg = e.statusMessage
    }
    
    // Handle specific error cases
    if (errorMsg.includes('User belum terdaftar') || errorMsg.includes('not found')) {
      errorMsg = 'Email tidak terdaftar'
    } else if (errorMsg.includes('Invalid credentials') || errorMsg.includes('Invalid')) {
      errorMsg = 'Email atau password salah'
    } else if (errorMsg.includes('Network') || errorMsg.includes('fetch')) {
      errorMsg = 'Tidak dapat terhubung ke server. Pastikan backend API berjalan.'
    }
    
    error.value = errorMsg
    password.value = '' // Clear password on error
    loading.value = false
  }
}
</script>

