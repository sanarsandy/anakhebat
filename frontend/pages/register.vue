<template>
  <div class="min-h-screen flex items-center justify-center bg-white py-12 px-4 sm:px-6 lg:px-8">
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
          Daftar Akun Baru
        </h2>
        <p class="mt-2 text-center text-sm text-jurnal-charcoal-600">
          Daftar menggunakan nomor WhatsApp Anda
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="rounded-soft shadow-sm -space-y-px">
          <div>
            <label for="full-name" class="sr-only">Nama Lengkap</label>
            <input 
              id="full-name" 
              name="full_name" 
              type="text" 
              required 
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-jurnal-charcoal-300 placeholder-jurnal-charcoal-400 text-jurnal-charcoal-800 rounded-t-soft focus:outline-none focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500 focus:z-10 sm:text-sm" 
              placeholder="Nama Lengkap" 
              v-model="fullName" 
            />
          </div>
          <div>
            <label for="phone-number" class="sr-only">Nomor WhatsApp</label>
            <input 
              id="phone-number" 
              name="phone_number" 
              type="tel" 
              autocomplete="tel" 
              required 
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500 focus:z-10 sm:text-sm" 
              placeholder="Nomor WhatsApp (contoh: 081234567890 atau +6281234567890)" 
              v-model="phoneNumber" 
            />
          </div>
        </div>

        <div v-if="error" class="text-red-500 text-sm text-center bg-red-50 p-3 rounded-soft border border-red-200">
          {{ error }}
        </div>

        <div v-if="success" class="text-green-600 text-sm text-center bg-green-50 p-3 rounded-soft border border-green-200">
          {{ success }}
        </div>

        <div>
          <button 
            type="submit" 
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-soft text-white bg-jurnal-teal-600 hover:bg-jurnal-teal-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-jurnal-teal-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading">Mendaftarkan...</span>
            <span v-else>Daftar</span>
          </button>
        </div>
        <div class="text-center">
          <NuxtLink to="/login" class="font-medium text-jurnal-teal-600 hover:text-jurnal-teal-500">
            Sudah punya akun? Masuk di sini
          </NuxtLink>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import logoImage from '~/assets/images/logo.png'

definePageMeta({
  layout: false
})

const logoUrl = ref(null)
const fullName = ref('')

onMounted(() => {
  // Set logo URL only on client to avoid hydration mismatch
  logoUrl.value = logoImage
})
const phoneNumber = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)
const config = useRuntimeConfig()
const router = useRouter()

const handleRegister = async () => {
  error.value = ''
  success.value = ''
  loading.value = true
  
  try {
    // Register user with phone number
    const registerData = await $fetch(`${config.public.apiBase}/api/auth/register`, {
      method: 'POST',
      body: {
        full_name: fullName.value,
        phone_number: phoneNumber.value
      }
    })

    // Show success message and redirect to login
    success.value = 'Pendaftaran berhasil! Silakan login dengan nomor WhatsApp Anda.'
    
    // Redirect to login page after 2 seconds
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (e) {
    console.error('Registration error:', e)
    const err: any = e || {}
    const errorMsg = (err.data && err.data.error) || err.message || 'Pendaftaran gagal'
    error.value = errorMsg
  } finally {
    loading.value = false
  }
}
</script>
