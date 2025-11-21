<template>
  <div class="p-6 max-w-7xl mx-auto">
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Grafik Denver II</h1>
        <p class="text-gray-600 mt-2">Visualisasi perkembangan anak berdasarkan 4 domain Denver II</p>
      </div>
      <NuxtLink to="/development" class="px-4 py-2 text-gray-700 hover:text-gray-900 font-medium">
        ← Kembali ke Dashboard
      </NuxtLink>
    </div>

    <!-- No Child Selected -->
    <div v-if="!childStore.selectedChild" class="bg-amber-50 border border-amber-200 rounded-xl p-8 text-center">
      <div class="text-5xl mb-4">⚠️</div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Pilih Anak Terlebih Dahulu</h2>
      <p class="text-gray-600 mb-4">Silakan pilih anak dari dropdown di header untuk melihat grafik Denver II</p>
      <NuxtLink to="/children" class="inline-block px-6 py-3 bg-indigo-600 text-white font-semibold rounded-lg hover:bg-indigo-700 transition">
        Kelola Profil Anak
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
      <p class="text-gray-600">Memuat data grafik Denver II...</p>
    </div>

    <!-- Denver II Grid Chart -->
    <div v-else-if="childStore.selectedChild && !loading" class="space-y-6">
      <!-- Header Info -->
      <div class="bg-white rounded-xl shadow-sm p-6">
        <div class="grid md:grid-cols-3 gap-4 text-sm">
          <div>
            <span class="font-semibold text-gray-700">Nama:</span>
            <span class="ml-2 text-gray-900">{{ childStore.selectedChild.name }}</span>
          </div>
          <div>
            <span class="font-semibold text-gray-700">Tanggal Lahir:</span>
            <span class="ml-2 text-gray-900">{{ formatDate(childStore.selectedChild.dateOfBirth) }}</span>
          </div>
          <div>
            <span class="font-semibold text-gray-700">Usia:</span>
            <span class="ml-2 text-gray-900">{{ currentAge }} bulan</span>
          </div>
        </div>
      </div>

      <!-- Legend -->
      <div class="bg-blue-50 border border-blue-200 rounded-xl p-4">
        <div class="flex flex-wrap items-center gap-6 text-sm">
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 bg-green-500 border border-gray-400"></div>
            <span>Pass (Lulus)</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 bg-red-500 border border-gray-400"></div>
            <span>Fail (Tidak Lulus)</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 bg-yellow-500 border border-gray-400"></div>
            <span>Sometimes (Kadang-kadang)</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 bg-gray-200 border border-gray-400"></div>
            <span>Belum Dinilai</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 rounded-full bg-red-600"></div>
            <span class="text-red-600 font-semibold">Red Flag</span>
          </div>
        </div>
      </div>

      <!-- Grid Chart -->
      <div class="bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="overflow-x-auto">
          <!-- Age Header -->
          <div class="sticky top-0 bg-gray-50 border-b border-gray-200 z-10">
            <div class="flex" :style="{ minWidth: `${chartWidth}px` }">
              <div class="w-64 flex-shrink-0 border-r border-gray-200 p-2 text-xs font-semibold text-gray-700">
                Domain / Milestone
              </div>
              <div class="flex-1 relative">
                <div class="absolute inset-0 flex">
                  <div 
                    v-for="age in ageMarkers" 
                    :key="age"
                    class="flex-1 border-r border-gray-300 text-center text-xs text-gray-600"
                    :style="{ width: `${getAgePosition(age)}%` }"
                  >
                    <div class="font-semibold">{{ age }}</div>
                    <div class="text-xs text-gray-500">{{ age < 24 ? 'bln' : 'thn' }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Domain Sections -->
          <div v-for="(domainKey, domainIndex) in domainOrder" :key="domainKey">
            <div 
              v-if="gridData && gridData.domains && gridData.domains[domainKey] && gridData.domains[domainKey].length > 0"
              class="border-b border-gray-200"
            >
              <!-- Domain Header -->
              <div class="bg-gray-100 border-b border-gray-300">
                <div class="flex" :style="{ minWidth: `${chartWidth}px` }">
                  <div class="w-64 flex-shrink-0 border-r border-gray-300 p-3 font-bold text-gray-900">
                    {{ getDomainName(domainKey) }}
                  </div>
                  <div class="flex-1"></div>
                </div>
              </div>

              <!-- Milestones -->
              <div 
                v-for="(milestone, idx) in (gridData && gridData.domains ? gridData.domains[domainKey] : [])" 
                :key="milestone.id"
                class="border-b border-gray-100 hover:bg-gray-50 transition"
              >
                <div class="flex items-center" :style="{ minWidth: `${chartWidth}px` }">
                  <!-- Milestone Label -->
                  <div class="w-64 flex-shrink-0 border-r border-gray-200 p-2 text-sm">
                    <div class="flex items-center gap-2">
                      <span v-if="milestone.is_red_flag" class="w-2 h-2 rounded-full bg-red-600 flex-shrink-0"></span>
                      <span class="text-gray-800">{{ milestone.question }}</span>
                    </div>
                  </div>

                  <!-- Age Bar -->
                  <div class="flex-1 relative h-12">
                    <div class="absolute inset-0 flex items-center">
                      <!-- Age Range Bar -->
                      <div 
                        v-if="milestone.min_age_range !== null && milestone.max_age_range !== null"
                        class="h-6 rounded"
                        :class="getBarColor(milestone.assessment_status)"
                        :style="{
                          left: `${getAgePosition(milestone.min_age_range)}%`,
                          width: `${getAgePosition(milestone.max_age_range) - getAgePosition(milestone.min_age_range)}%`,
                          border: '1px solid #9CA3AF'
                        }"
                      ></div>
                      
                      <!-- Current Age Indicator -->
                      <div 
                        v-if="currentAge >= 0"
                        class="absolute top-0 bottom-0 w-0.5 bg-blue-600 z-10"
                        :style="{ left: `${getAgePosition(currentAge)}%` }"
                      >
                        <div class="absolute -top-1 -left-1 w-3 h-3 bg-blue-600 rounded-full border-2 border-white"></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No Data Message -->
      <div v-if="!loading && (!gridData || !gridData.domains || Object.keys(gridData.domains).length === 0)" class="bg-white rounded-xl shadow-sm p-12 text-center">
        <div class="text-6xl mb-4">📊</div>
        <h3 class="text-xl font-bold text-gray-900 mb-2">Belum Ada Data Milestone</h3>
        <p class="text-gray-600 mb-6">Silakan isi penilaian Denver II terlebih dahulu untuk melihat grafik.</p>
        <NuxtLink to="/development/assess-denver" class="inline-block px-6 py-3 bg-indigo-600 text-white font-bold rounded-lg hover:bg-indigo-700">
          Mulai Penilaian Denver II
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useChildStore } from '~/stores/child'
import { useAuthStore } from '~/stores/auth'

const childStore = useChildStore()
const authStore = useAuthStore()
const apiBase = useApiUrl()

const loading = ref(true)
const gridData = ref(null)
const chartWidth = 1200

const domainOrder = ['PS', 'FM', 'L', 'GM']

// Age markers for the chart (0-72 months)
const ageMarkers = computed(() => {
  const markers = []
  // Months: 0, 2, 4, 6, 9, 12, 15, 18, 24
  markers.push(0, 2, 4, 6, 9, 12, 15, 18, 24)
  // Years: 3, 4, 5, 6 (in months: 36, 48, 60, 72)
  markers.push(36, 48, 60, 72)
  return markers
})

const currentAge = computed(() => {
  if (!childStore.selectedChild || !childStore.selectedChild.dateOfBirth) return 0
  const dob = new Date(childStore.selectedChild.dateOfBirth)
  const today = new Date()
  const months = (today.getFullYear() - dob.getFullYear()) * 12 + (today.getMonth() - dob.getMonth())
  return Math.max(0, months)
})

const getDomainName = (domain) => {
  const names = {
    'PS': 'PERSONAL SOSIAL',
    'FM': 'ADAPTIF - MOTORIK HALUS',
    'L': 'BAHASA',
    'GM': 'MOTORIK KASAR'
  }
  return names[domain] || domain
}

const getAgePosition = (age) => {
  // Convert age in months to percentage position (0-72 months = 0-100%)
  const maxAge = 72
  return Math.min(100, (age / maxAge) * 100)
}

const getBarColor = (status) => {
  switch (status) {
    case 'yes':
      return 'bg-green-500'
    case 'no':
      return 'bg-red-500'
    case 'sometimes':
      return 'bg-yellow-500'
    default:
      return 'bg-gray-200'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const fetchGridData = async () => {
  if (!childStore.selectedChild) return null

  try {
    const response = await $fetch(`${apiBase}/api/children/${childStore.selectedChild.id}/denver-ii/grid-data`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    return response
  } catch (error) {
    console.error('Failed to fetch Denver II grid data:', error)
    return null
  }
}

// Track if component is mounted
const isMounted = ref(false)

onMounted(async () => {
  isMounted.value = true
  
  if (childStore.selectedChild) {
    loading.value = true
    try {
      const data = await fetchGridData()
      
      // Guard: Don't update if component is unmounted
      if (!isMounted.value) return
      
      gridData.value = data
      console.log('Fetched Denver II grid data:', data)
    } catch (error) {
      if (isMounted.value) {
        console.error('Error loading Denver II grid:', error)
      }
    } finally {
      if (isMounted.value) {
        loading.value = false
      }
    }
  } else {
    loading.value = false
  }
})

const stopWatcher = watch(() => childStore.selectedChild, async (newChild) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  
  if (newChild) {
    loading.value = true
    try {
      const data = await fetchGridData()
      
      // Guard again after async operation
      if (!isMounted.value) return
      
      gridData.value = data
    } catch (error) {
      if (isMounted.value) {
        console.error('Error reloading Denver II grid:', error)
      }
    } finally {
      if (isMounted.value) {
        loading.value = false
      }
    }
  }
}, { immediate: false })

onUnmounted(() => {
  isMounted.value = false
  stopWatcher()
})
</script>
