<template>
  <div ref="dropdownRef" class="relative">
    <button 
      @click="isOpen = !isOpen"
      class="flex items-center gap-2 px-3 py-1.5 bg-gray-50 hover:bg-gray-100 text-gray-700 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-jurnal-teal-500 focus:ring-offset-2 border border-gray-200/50"
    >
      <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <span class="text-sm font-medium max-w-[120px] truncate">
        {{ selectedChildName || 'Pilih Anak' }}
      </span>
      <svg class="w-3.5 h-3.5 text-gray-400 transition-transform duration-200" :class="isOpen ? 'rotate-180' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown -->
    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 scale-95 translate-y-1"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 translate-y-1"
    >
      <div 
        v-if="isOpen" 
        class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-lg border border-gray-100 py-1.5 z-30 overflow-hidden"
        @click.stop
      >
        <div v-if="childStore.children.length === 0" class="px-4 py-4 text-center">
          <div class="w-12 h-12 mx-auto mb-2 rounded-full bg-gray-100 flex items-center justify-center">
            <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <p class="text-sm text-gray-500 mb-3">Belum ada data anak</p>
        </div>
        <div v-else class="max-h-64 overflow-y-auto">
          <button
            v-for="child in childStore.children"
            :key="child.id"
            @click="selectChild(child)"
            class="w-full text-left px-4 py-2.5 hover:bg-gray-50 transition-colors group"
            :class="childStore.selectedChild?.id === child.id ? 'bg-jurnal-teal-50 border-l-2 border-jurnal-teal-500' : ''"
          >
            <div class="flex items-center justify-between mb-1">
              <span class="font-medium text-gray-900" :class="childStore.selectedChild?.id === child.id ? 'text-jurnal-teal-600' : ''">
                {{ child.name }}
              </span>
              <svg v-if="childStore.selectedChild?.id === child.id" class="w-4 h-4 text-jurnal-teal-600" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="text-xs text-gray-500 space-y-0.5">
              <div>{{ getAgeDisplay(child) }}</div>
              <div v-if="getAgeInfo(child).useCorrected" class="text-amber-600">
                Koreksi: {{ getAgeInfo(child).correctedDisplay }}
              </div>
            </div>
          </button>
        </div>
        <div class="border-t border-gray-100 mt-1 pt-1.5">
          <NuxtLink 
            to="/children/add" 
            class="flex items-center gap-3 px-4 py-2.5 text-sm text-jurnal-teal-600 hover:bg-jurnal-teal-50 transition-colors group/item"
            @click="isOpen = false"
          >
            <div class="w-8 h-8 rounded-lg bg-jurnal-teal-50 flex items-center justify-center group-hover/item:bg-jurnal-teal-100 transition-colors">
              <svg class="w-4 h-4 text-jurnal-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
            </div>
            <span class="font-medium">Tambah Anak Baru</span>
          </NuxtLink>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
const childStore = useChildStore()
const isOpen = ref(false)
const dropdownRef = ref(null)

const selectedChildName = computed(() => {
  return childStore.selectedChild?.name
})

const selectChild = (child) => {
  childStore.selectChild(child.id)
  isOpen.value = false
}

const { calculateCorrectedAge } = useCorrectedAge()

const getAgeInfo = (child) => {
  if (!child?.dob) {
    return {
      chronologicalDisplay: '',
      correctedDisplay: null,
      useCorrected: false
    }
  }
  
  return calculateCorrectedAge(
    child.dob,
    child.is_premature || false,
    child.gestational_age,
    undefined // Use current date
  )
}

const getAgeDisplay = (child) => {
  const info = getAgeInfo(child)
  if (info.useCorrected) {
    return `${info.chronologicalDisplay} (Usia Koreksi: ${info.correctedDisplay})`
  }
  return info.chronologicalDisplay
}

onMounted(() => {
  childStore.fetchChildren()
  
  // Close dropdown when clicking outside
  const handleClickOutside = (event) => {
    if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
      isOpen.value = false
    }
  }
  document.addEventListener('click', handleClickOutside)
  
  onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
  })
})
</script>
