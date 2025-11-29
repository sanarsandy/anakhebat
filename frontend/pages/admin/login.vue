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
          Masuk ke Admin Panel dengan OTP WhatsApp
        </p>
      </div>

      <!-- Login Form -->
      <div class="bg-white p-8 rounded-soft-lg shadow-md border border-gray-200">
        <!-- Step 1: Phone Number Input -->
        <form v-if="step === 'phone'" @submit.prevent="handleRequestOTP" class="space-y-6">
          <div>
            <label for="phone" class="block text-sm font-medium text-gray-700 mb-1">
              Nomor WhatsApp Admin
            </label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 text-sm">+62</span>
              <input 
                id="phone" 
                v-model="phoneNumber" 
                type="tel" 
                placeholder="81234567890" 
                required 
                autocomplete="tel"
                class="appearance-none relative block w-full pl-12 pr-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-soft focus:outline-none focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500 focus:z-10 sm:text-sm" 
              />
            </div>
            <p class="mt-1 text-xs text-gray-500">Masukkan nomor WhatsApp yang terdaftar sebagai admin</p>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="text-red-500 text-sm">
            <div class="bg-red-50 border border-red-200 rounded-soft p-3">
              <p class="font-medium">{{ error }}</p>
            </div>
          </div>

          <button 
            type="submit"
            :disabled="loading || !phoneNumber"
            class="w-full flex justify-center items-center gap-2 py-2 px-4 border border-transparent text-sm font-medium rounded-soft text-white bg-jurnal-teal-600 hover:bg-jurnal-teal-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-jurnal-teal-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Mengirim OTP...
            </span>
            <span v-else class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
              </svg>
              Kirim OTP ke WhatsApp
            </span>
          </button>
        </form>

        <!-- Step 2: OTP Input -->
        <form v-else-if="step === 'otp'" @submit.prevent="handleVerifyOTP" class="space-y-6">
          <div class="text-center mb-4">
            <p class="text-sm text-gray-600">
              Kode OTP telah dikirim ke WhatsApp
            </p>
            <p class="font-medium text-jurnal-charcoal-800">+62{{ phoneNumber }}</p>
          </div>

          <div>
            <label for="otp" class="block text-sm font-medium text-gray-700 mb-1">
              Masukkan Kode OTP
            </label>
            <input 
              id="otp" 
              v-model="otpCode" 
              type="text" 
              inputmode="numeric"
              pattern="[0-9]*"
              maxlength="6"
              placeholder="000000" 
              required 
              autocomplete="one-time-code"
              class="appearance-none relative block w-full px-3 py-3 border border-gray-300 placeholder-gray-400 text-gray-900 text-center text-2xl tracking-[0.5em] font-mono rounded-soft focus:outline-none focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500 focus:z-10" 
            />
          </div>

          <!-- Countdown Timer -->
          <div class="text-center text-sm text-gray-500">
            <span v-if="countdown > 0">
              Kode berlaku {{ Math.floor(countdown / 60) }}:{{ String(countdown % 60).padStart(2, '0') }}
            </span>
            <span v-else class="text-red-500">
              Kode sudah kadaluarsa
            </span>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="text-red-500 text-sm">
            <div class="bg-red-50 border border-red-200 rounded-soft p-3">
              <p class="font-medium">{{ error }}</p>
            </div>
          </div>

          <button 
            type="submit"
            :disabled="loading || !otpCode || otpCode.length < 6"
            class="w-full flex justify-center items-center gap-2 py-2 px-4 border border-transparent text-sm font-medium rounded-soft text-white bg-jurnal-teal-600 hover:bg-jurnal-teal-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-jurnal-teal-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memverifikasi...
            </span>
            <span v-else>Verifikasi & Masuk</span>
          </button>

          <!-- Resend OTP -->
          <div class="text-center">
            <button 
              type="button"
              @click="handleResendOTP"
              :disabled="resendCooldown > 0 || loading"
              class="text-sm text-jurnal-teal-600 hover:text-jurnal-teal-500 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="resendCooldown > 0">
                Kirim ulang dalam {{ resendCooldown }} detik
              </span>
              <span v-else>
                Kirim Ulang OTP
              </span>
            </button>
          </div>

          <!-- Back Button -->
          <div class="text-center">
            <button 
              type="button"
              @click="goBackToPhone"
              class="text-sm text-gray-500 hover:text-gray-700"
            >
              ← Ganti nomor WhatsApp
            </button>
          </div>
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
const phoneNumber = ref('')
const otpCode = ref('')
const error = ref('')
const loading = ref(false)
const step = ref<'phone' | 'otp'>('phone')
const countdown = ref(300) // 5 minutes
const resendCooldown = ref(0)

const authStore = useAuthStore()
const config = useRuntimeConfig()
const router = useRouter()

let countdownInterval: ReturnType<typeof setInterval> | null = null
let resendInterval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  logoUrl.value = logoImage
})

onUnmounted(() => {
  if (countdownInterval) clearInterval(countdownInterval)
  if (resendInterval) clearInterval(resendInterval)
})

const startCountdown = () => {
  countdown.value = 300
  if (countdownInterval) clearInterval(countdownInterval)
  countdownInterval = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--
    } else {
      if (countdownInterval) clearInterval(countdownInterval)
    }
  }, 1000)
}

const startResendCooldown = () => {
  resendCooldown.value = 60
  if (resendInterval) clearInterval(resendInterval)
  resendInterval = setInterval(() => {
    if (resendCooldown.value > 0) {
      resendCooldown.value--
    } else {
      if (resendInterval) clearInterval(resendInterval)
    }
  }, 1000)
}

const formatPhoneNumber = (phone: string): string => {
  // Remove all non-digits
  let cleaned = phone.replace(/\D/g, '')
  
  // If starts with 0, replace with 62
  if (cleaned.startsWith('0')) {
    cleaned = '62' + cleaned.substring(1)
  }
  // If doesn't start with 62, add it
  if (!cleaned.startsWith('62')) {
    cleaned = '62' + cleaned
  }
  
  return cleaned
}

const handleRequestOTP = async () => {
  error.value = ''
  loading.value = true

  try {
    const apiBase = config.public.apiBase || 'http://localhost:8080'
    const formattedPhone = formatPhoneNumber(phoneNumber.value)
    
    const data = await $fetch<{ success: boolean; message?: string; error?: string }>(`${apiBase}/api/auth/admin/request-otp`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: {
        phone_number: formattedPhone
      }
    })

    if (data && data.success) {
      step.value = 'otp'
      startCountdown()
      startResendCooldown()
    } else {
      error.value = data?.error || 'Gagal mengirim OTP'
    }
  } catch (e: any) {
    console.error('Request OTP error:', e)
    
    let errorMsg = 'Gagal mengirim OTP'
    if (e.data?.error) {
      errorMsg = e.data.error
    } else if (e.message) {
      errorMsg = e.message
    }
    
    error.value = errorMsg
  } finally {
    loading.value = false
  }
}

const handleVerifyOTP = async () => {
  error.value = ''
  loading.value = true

  try {
    const apiBase = config.public.apiBase || 'http://localhost:8080'
    const formattedPhone = formatPhoneNumber(phoneNumber.value)
    
    const data = await $fetch<{ token?: string; user?: any; success?: boolean; error?: string }>(`${apiBase}/api/auth/admin/verify-otp`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: {
        phone_number: formattedPhone,
        otp: otpCode.value
      }
    })

    if (data && data.token && data.user) {
      // Check if user is admin
      if (data.user.role !== 'admin' && data.user.role !== 'super_admin') {
        error.value = 'Akses ditolak. Hanya admin yang dapat mengakses halaman ini.'
        loading.value = false
        otpCode.value = ''
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
      error.value = data?.error || 'Verifikasi OTP gagal'
      otpCode.value = ''
    }
  } catch (e: any) {
    console.error('Verify OTP error:', e)
    
    let errorMsg = 'OTP tidak valid'
    if (e.data?.error) {
      errorMsg = e.data.error
    } else if (e.message) {
      errorMsg = e.message
    }
    
    error.value = errorMsg
    otpCode.value = ''
  } finally {
    loading.value = false
  }
}

const handleResendOTP = async () => {
  error.value = ''
  loading.value = true

  try {
    const apiBase = config.public.apiBase || 'http://localhost:8080'
    const formattedPhone = formatPhoneNumber(phoneNumber.value)
    
    const data = await $fetch<{ success: boolean; message?: string; error?: string }>(`${apiBase}/api/auth/admin/request-otp`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: {
        phone_number: formattedPhone
      }
    })

    if (data && data.success) {
      startCountdown()
      startResendCooldown()
      otpCode.value = ''
    } else {
      error.value = data?.error || 'Gagal mengirim ulang OTP'
    }
  } catch (e: any) {
    console.error('Resend OTP error:', e)
    error.value = e.data?.error || 'Gagal mengirim ulang OTP'
  } finally {
    loading.value = false
  }
}

const goBackToPhone = () => {
  step.value = 'phone'
  otpCode.value = ''
  error.value = ''
  if (countdownInterval) clearInterval(countdownInterval)
}
</script>
