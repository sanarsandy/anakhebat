<template>
  <div class="p-6 max-w-4xl mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Perkembangan Anak</h1>
      <p class="text-gray-600 mt-2">Pantau tumbuh kembang si kecil berdasarkan Piramida Belajar</p>
    </div>

    <!-- No Child Selected -->
    <div v-if="!childStore.selectedChild" class="bg-amber-50 border border-amber-200 rounded-xl p-8 text-center">
      <div class="text-5xl mb-4">⚠️</div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Pilih Anak Terlebih Dahulu</h2>
      <p class="text-gray-600 mb-4">Silakan pilih anak dari dropdown di header untuk melihat perkembangan.</p>
    </div>

    <div v-else class="space-y-8">
      <!-- Pyramid Summary -->
      <div class="bg-white rounded-xl shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-6 text-center">Status Piramida Belajar</h2>
        
        <div v-if="milestoneStore.loading" class="h-64 flex items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
        </div>
        
        <div v-else>
          <PyramidVisualizer :progress="milestoneStore.summary?.progress_by_category || {}" />
          
          <!-- Warnings -->
          <div v-if="milestoneStore.summary?.pyramid_warnings?.length" class="mt-6 space-y-3">
            <div v-for="(warning, idx) in milestoneStore.summary.pyramid_warnings" :key="idx" 
                 class="bg-orange-50 border-l-4 border-orange-500 p-4 rounded-r-lg">
              <div class="flex">
                <div class="flex-shrink-0">
                  <Icon name="mdi:alert" class="h-5 w-5 text-orange-500" />
                </div>
                <div class="ml-3">
                  <p class="text-sm text-orange-700">{{ warning }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Red Flags -->
          <div v-if="milestoneStore.summary?.red_flags_detected?.length" class="mt-4 space-y-3">
            <div v-for="(flag, idx) in milestoneStore.summary.red_flags_detected" :key="idx" 
                 class="bg-red-50 border-l-4 border-red-500 p-4 rounded-r-lg">
              <div class="flex">
                <div class="flex-shrink-0">
                  <Icon name="mdi:alert-circle" class="h-5 w-5 text-red-500" />
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-red-800">Tanda Bahaya Terdeteksi</h3>
                  <p class="text-sm text-red-700 mt-1">{{ flag.question }}</p>
                  <p class="text-xs text-red-600 mt-2 font-semibold">Segera konsultasikan ke Dokter Anak.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <NuxtLink to="/development/assess" class="block p-6 bg-indigo-600 rounded-xl text-white hover:bg-indigo-700 transition shadow-lg transform hover:-translate-y-1">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-xl font-bold mb-2">Penilaian KPSP</h3>
              <p class="text-indigo-100 text-sm">Checklist berdasarkan Piramida Belajar.</p>
            </div>
            <Icon name="mdi:clipboard-check" class="h-12 w-12 opacity-80" />
          </div>
        </NuxtLink>

        <NuxtLink to="/development/assess-denver" class="block p-6 bg-gradient-to-br from-purple-600 to-indigo-600 rounded-xl text-white hover:from-purple-700 hover:to-indigo-700 transition shadow-lg transform hover:-translate-y-1">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-xl font-bold mb-2">Penilaian Denver II</h3>
              <p class="text-purple-100 text-sm">Checklist berdasarkan 4 domain Denver II.</p>
            </div>
            <Icon name="mdi:clipboard-list" class="h-12 w-12 opacity-80" />
          </div>
        </NuxtLink>

        <NuxtLink to="/development/history" class="block p-6 bg-white border border-gray-200 rounded-xl hover:border-indigo-300 transition hover:shadow-md">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-xl font-bold text-gray-900 mb-2">Riwayat Penilaian</h3>
              <p class="text-gray-600 text-sm">Lihat grafik kemajuan dari waktu ke waktu.</p>
            </div>
            <Icon name="mdi:history" class="h-12 w-12 text-gray-400" />
          </div>
        </NuxtLink>

        <NuxtLink to="/development/denver-ii" class="block p-6 bg-white border border-purple-200 rounded-xl hover:border-purple-300 transition hover:shadow-md">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-xl font-bold text-gray-900 mb-2">Grafik Denver II</h3>
              <p class="text-gray-600 text-sm">Visualisasi perkembangan berdasarkan 4 domain.</p>
            </div>
            <Icon name="mdi:chart-line" class="h-12 w-12 text-purple-400" />
          </div>
        </NuxtLink>
      </div>

      <!-- Rekomendasi Stimulasi Section -->
      <div v-if="recommendationStore.loading" class="bg-white rounded-xl shadow-sm p-6">
        <div class="flex items-center justify-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          <span class="ml-3 text-gray-600">Memuat rekomendasi...</span>
        </div>
      </div>
      
      <!-- Recommendations Section -->
      <div v-else-if="recommendationStore.hasRecommendations" class="bg-white rounded-xl shadow-sm p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <span class="text-3xl">💡</span>
            Rekomendasi Stimulasi
          </h2>
        </div>
        
        <p class="text-sm text-gray-600 mb-6">
          Berdasarkan milestone yang belum tercapai untuk anak usia {{ recommendationStore.ageMonths }} bulan
        </p>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <RecommendationCard 
            v-for="rec in recommendationStore.recommendations" 
            :key="rec.content.id"
            :recommendation="rec"
          />
        </div>
      </div>
      
      <!-- Empty State - No Recommendations -->
      <div v-else-if="!recommendationStore.loading" class="bg-white rounded-xl shadow-sm p-6">
        <div class="text-center py-8">
          <span class="text-4xl mb-4 block">✨</span>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Belum Ada Rekomendasi</h3>
          <p class="text-sm text-gray-600">
            Lakukan penilaian perkembangan terlebih dahulu untuk mendapatkan rekomendasi stimulasi yang sesuai.
          </p>
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
const milestoneStore = useMilestoneStore()
const recommendationStore = useRecommendationStore()

onMounted(async () => {
  if (childStore.selectedChild) {
    await Promise.all([
      milestoneStore.fetchSummary(childStore.selectedChild.id),
      recommendationStore.fetchRecommendations(childStore.selectedChild.id)
    ])
  }
})

// Track if component is mounted
const isMounted = ref(false)

// Watch for child change
const stopWatcher = watch(() => childStore.selectedChild, async (newChild) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  
  if (newChild) {
    try {
      await Promise.all([
        milestoneStore.fetchSummary(newChild.id),
        recommendationStore.fetchRecommendations(newChild.id)
      ])
      
      // Guard again after async operation
      if (!isMounted.value) return
    } catch (error) {
      if (isMounted.value) {
        console.error('Error fetching data:', error)
      }
    }
  } else {
    recommendationStore.clearRecommendations()
  }
})

onMounted(() => {
  isMounted.value = true
})

onUnmounted(() => {
  isMounted.value = false
  stopWatcher()
})

// Force update for HMR
console.log('Development page loaded')
</script>
