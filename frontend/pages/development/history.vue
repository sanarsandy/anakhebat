<template>
  <div class="p-6 max-w-4xl mx-auto">
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Riwayat Perkembangan</h1>
      <NuxtLink to="/development" class="text-indigo-600 hover:text-indigo-800 font-medium">
        &larr; Kembali ke Dashboard
      </NuxtLink>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-4 text-gray-500">Memuat riwayat...</p>
    </div>

    <div v-else-if="!assessments.length" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="text-6xl mb-4">📅</div>
      <h3 class="text-xl font-bold text-gray-900 mb-2">Belum Ada Riwayat</h3>
      <p class="text-gray-600 mb-6">Belum ada penilaian perkembangan yang tercatat untuk anak ini.</p>
      <NuxtLink to="/development/assess" class="inline-block px-6 py-3 bg-indigo-600 text-white font-bold rounded-lg hover:bg-indigo-700">
        Mulai Penilaian Sekarang
      </NuxtLink>
    </div>

    <div v-else class="space-y-6">
      <!-- Group by Date -->
      <div v-for="(group, date) in groupedAssessments" :key="date" class="bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-6 py-4 bg-gray-50 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-bold text-gray-900">{{ formatDate(date) }}</h3>
          <span class="text-sm text-gray-500">{{ group.length }} Item Dinilai</span>
        </div>
        
        <div class="divide-y divide-gray-100">
          <div v-for="item in group" :key="item.id" class="p-4 hover:bg-gray-50 transition flex items-start">
            <div class="flex-shrink-0 mt-1">
              <Icon v-if="item.status === 'yes'" name="mdi:check-circle" class="text-green-500 h-5 w-5" />
              <Icon v-else-if="item.status === 'sometimes'" name="mdi:minus-circle" class="text-yellow-500 h-5 w-5" />
              <Icon v-else name="mdi:close-circle" class="text-red-500 h-5 w-5" />
            </div>
            <div class="ml-4 flex-1">
              <p class="text-gray-900 font-medium">{{ item.milestone?.question || item['milestone.question'] || 'N/A' }}</p>
              <div class="flex items-center mt-1 space-x-2">
                <span class="text-xs px-2 py-0.5 rounded bg-gray-100 text-gray-600 capitalize">
                  {{ item.milestone?.category || item['milestone.category'] || 'N/A' }}
                </span>
                <span v-if="item.milestone?.pyramid_level || item['milestone.pyramid_level']" class="text-xs px-2 py-0.5 rounded bg-indigo-50 text-indigo-600">
                  Level {{ item.milestone?.pyramid_level || item['milestone.pyramid_level'] }}
                </span>
                <span v-if="item.milestone?.denver_domain || item['milestone.denver_domain']" class="text-xs px-2 py-0.5 rounded bg-purple-50 text-purple-600">
                  {{ item.milestone?.denver_domain || item['milestone.denver_domain'] }}
                </span>
                <span v-if="item.milestone?.source || item['milestone.source']" class="text-xs px-2 py-0.5 rounded bg-blue-50 text-blue-600">
                  {{ item.milestone?.source || item['milestone.source'] }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const childStore = useChildStore()
const loading = ref(false)
const assessments = ref([])

const groupedAssessments = computed(() => {
  const groups = {}
  assessments.value.forEach(a => {
    const date = a.assessment_date.split('T')[0]
    if (!groups[date]) groups[date] = []
    groups[date].push(a)
  })
  return groups
})

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const fetchHistory = async () => {
  if (!childStore.selectedChild) return
  
  loading.value = true
  try {
    const apiBase = useApiUrl()
    const authStore = useAuthStore()
    
    const data = await $fetch(`${apiBase}/api/children/${childStore.selectedChild.id}/assessments`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    
    assessments.value = data
  } catch (e) {
    console.error('Failed to fetch history:', e)
  } finally {
    loading.value = false
  }
}

// Track if component is mounted
const isMounted = ref(false)

onMounted(() => {
  isMounted.value = true
  fetchHistory()
})

const stopWatcher = watch(() => childStore.selectedChild, () => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  fetchHistory()
})

onUnmounted(() => {
  isMounted.value = false
  stopWatcher()
})
</script>
