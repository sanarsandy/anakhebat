<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Masuk ke Akun Anda
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Pilih metode login yang Anda inginkan
        </p>
      </div>

      <!-- Login Method Selection -->
      <div class="mt-8 space-y-4">
        <!-- WhatsApp OTP Login -->
        <div v-if="loginMethod === 'whatsapp' || loginMethod === null" class="space-y-4">
          <div class="bg-white p-6 rounded-lg shadow-md border border-gray-200">
            <div class="flex items-center justify-center mb-4">
              <svg class="w-8 h-8 text-green-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
              </svg>
              <h3 class="ml-2 text-lg font-semibold text-gray-900">Login dengan WhatsApp</h3>
            </div>

            <!-- Step 1: Request OTP -->
            <div v-if="otpStep === 'request'" class="space-y-4">
              <div>
                <label for="phone-number" class="block text-sm font-medium text-gray-700 mb-1">
                  Nomor WhatsApp
                </label>
                <input 
                  id="phone-number" 
                  v-model="phoneNumber" 
                  type="tel" 
                  placeholder="+6281234567890" 
                  required 
                  class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm" 
                />
                <p class="mt-1 text-xs text-gray-500">Format: +62xxxxxxxxxxx</p>
              </div>

              <div v-if="error" class="text-red-500 text-sm text-center">
                {{ error }}
              </div>

              <button 
                @click="handleRequestOTP" 
                :disabled="loading || !phoneNumber"
                class="w-full flex justify-center items-center gap-2 py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span v-if="loading">Mengirim OTP...</span>
                <span v-else>Kirim OTP via WhatsApp</span>
              </button>
            </div>

            <!-- Step 2: Verify OTP -->
            <div v-if="otpStep === 'verify'" class="space-y-4">
              <div>
                <p class="text-sm text-gray-600 mb-2">
                  Kode OTP telah dikirim ke <strong>{{ phoneNumber }}</strong>
                </p>
                <label for="otp-code" class="block text-sm font-medium text-gray-700 mb-1">
                  Masukkan Kode OTP (6 digit)
                </label>
                <input 
                  id="otp-code" 
                  v-model="otpCode" 
                  type="text" 
                  maxlength="6" 
                  placeholder="123456" 
                  required 
                  class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm text-center text-2xl tracking-widest" 
                  @input="handleOTPInput"
                />
                <p class="mt-1 text-xs text-gray-500">
                  Kode berlaku selama 5 menit
                  <span v-if="otpExpiresIn > 0" class="ml-2">({{ formatTime(otpExpiresIn) }})</span>
                </p>
              </div>

              <div v-if="error" class="text-red-500 text-sm text-center">
                {{ error }}
              </div>

              <div class="flex gap-2">
                <button 
                  @click="handleVerifyOTP" 
                  :disabled="loading || otpCode.length !== 6"
                  class="flex-1 flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span v-if="loading">Memverifikasi...</span>
                  <span v-else>Verifikasi OTP</span>
                </button>
                <button 
                  @click="handleResendOTP" 
                  :disabled="loading || resendCooldown > 0"
                  class="px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span v-if="resendCooldown > 0">{{ resendCooldown }}s</span>
                  <span v-else>Kirim Ulang</span>
                </button>
              </div>

              <button 
                @click="otpStep = 'request'; otpCode = ''; error = ''" 
                class="w-full text-sm text-gray-600 hover:text-gray-800"
              >
                ← Kembali ke input nomor
              </button>
            </div>
          </div>
        </div>

        <!-- Divider -->
        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-300"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-gray-50 text-gray-500">atau</span>
          </div>
        </div>

        <!-- Google Sign-In Button -->
        <div class="space-y-2">
          <button 
            type="button"
            @click="handleGoogleLogin"
            :disabled="loading"
            class="w-full flex items-center justify-center gap-3 py-3 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            <span>Masuk dengan Google</span>
          </button>
          
          <!-- Error message for Google login -->
          <div v-if="error && loginMethod === 'google'" class="text-red-500 text-sm text-center bg-red-50 p-3 rounded-md border border-red-200">
            {{ error }}
          </div>
        </div>

        <!-- Register Link -->
        <div class="text-center mt-6">
          <p class="text-sm text-gray-600">
            Belum punya akun?
            <NuxtLink to="/register" class="font-medium text-indigo-600 hover:text-indigo-500">
              Daftar di sini
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false
})

const phoneNumber = ref('')
const otpCode = ref('')
const otpStep = ref('request') // 'request' or 'verify'
const error = ref('')
const loading = ref(false)
const loginMethod = ref(null) // 'whatsapp' or 'google' or null
const authStore = useAuthStore()
const config = useRuntimeConfig()
const router = useRouter()
const otpExpiresIn = ref(300) // 5 minutes in seconds
const resendCooldown = ref(0)

// OTP countdown timer
let otpTimer = null
let resendTimer = null

onMounted(() => {
  // Start OTP expiration countdown if in verify step
  if (otpStep.value === 'verify') {
    startOTPTimer()
  }
})

onUnmounted(() => {
  if (otpTimer) clearInterval(otpTimer)
  if (resendTimer) clearInterval(resendTimer)
})

const startOTPTimer = () => {
  if (otpTimer) clearInterval(otpTimer)
  otpTimer = setInterval(() => {
    if (otpExpiresIn.value > 0) {
      otpExpiresIn.value--
    } else {
      clearInterval(otpTimer)
      error.value = 'OTP sudah kadaluarsa. Silakan request OTP baru.'
    }
  }, 1000)
}

const startResendCooldown = (seconds = 60) => {
  resendCooldown.value = seconds
  if (resendTimer) clearInterval(resendTimer)
  resendTimer = setInterval(() => {
    if (resendCooldown.value > 0) {
      resendCooldown.value--
    } else {
      clearInterval(resendTimer)
    }
  }, 1000)
}

const formatTime = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const handleOTPInput = (e) => {
  // Only allow numbers
  otpCode.value = e.target.value.replace(/\D/g, '').slice(0, 6)
}

const handleRequestOTP = async () => {
  error.value = ''
  loading.value = true
  loginMethod.value = 'whatsapp'

  try {
    const response = await $fetch(`${config.public.apiBase}/api/auth/request-otp`, {
      method: 'POST',
      body: {
        phone_number: phoneNumber.value
      }
    })

    if (response.success) {
      otpStep.value = 'verify'
      otpExpiresIn.value = response.expires_in || 300
      startOTPTimer()
      startResendCooldown(60)
    } else {
      error.value = response.error || 'Gagal mengirim OTP'
    }
  } catch (e) {
    console.error('Request OTP error:', e)
    const err: any = e || {}
    const errorMsg = (err.data && err.data.error) || err.message || 'Gagal mengirim OTP'
    error.value = errorMsg
  } finally {
    loading.value = false
  }
}

const handleVerifyOTP = async () => {
  error.value = ''
  loading.value = true

  try {
    const data = await $fetch(`${config.public.apiBase}/api/auth/verify-otp`, {
      method: 'POST',
      body: {
        phone_number: phoneNumber.value,
        otp: otpCode.value
      }
    })

    if (data.token) {
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

      // Show welcome message for new users
      if (data.is_new_user) {
        // You can show a welcome message here
        console.log('New user registered:', data.user)
      }

      // Navigate to dashboard
      router.push('/dashboard')
    }
  } catch (e) {
    console.error('Verify OTP error:', e)
    const err: any = e || {}
    const errorMsg = (err.data && err.data.error) || err.message || 'OTP tidak valid'
    error.value = errorMsg
    otpCode.value = '' // Clear OTP on error
  } finally {
    loading.value = false
  }
}

const handleResendOTP = async () => {
  error.value = ''
  loading.value = true

  try {
    const response = await $fetch(`${config.public.apiBase}/api/auth/resend-otp`, {
      method: 'POST',
      body: {
        phone_number: phoneNumber.value
      }
    })

    if (response.success) {
      otpCode.value = ''
      otpExpiresIn.value = response.expires_in || 300
      startOTPTimer()
      startResendCooldown(60)
      error.value = '' // Clear any previous errors
    } else {
      error.value = response.error || 'Gagal mengirim ulang OTP'
    }
  } catch (e) {
    console.error('Resend OTP error:', e)
    const err: any = e || {}
    const errorMsg = (err.data && err.data.error) || err.message || 'Gagal mengirim ulang OTP'
    error.value = errorMsg
  } finally {
    loading.value = false
  }
}

const handleGoogleLogin = async () => {
  error.value = ''
  loading.value = true
  loginMethod.value = 'google'

  try {
    // Get Google OAuth URL from backend
    const response = await $fetch(`${config.public.apiBase}/api/auth/google`, {
      method: 'GET'
    })

    console.log('Google OAuth response:', response)

    if (response && response.auth_url) {
      // Redirect immediately to Google OAuth - don't wait
      console.log('Redirecting to Google OAuth:', response.auth_url)
      // Use window.location.replace to avoid back button issues
      window.location.replace(response.auth_url)
      // Don't set loading to false here - we're redirecting
      return
    } else {
      error.value = 'Gagal mendapatkan URL Google OAuth. Response tidak valid.'
      loading.value = false
    }
  } catch (e) {
    console.error('Google login error:', e)
    
    // Check for 503 status or "not configured" error
    const err: any = e || {}
    const status = err.status || err.statusCode || (err.response && err.response.status)
    const errorMessage = (err.data && err.data.error) || err.message || 'Gagal login dengan Google'
    
    if (status === 503 || (errorMessage && errorMessage.includes('not configured'))) {
      error.value = 'Google OAuth belum dikonfigurasi. Silakan hubungi administrator atau gunakan login dengan WhatsApp.'
    } else {
      error.value = `Gagal login dengan Google: ${errorMessage}`
    }
    
    loading.value = false
    loginMethod.value = null // Reset login method on error
  }
}
</script>
