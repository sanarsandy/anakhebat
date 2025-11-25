<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Profil Saya</h1>
      <p class="text-gray-600 mt-2">Kelola informasi profil dan pengaturan akun Anda</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !profile" class="bg-white rounded-xl shadow-sm p-8 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-4 text-gray-600">Memuat profil...</p>
    </div>

    <!-- Profile Content -->
    <div v-else class="space-y-6">
      <!-- Avatar & Basic Info Card -->
      <div class="bg-white rounded-xl shadow-sm p-8">
        <div class="flex items-center space-x-6 mb-6">
          <!-- Avatar -->
          <div class="w-20 h-20 rounded-full bg-indigo-100 flex items-center justify-center">
            <span class="text-3xl font-bold text-indigo-600">
              {{ getInitials(profile?.user?.full_name || '') }}
            </span>
          </div>
          
          <!-- User Info -->
          <div class="flex-1">
            <h2 class="text-2xl font-bold text-gray-900">{{ profile?.user?.full_name || 'User' }}</h2>
            <div class="flex items-center space-x-4 mt-2">
              <!-- Auth Provider Badge -->
              <span 
                class="px-3 py-1 rounded-full text-xs font-medium"
                :class="getAuthProviderBadgeClass(profile?.user?.auth_provider)"
              >
                {{ getAuthProviderLabel(profile?.user?.auth_provider) }}
              </span>
              
              <!-- Phone Verified Badge -->
              <span 
                v-if="profile?.user?.phone_verified"
                class="px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800"
              >
                ✓ Nomor Terverifikasi
              </span>
              <span 
                v-else-if="profile?.user?.phone_number"
                class="px-3 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
              >
                ⚠ Belum Terverifikasi
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Statistics Cards -->
      <div v-if="profile?.stats" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white rounded-xl shadow-sm p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600">Jumlah Anak</p>
              <p class="text-3xl font-bold text-gray-900 mt-2">{{ profile.stats.children_count || 0 }}</p>
            </div>
            <div class="text-4xl">👶</div>
          </div>
        </div>
        
        <div class="bg-white rounded-xl shadow-sm p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600">Data Pengukuran</p>
              <p class="text-3xl font-bold text-gray-900 mt-2">{{ profile.stats.measurements_count || 0 }}</p>
            </div>
            <div class="text-4xl">📊</div>
          </div>
        </div>
        
        <div class="bg-white rounded-xl shadow-sm p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600">Milestone</p>
              <p class="text-3xl font-bold text-gray-900 mt-2">{{ profile.stats.milestones_count || 0 }}</p>
            </div>
            <div class="text-4xl">✅</div>
          </div>
        </div>
      </div>

      <!-- Edit Profile Form -->
      <div class="bg-white rounded-xl shadow-sm p-8">
        <h3 class="text-xl font-bold text-gray-900 mb-6">Informasi Pribadi</h3>
        
        <form @submit.prevent="handleUpdateProfile" class="space-y-6">
          <!-- Full Name -->
          <div>
            <label for="full-name" class="block text-sm font-medium text-gray-700 mb-2">
              Nama Lengkap *
            </label>
            <input
              id="full-name"
              v-model="form.fullName"
              type="text"
              required
              minlength="3"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="Masukkan nama lengkap"
            />
          </div>

          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="email@example.com"
            />
            <p class="mt-1 text-xs text-gray-500">Email opsional, digunakan untuk notifikasi</p>
          </div>

          <!-- Phone Number -->
          <div>
            <label for="phone-number" class="block text-sm font-medium text-gray-700 mb-2">
              Nomor WhatsApp
            </label>
            <div class="flex items-center space-x-3">
              <input
                id="phone-number"
                v-model="form.phoneNumber"
                type="tel"
                class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="+6281234567890"
                :disabled="phoneVerificationStep !== 'none'"
              />
              
              <!-- Verified Badge -->
              <span 
                v-if="profile?.user?.phone_verified && profile?.user?.phone_number"
                class="px-3 py-3 rounded-lg bg-green-100 text-green-800 text-sm font-medium whitespace-nowrap"
              >
                ✓ Terverifikasi
              </span>
              
              <!-- Verify Button -->
              <button
                v-else-if="profile?.user?.phone_number && !profile?.user?.phone_verified"
                type="button"
                @click="startPhoneVerification"
                class="px-4 py-3 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition whitespace-nowrap"
              >
                Verifikasi
              </button>
              
              <!-- Add Phone Button -->
              <button
                v-else-if="!profile?.user?.phone_number"
                type="button"
                @click="startPhoneVerification"
                class="px-4 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition whitespace-nowrap"
              >
                Tambah Nomor
              </button>
            </div>
            <p class="mt-1 text-xs text-gray-500">Format: +62xxxxxxxxxxx</p>
          </div>

          <!-- Phone Verification Flow -->
          <div v-if="phoneVerificationStep === 'request'" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p class="text-sm text-blue-800 mb-3">
              Kode OTP akan dikirim ke nomor <strong>{{ form.phoneNumber }}</strong>
            </p>
            <button
              type="button"
              @click="requestPhoneOTP"
              :disabled="phoneOTPLoading || !form.phoneNumber"
              class="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
            >
              <span v-if="phoneOTPLoading">Mengirim OTP...</span>
              <span v-else>Kirim OTP via WhatsApp</span>
            </button>
          </div>

          <div v-if="phoneVerificationStep === 'verify'" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p class="text-sm text-blue-800 mb-3">
              Masukkan kode OTP yang dikirim ke <strong>{{ form.phoneNumber }}</strong>
            </p>
            <div class="space-y-3">
              <input
                v-model="phoneOTPCode"
                type="text"
                maxlength="6"
                placeholder="123456"
                class="w-full px-4 py-3 border border-blue-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-center text-2xl tracking-widest"
                @input="handleOTPInput"
              />
              <div class="flex gap-2">
                <button
                  type="button"
                  @click="confirmPhoneOTP"
                  :disabled="phoneOTPLoading || phoneOTPCode.length !== 6"
                  class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
                >
                  <span v-if="phoneOTPLoading">Memverifikasi...</span>
                  <span v-else>Verifikasi</span>
                </button>
                <button
                  type="button"
                  @click="cancelPhoneVerification"
                  class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                >
                  Batal
                </button>
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
            <p class="text-sm text-red-800">{{ error }}</p>
          </div>

          <!-- Success Message -->
          <div v-if="success" class="bg-green-50 border border-green-200 rounded-lg p-4">
            <p class="text-sm text-green-800">{{ success }}</p>
          </div>

          <!-- Submit Button -->
          <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
            <button
              type="button"
              @click="resetForm"
              class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loading || !hasChanges"
              class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading">Menyimpan...</span>
              <span v-else>Simpan Perubahan</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()
const config = useRuntimeConfig()
const router = useRouter()

const profile = ref<any>(null)
const loading = ref(false)
const error = ref('')
const success = ref('')

// Form data
const form = ref({
  fullName: '',
  email: '',
  phoneNumber: ''
})

// Phone verification
const phoneVerificationStep = ref<'none' | 'request' | 'verify'>('none')
const phoneOTPLoading = ref(false)
const phoneOTPCode = ref('')

// Computed
const hasChanges = computed(() => {
  if (!profile.value?.user) return false
  return (
    form.value.fullName !== profile.value.user.full_name ||
    form.value.email !== (profile.value.user.email || '') ||
    form.value.phoneNumber !== (profile.value.user.phone_number || '')
  )
})

// Methods
const getInitials = (name: string) => {
  if (!name) return 'U'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

const getAuthProviderLabel = (provider: string) => {
  if (!provider || typeof provider !== 'string') return 'Unknown'
  const labels: Record<string, string> = {
    'phone': 'WhatsApp',
    'google': 'Google',
    'phone_google': 'WhatsApp & Google',
    'email': 'Email'
  }
  return labels[provider] || 'Unknown'
}

const getAuthProviderBadgeClass = (provider: string) => {
  if (!provider || typeof provider !== 'string') return 'bg-gray-100 text-gray-800'
  const classes: Record<string, string> = {
    'phone': 'bg-green-100 text-green-800',
    'google': 'bg-blue-100 text-blue-800',
    'phone_google': 'bg-purple-100 text-purple-800',
    'email': 'bg-gray-100 text-gray-800'
  }
  return classes[provider] || 'bg-gray-100 text-gray-800'
}

const loadProfile = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const token = authStore.token
    if (!token) {
      error.value = 'Anda harus login terlebih dahulu'
      router.push('/login')
      return
    }
    
    const data = await $fetch(`${config.public.apiBase}/api/user/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    profile.value = data
    form.value = {
      fullName: data.user?.full_name || '',
      email: data.user?.email || '',
      phoneNumber: data.user?.phone_number || ''
    }
  } catch (e: any) {
    console.error('Load profile error:', e)
    const err: any = e || {}
    error.value = (err.data && err.data.error) || err.message || 'Gagal memuat profil'
  } finally {
    loading.value = false
  }
}

const handleUpdateProfile = async () => {
  loading.value = true
  error.value = ''
  success.value = ''
  
  try {
    const updateData: any = {}
    if (form.value.fullName !== profile.value?.user?.full_name) {
      updateData.full_name = form.value.fullName
    }
    if (form.value.email !== (profile.value?.user?.email || '')) {
      updateData.email = form.value.email || null
    }
    
    if (Object.keys(updateData).length === 0) {
      error.value = 'Tidak ada perubahan yang disimpan'
      loading.value = false
      return
    }
    
    const token = authStore.token
    if (!token) {
      error.value = 'Anda harus login terlebih dahulu'
      return
    }
    
    const data = await $fetch(`${config.public.apiBase}/api/user/profile`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: updateData
    })
    
    success.value = 'Profil berhasil diupdate'
    
    // Update auth store
    if (data.user) {
      authStore.setUser(data.user)
    }
    
    // Reload profile
    await loadProfile()
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      success.value = ''
    }, 3000)
  } catch (e: any) {
    console.error('Update profile error:', e)
    const err: any = e || {}
    error.value = (err.data && err.data.error) || err.message || 'Gagal mengupdate profil'
  } finally {
    loading.value = false
  }
}

const startPhoneVerification = () => {
  if (!form.value.phoneNumber) {
    error.value = 'Masukkan nomor WhatsApp terlebih dahulu'
    return
  }
  phoneVerificationStep.value = 'request'
  error.value = ''
  success.value = ''
}

const requestPhoneOTP = async () => {
  phoneOTPLoading.value = true
  error.value = ''
  
  try {
    const token = authStore.token
    if (!token) {
      error.value = 'Anda harus login terlebih dahulu'
      return
    }
    
    const data = await $fetch(`${config.public.apiBase}/api/user/verify-phone`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        phone_number: form.value.phoneNumber
      }
    })
    
    if (data.success) {
      phoneVerificationStep.value = 'verify'
      success.value = 'OTP telah dikirim ke WhatsApp Anda'
      setTimeout(() => {
        success.value = ''
      }, 3000)
    } else {
      error.value = data.error || 'Gagal mengirim OTP'
    }
  } catch (e: any) {
    console.error('Request OTP error:', e)
    const err: any = e || {}
    error.value = (err.data && err.data.error) || err.message || 'Gagal mengirim OTP'
  } finally {
    phoneOTPLoading.value = false
  }
}

const confirmPhoneOTP = async () => {
  phoneOTPLoading.value = true
  error.value = ''
  
  try {
    const token = authStore.token
    if (!token) {
      error.value = 'Anda harus login terlebih dahulu'
      return
    }
    
    const data = await $fetch(`${config.public.apiBase}/api/user/verify-phone/confirm`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        phone_number: form.value.phoneNumber,
        otp: phoneOTPCode.value
      }
    })
    
    if (data.success || data.user) {
      success.value = 'Nomor WhatsApp berhasil diverifikasi'
      phoneVerificationStep.value = 'none'
      phoneOTPCode.value = ''
      
      // Update auth store
      if (data.user) {
        authStore.setUser(data.user)
      }
      
      // Reload profile
      await loadProfile()
      
      setTimeout(() => {
        success.value = ''
      }, 3000)
    } else {
      error.value = data.error || 'OTP tidak valid'
    }
  } catch (e: any) {
    console.error('Confirm OTP error:', e)
    const err: any = e || {}
    error.value = (err.data && err.data.error) || err.message || 'Gagal memverifikasi OTP'
  } finally {
    phoneOTPLoading.value = false
  }
}

const cancelPhoneVerification = () => {
  phoneVerificationStep.value = 'none'
  phoneOTPCode.value = ''
  error.value = ''
  success.value = ''
}

const handleOTPInput = (e: any) => {
  phoneOTPCode.value = e.target.value.replace(/\D/g, '').slice(0, 6)
}

const resetForm = () => {
  if (profile.value?.user) {
    form.value = {
      fullName: profile.value.user.full_name || '',
      email: profile.value.user.email || '',
      phoneNumber: profile.value.user.phone_number || ''
    }
  }
  error.value = ''
  success.value = ''
  phoneVerificationStep.value = 'none'
  phoneOTPCode.value = ''
}

// Lifecycle
onMounted(() => {
  loadProfile()
})
</script>

