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
      <NuxtLink to="/children" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition">
        Kelola Profil Anak
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto mb-4"></div>
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
        <div class="space-y-3">
          <div class="text-sm font-semibold text-gray-700 mb-2">Keterangan Grafik Denver II:</div>
          
          <!-- Percentile Zones -->
          <div class="flex flex-wrap items-center gap-4 text-xs">
            <div class="flex items-center gap-2">
              <div class="w-12 h-4 rounded" style="background: linear-gradient(to right, #ffffff, #e5e7eb, #9ca3af, #4b5563); border: 1px solid #9ca3af;"></div>
              <span>Bar Persentil: 25% → 50% → 75% → 90%</span>
            </div>
          </div>
          
          <!-- Status Markers -->
          <div class="flex flex-wrap items-center gap-4 text-xs">
            <div class="flex items-center gap-2">
              <span class="text-green-600 font-bold text-lg">✓</span>
              <span>Pass (Lulus)</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-red-600 font-bold text-lg">✗</span>
              <span>Fail (Tidak Lulus)</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-yellow-600 font-bold text-lg">R</span>
              <span>Refuse (Menolak)</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-gray-500 font-bold text-lg">—</span>
              <span>Belum Dinilai</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 rounded-full bg-red-600"></div>
              <span class="text-red-600 font-semibold">Red Flag</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-0.5 h-4 bg-blue-600"></div>
              <div class="w-3 h-3 rounded-full bg-blue-600 border-2 border-white"></div>
              <span>Usia Saat Ini</span>
            </div>
          </div>
          
          <!-- Domain Colors -->
          <div class="flex flex-wrap items-center gap-4 text-xs pt-2 border-t border-blue-300">
            <div class="flex items-center gap-2">
              <div class="w-4 h-4 rounded bg-blue-500"></div>
              <span><strong>PS</strong> - Personal Sosial</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-4 h-4 rounded bg-red-500"></div>
              <span><strong>FM</strong> - Adaptif Motorik Halus</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-4 h-4 rounded bg-green-500"></div>
              <span><strong>L</strong> - Bahasa</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-4 h-4 rounded bg-yellow-500"></div>
              <span><strong>GM</strong> - Motorik Kasar</span>
            </div>
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
              <div class="flex-1 relative h-12">
                <div class="absolute inset-0">
                  <div 
                    v-for="(age, idx) in ageMarkers" 
                    :key="age"
                    class="absolute text-center text-xs"
                    :style="{ 
                      left: `${getAgePosition(age)}%`,
                      transform: 'translateX(-50%)',
                      width: '60px'
                    }"
                  >
                    <div class="absolute top-0 left-1/2 -translate-x-1/2 w-px h-3 bg-gray-300"></div>
                    <div class="font-semibold pt-4 text-gray-600">{{ age }}</div>
                    <div class="text-xs text-gray-500">{{ age < 24 ? 'bln' : 'thn' }}</div>
                  </div>
                  <!-- Vertical lines between markers -->
                  <div 
                    v-for="(age, idx) in ageMarkers.slice(0, -1)" 
                    :key="`line-${age}`"
                    class="absolute top-0 bottom-0 w-px bg-gray-300"
                    :style="{ left: `${getAgePosition(age)}%` }"
                  ></div>
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
              <div class="border-b border-gray-300" :style="{ backgroundColor: getDomainColor(domainKey, 0.1) }">
                <div class="flex" :style="{ minWidth: `${chartWidth}px` }">
                  <div class="w-64 flex-shrink-0 border-r border-gray-300 p-3 font-bold" :style="{ color: getDomainColor(domainKey, 1) }">
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
                      <!-- Percentile Bar with 4 zones (25%, 50%, 75%, 90%) -->
                      <div 
                        v-if="milestone.age_25_percentile !== undefined && milestone.age_90_percentile !== undefined"
                        class="absolute h-6 rounded border border-gray-400"
                        :style="{
                          left: `${getAgePosition(milestone.age_25_percentile)}%`,
                          width: `${getAgePosition(milestone.age_90_percentile) - getAgePosition(milestone.age_25_percentile)}%`,
                          background: getPercentileGradient(milestone, domainKey)
                        }"
                      >
                        <!-- Zone dividers -->
                        <div 
                          class="absolute top-0 bottom-0 w-px bg-gray-600 opacity-30"
                          :style="{ left: `${((milestone.age_50_percentile - milestone.age_25_percentile) / (milestone.age_90_percentile - milestone.age_25_percentile)) * 100}%` }"
                        ></div>
                        <div 
                          class="absolute top-0 bottom-0 w-px bg-gray-600 opacity-30"
                          :style="{ left: `${((milestone.age_75_percentile - milestone.age_25_percentile) / (milestone.age_90_percentile - milestone.age_25_percentile)) * 100}%` }"
                        ></div>
                      </div>
                      
                      <!-- Status Marker at Current Age -->
                      <div 
                        v-if="currentAge >= 0 && milestone.assessment_status"
                        class="absolute top-1/2 -translate-y-1/2 z-20"
                        :style="{ left: `${getAgePosition(currentAge)}%`, transform: 'translate(-50%, -50%)' }"
                      >
                        <span 
                          class="text-xl font-bold drop-shadow-lg"
                          :class="getStatusMarkerClass(milestone.assessment_status)"
                        >
                          {{ getStatusMarker(milestone.assessment_status) }}
                        </span>
                      </div>
                      
                      <!-- Current Age Indicator -->
                      <div 
                        v-if="currentAge >= 0"
                        class="absolute top-0 bottom-0 w-0.5 bg-blue-600 z-10"
                        :style="{ left: `${getAgePosition(currentAge)}%` }"
                      >
                        <div class="absolute -top-1 -left-1 w-3 h-3 bg-blue-600 rounded-full border-2 border-white shadow-sm"></div>
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
        <NuxtLink to="/development/assess-denver" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-bold rounded-lg hover:bg-jurnal-teal-700">
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
  // Denver II uses logarithmic scale for 0-24 months, then linear for 24-72 months
  if (age <= 24) {
    // Logarithmic scale for 0-24 months (maps to 0-50% of chart)
    // Using log base 25 to map 0-24 months to 0-50%
    if (age === 0) return 0
    const logValue = Math.log(age + 1) / Math.log(25)
    return Math.min(50, logValue * 50)
  } else {
    // Linear scale for 24-72 months (maps to 50-100% of chart)
    const linearValue = 50 + ((age - 24) / 48) * 50
    return Math.min(100, linearValue)
  }
}

// Domain colors (Denver II standard)
const getDomainColor = (domain, opacity = 1) => {
  const colors = {
    'PS': `rgba(59, 130, 246, ${opacity})`,    // Blue
    'FM': `rgba(239, 68, 68, ${opacity})`,    // Red
    'L': `rgba(16, 185, 129, ${opacity})`,    // Green
    'GM': `rgba(245, 158, 11, ${opacity})`    // Yellow/Orange
  }
  return colors[domain] || `rgba(156, 163, 175, ${opacity})`
}

// Get gradient for percentile bar (25% -> 50% -> 75% -> 90%)
const getPercentileGradient = (milestone, domain) => {
  const baseColor = getDomainColor(domain, 1)
  // Extract RGB values
  const rgbMatch = baseColor.match(/\d+/g)
  if (!rgbMatch || rgbMatch.length < 3) {
    return 'linear-gradient(to right, #ffffff, #e5e7eb, #9ca3af, #4b5563)'
  }
  
  const r = parseInt(rgbMatch[0])
  const g = parseInt(rgbMatch[1])
  const b = parseInt(rgbMatch[2])
  
  // Create gradient: 25% (light) -> 50% -> 75% -> 90% (darker)
  return `linear-gradient(to right, 
    rgba(255, 255, 255, 0.9) 0%,
    rgba(${r}, ${g}, ${b}, 0.3) 33%,
    rgba(${r}, ${g}, ${b}, 0.6) 66%,
    rgba(${r}, ${g}, ${b}, 0.9) 100%
  )`
}

// Get status marker symbol
const getStatusMarker = (status) => {
  switch (status) {
    case 'yes':
      return '✓'
    case 'no':
      return '✗'
    case 'sometimes':
      return 'R'
    default:
      return '—'
  }
}

// Get status marker CSS class
const getStatusMarkerClass = (status) => {
  switch (status) {
    case 'yes':
      return 'text-green-600'
    case 'no':
      return 'text-red-600'
    case 'sometimes':
      return 'text-yellow-600'
    default:
      return 'text-gray-500'
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
  if (stopWatcher) {
    stopWatcher()
  }
})
</script>
