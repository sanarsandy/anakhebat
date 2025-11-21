<template>
  <div class="relative">
    <button 
      @click="isOpen = !isOpen"
      class="flex items-center space-x-2 px-4 py-2 bg-indigo-50 text-indigo-600 rounded-lg hover:bg-indigo-100 transition"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      </svg>
      <span class="text-sm font-medium">
        {{ selectedChildName || 'Pilih Anak' }}
      </span>
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown -->
    <div 
      v-if="isOpen" 
      class="absolute right-0 mt-2 w-64 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-20"
    >
      <div v-if="childStore.children.length === 0" class="px-4 py-3 text-sm text-gray-500">
        Belum ada data anak
      </div>
      <button
        v-for="child in childStore.children"
        :key="child.id"
        @click="selectChild(child)"
        class="w-full text-left px-4 py-2 hover:bg-gray-50 transition"
        :class="childStore.selectedChild?.id === child.id ? 'bg-indigo-50 text-indigo-600' : 'text-gray-700'"
      >
        <div class="font-medium">{{ child.name }}</div>
        <div class="text-xs text-gray-500 space-y-0.5">
          <div>{{ getAgeDisplay(child) }}</div>
          <div v-if="getAgeInfo(child).useCorrected" class="text-amber-600 text-xs">
            Koreksi: {{ getAgeInfo(child).correctedDisplay }}
          </div>
        </div>
      </button>
      <div class="border-t border-gray-200 mt-2 pt-2">
        <NuxtLink 
          to="/children/add" 
          class="flex items-center space-x-2 px-4 py-2 text-indigo-600 hover:bg-indigo-50 transition"
          @click="isOpen = false"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span class="text-sm font-medium">Tambah Anak Baru</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
const childStore = useChildStore()
const isOpen = ref(false)

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
})
</script>
