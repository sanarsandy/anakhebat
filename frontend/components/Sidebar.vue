<template>
  <aside class="fixed inset-y-0 left-0 w-64 bg-white border-r border-jurnal-charcoal-200 flex flex-col">
    <!-- Logo -->
    <div class="p-6 border-b border-jurnal-charcoal-200">
      <NuxtLink to="/dashboard" class="flex items-center gap-3">
        <img 
          v-if="logoUrl"
          :src="logoUrl" 
          alt="Jurnal Si Kecil" 
          class="h-10 w-auto object-contain"
        />
        <div>
          <h2 class="text-lg font-bold text-jurnal-teal-600">Jurnal Si Kecil</h2>
          <p class="text-xs text-jurnal-charcoal-500">Tumbuh Kembang</p>
        </div>
      </NuxtLink>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 p-4 space-y-2">
      <NuxtLink 
        to="/dashboard" 
        class="flex items-center space-x-3 px-4 py-3 rounded-soft transition-colors"
        :class="isActive('/dashboard') ? 'bg-jurnal-teal-50 text-jurnal-teal-600' : 'text-jurnal-charcoal-700 hover:bg-jurnal-charcoal-50'"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>
        <span class="font-medium">Dashboard</span>
      </NuxtLink>

      <NuxtLink 
        to="/children" 
        class="flex items-center space-x-3 px-4 py-3 rounded-soft transition-colors"
        :class="isActive('/children') ? 'bg-jurnal-teal-50 text-jurnal-teal-600' : 'text-gray-700 hover:bg-gray-50'"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <span class="font-medium">Profil Anak</span>
      </NuxtLink>

      <NuxtLink 
        to="/growth" 
        class="flex items-center space-x-3 px-4 py-3 rounded-soft transition-colors"
        :class="isActive('/growth') ? 'bg-jurnal-teal-50 text-jurnal-teal-600' : 'text-gray-700 hover:bg-gray-50'"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <span class="font-medium">Pertumbuhan</span>
      </NuxtLink>

      <NuxtLink 
        to="/development" 
        class="flex items-center space-x-3 px-4 py-3 rounded-soft transition-colors"
        :class="isActive('/development') ? 'bg-jurnal-teal-50 text-jurnal-teal-600' : 'text-gray-700 hover:bg-gray-50'"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span class="font-medium">Perkembangan</span>
      </NuxtLink>
    </nav>

    <!-- Footer -->
    <div class="p-4 border-t border-jurnal-charcoal-200">
      <p class="text-xs text-jurnal-charcoal-500 text-center">WHO Standards • CDC Milestones</p>
    </div>
  </aside>
</template>

<script setup>
import logoImage from '~/assets/images/logo.png'

const route = useRoute()
const logoUrl = ref(null)

onMounted(() => {
  // Set logo URL only on client to avoid hydration mismatch
  logoUrl.value = logoImage
})

const isActive = (path) => {
  return route.path.startsWith(path)
}
</script>
