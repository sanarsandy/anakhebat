<template>
  <nav class="fixed bottom-0 inset-x-0 bg-white border-t border-jurnal-charcoal-200 z-50 shadow-sm" style="padding-bottom: env(safe-area-inset-bottom, 0px);">
    <div class="flex justify-around items-center h-16 max-w-md mx-auto">
      <NuxtLink 
        to="/dashboard" 
        class="flex flex-col items-center justify-center flex-1 py-2 relative group transition-colors"
        :class="isActive('/dashboard') ? 'text-jurnal-teal-600' : 'text-jurnal-charcoal-500 hover:text-jurnal-charcoal-700'"
      >
        <div class="relative">
          <svg class="w-6 h-6 transition-transform group-active:scale-90" fill="none" stroke="currentColor" viewBox="0 0 24 24" :stroke-width="isActive('/dashboard') ? '2.5' : '2'">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <div v-if="isActive('/dashboard')" class="absolute -top-1 -right-1 w-2 h-2 bg-jurnal-teal-600 rounded-full"></div>
        </div>
        <span class="text-xs mt-1 font-medium" :class="isActive('/dashboard') ? 'text-jurnal-teal-600' : 'text-jurnal-charcoal-500'">Home</span>
      </NuxtLink>

      <NuxtLink 
        to="/children" 
        class="flex flex-col items-center justify-center flex-1 py-2 relative group transition-colors"
        :class="isActive('/children') ? 'text-jurnal-teal-600' : 'text-gray-500 hover:text-gray-700'"
      >
        <div class="relative">
          <svg class="w-6 h-6 transition-transform group-active:scale-90" fill="none" stroke="currentColor" viewBox="0 0 24 24" :stroke-width="isActive('/children') ? '2.5' : '2'">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <div v-if="isActive('/children')" class="absolute -top-1 -right-1 w-2 h-2 bg-jurnal-teal-600 rounded-full"></div>
        </div>
        <span class="text-xs mt-1 font-medium" :class="isActive('/children') ? 'text-jurnal-teal-600' : 'text-gray-500'">Anak</span>
      </NuxtLink>

      <NuxtLink 
        to="/growth" 
        class="flex flex-col items-center justify-center flex-1 py-2 relative group transition-colors"
        :class="isActive('/growth') ? 'text-jurnal-teal-600' : 'text-gray-500 hover:text-gray-700'"
      >
        <div class="relative">
          <svg class="w-6 h-6 transition-transform group-active:scale-90" fill="none" stroke="currentColor" viewBox="0 0 24 24" :stroke-width="isActive('/growth') ? '2.5' : '2'">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          <div v-if="isActive('/growth')" class="absolute -top-1 -right-1 w-2 h-2 bg-jurnal-teal-600 rounded-full"></div>
        </div>
        <span class="text-xs mt-1 font-medium" :class="isActive('/growth') ? 'text-jurnal-teal-600' : 'text-gray-500'">Tumbuh</span>
      </NuxtLink>

      <NuxtLink 
        to="/development" 
        class="flex flex-col items-center justify-center flex-1 py-2 relative group transition-colors"
        :class="isActive('/development') ? 'text-jurnal-teal-600' : 'text-gray-500 hover:text-gray-700'"
      >
        <div class="relative">
          <svg class="w-6 h-6 transition-transform group-active:scale-90" fill="none" stroke="currentColor" viewBox="0 0 24 24" :stroke-width="isActive('/development') ? '2.5' : '2'">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div v-if="isActive('/development')" class="absolute -top-1 -right-1 w-2 h-2 bg-jurnal-teal-600 rounded-full"></div>
        </div>
        <span class="text-xs mt-1 font-medium" :class="isActive('/development') ? 'text-jurnal-teal-600' : 'text-gray-500'">Kembang</span>
      </NuxtLink>

      <NuxtLink 
        to="/profile" 
        class="flex flex-col items-center justify-center flex-1 py-2 relative group transition-colors"
        :class="isActive('/profile') ? 'text-jurnal-teal-600' : 'text-gray-500 hover:text-gray-700'"
      >
        <div class="relative">
          <svg class="w-6 h-6 transition-transform group-active:scale-90" fill="none" stroke="currentColor" viewBox="0 0 24 24" :stroke-width="isActive('/profile') ? '2.5' : '2'">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <div v-if="isActive('/profile')" class="absolute -top-1 -right-1 w-2 h-2 bg-jurnal-teal-600 rounded-full"></div>
        </div>
        <span class="text-xs mt-1 font-medium" :class="isActive('/profile') ? 'text-jurnal-teal-600' : 'text-gray-500'">Profil</span>
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup>
const route = useRoute()

const isActive = (path) => {
  return route.path.startsWith(path)
}
</script>
