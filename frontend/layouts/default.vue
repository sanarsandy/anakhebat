<template>
  <div class="min-h-screen bg-white">
    <!-- Desktop Sidebar -->
    <Sidebar class="hidden lg:block" />
    
    <!-- Main Content Area -->
    <div class="lg:ml-64">
      <!-- Header -->
      <header class="bg-white border-b border-jurnal-charcoal-200 sticky top-0 z-40">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex items-center justify-between h-16">
            <!-- Logo Section -->
            <div class="flex items-center gap-3">
              <NuxtLink to="/" class="flex items-center gap-3">
                <img 
                  v-if="logoUrl"
                  :src="logoUrl" 
                  alt="Jurnal Si Kecil" 
                  class="h-10 w-auto object-contain"
                />
                <h1 class="text-xl font-bold text-jurnal-teal-600 hidden sm:block">
                  Jurnal Si Kecil
                </h1>
              </NuxtLink>
            </div>
            
            <!-- Right Section -->
            <div class="flex items-center gap-3">
              <!-- Child Selector -->
              <ChildSelector v-if="authStore.isAuthenticated" />
              
              <!-- User Menu with Dropdown -->
              <div v-if="authStore.isAuthenticated" class="relative group">
                <button class="flex items-center gap-2 px-2 py-1.5 rounded-soft hover:bg-jurnal-charcoal-50 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-jurnal-teal-500 focus:ring-offset-2">
                  <div class="w-9 h-9 rounded-full bg-jurnal-teal-600 flex items-center justify-center shadow-sm">
                    <span class="text-xs font-semibold text-white">
                      {{ getInitials(authStore.user?.full_name || '') }}
                    </span>
                  </div>
                  <div class="hidden sm:block text-left">
                    <p class="text-sm font-medium text-jurnal-charcoal-800 leading-tight">
                      {{ authStore.user?.full_name }}
                    </p>
                    <p v-if="authStore.user?.email" class="text-xs text-jurnal-charcoal-500 truncate max-w-[120px]">
                      {{ authStore.user.email }}
                    </p>
                  </div>
                  <svg class="w-4 h-4 text-gray-400 transition-transform duration-200 group-hover:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                
                <!-- Dropdown Menu -->
                <div class="absolute right-0 mt-2 w-72 bg-white rounded-soft-lg shadow-md border border-jurnal-charcoal-200 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50 overflow-hidden">
                  <!-- User Info Section -->
                  <div class="p-4 bg-jurnal-charcoal-50 border-b border-jurnal-charcoal-200">
                    <div class="flex items-center gap-3 mb-3">
                      <div class="w-14 h-14 rounded-full bg-jurnal-teal-600 flex items-center justify-center shadow-sm ring-2 ring-white">
                        <span class="text-base font-bold text-white">
                          {{ getInitials(authStore.user?.full_name || '') }}
                        </span>
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold text-jurnal-charcoal-800 truncate mb-0.5">
                          {{ authStore.user?.full_name }}
                        </p>
                        <p v-if="authStore.user?.email" class="text-xs text-jurnal-charcoal-600 truncate mb-0.5">
                          {{ authStore.user.email }}
                        </p>
                        <p v-if="authStore.user?.phone_number" class="text-xs text-jurnal-charcoal-600 truncate">
                          {{ authStore.user.phone_number }}
                        </p>
                      </div>
                    </div>
                    <!-- Badges -->
                    <div class="flex items-center gap-2 flex-wrap">
                      <span 
                        class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                        :class="getAuthProviderBadgeClass(authStore.user?.auth_provider)"
                      >
                        {{ getAuthProviderLabel(authStore.user?.auth_provider) }}
                      </span>
                      <span
                        v-if="authStore.user?.phone_verified"
                        class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-700"
                      >
                        <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                        </svg>
                        Terverifikasi
                      </span>
                    </div>
                  </div>
                  
                  <!-- Menu Items -->
                  <div class="py-1.5">
                    <NuxtLink 
                      to="/profile" 
                      class="flex items-center gap-3 px-4 py-2.5 text-sm text-jurnal-charcoal-700 hover:bg-jurnal-charcoal-50 transition-colors group/item"
                    >
                      <div class="w-8 h-8 rounded-soft bg-jurnal-teal-50 flex items-center justify-center group-hover/item:bg-jurnal-teal-100 transition-colors">
                        <svg class="w-4 h-4 text-jurnal-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                      </div>
                      <span class="font-medium">Edit Profil</span>
                    </NuxtLink>
                    <button 
                      @click="handleLogout" 
                      class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors group/item"
                    >
                      <div class="w-8 h-8 rounded-lg bg-red-50 flex items-center justify-center group-hover/item:bg-red-100 transition-colors">
                        <svg class="w-4 h-4 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                        </svg>
                      </div>
                      <span class="font-medium">Logout</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <slot />
      </main>
    </div>

    <!-- Mobile Bottom Navigation -->
    <BottomNav class="lg:hidden" />
  </div>
</template>

<script setup>
import logoImage from '~/assets/images/logo.png'

const authStore = useAuthStore()
const logoUrl = ref(null)

onMounted(() => {
  // Set logo URL only on client to avoid hydration mismatch
  logoUrl.value = logoImage
})

const handleLogout = () => {
  authStore.logout()
}

const getInitials = (name) => {
  if (!name) return 'U'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

const getAuthProviderLabel = (provider) => {
  if (!provider || typeof provider !== 'string') return 'Unknown'
  const labels = {
    'phone': 'WhatsApp',
    'google': 'Google',
    'phone_google': 'WhatsApp & Google',
    'email': 'Email'
  }
  return labels[provider] || 'Unknown'
}

const getAuthProviderBadgeClass = (provider) => {
  if (!provider || typeof provider !== 'string') return 'bg-gray-100 text-gray-800'
  const classes = {
    'phone': 'bg-green-100 text-green-800',
    'google': 'bg-blue-100 text-blue-800',
    'phone_google': 'bg-purple-100 text-purple-800',
    'email': 'bg-gray-100 text-gray-800'
  }
  return classes[provider] || 'bg-gray-100 text-gray-800'
}
</script>
