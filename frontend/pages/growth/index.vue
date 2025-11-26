<template>
  <div>
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Pertumbuhan</h1>
        <p class="text-gray-600 mt-2">Pantau pertumbuhan fisik anak Anda</p>
      </div>
      <NuxtLink to="/growth/add" class="px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition flex items-center space-x-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span>Tambah Pengukuran</span>
      </NuxtLink>
    </div>

    <!-- No Child Selected -->
    <div v-if="!childStore.selectedChild" class="bg-amber-50 border border-amber-200 rounded-xl p-8 text-center">
      <div class="text-5xl mb-4">⚠️</div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Pilih Anak Terlebih Dahulu</h2>
      <p class="text-gray-600 mb-4">Silakan pilih anak dari dropdown di header untuk melihat data pertumbuhan</p>
      <NuxtLink to="/children" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition">
        Kelola Profil Anak
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="measurementStore.loading" class="bg-white rounded-xl shadow-sm p-12 text-center mb-8">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto mb-4"></div>
      <p class="text-gray-600">Memuat data pengukuran...</p>
    </div>

    <!-- Latest Measurement Summary -->
    <div v-if="measurementStore.latestMeasurement && !measurementStore.loading" class="bg-gradient-to-r from-emerald-500 to-emerald-600 rounded-xl shadow-lg p-8 text-white mb-8">
      <h2 class="text-2xl font-bold mb-4">Pengukuran Terakhir</h2>
      <div class="grid md:grid-cols-3 gap-6">
        <div>
          <p class="text-emerald-100 text-sm">Berat Badan</p>
          <p class="text-3xl font-bold">{{ measurementStore.latestMeasurement.weight }} kg</p>
          <p class="text-emerald-100 text-sm mt-1">{{ measurementStore.latestMeasurement.nutritional_status || '-' }}</p>
          <p v-if="isValidNumber(measurementStore.latestMeasurement.weight_for_age_zscore)" class="text-emerald-200 text-xs mt-1">
            Z-Score: {{ formatZScore(measurementStore.latestMeasurement.weight_for_age_zscore) }}
          </p>
        </div>
        <div>
          <p class="text-emerald-100 text-sm">Tinggi Badan</p>
          <p class="text-3xl font-bold">{{ measurementStore.latestMeasurement.height }} cm</p>
          <p class="text-emerald-100 text-sm mt-1">{{ measurementStore.latestMeasurement.height_status || '-' }}</p>
          <p v-if="isValidNumber(measurementStore.latestMeasurement.height_for_age_zscore)" class="text-emerald-200 text-xs mt-1">
            Z-Score: {{ formatZScore(measurementStore.latestMeasurement.height_for_age_zscore) }}
          </p>
        </div>
        <div>
          <p class="text-emerald-100 text-sm">Tanggal</p>
          <p class="text-xl font-bold">{{ formatDate(measurementStore.latestMeasurement.measurement_date) }}</p>
          <div class="text-emerald-100 text-sm mt-1 space-y-0.5">
            <p>Usia: {{ measurementStore.latestMeasurement.age_display }}</p>
            <p v-if="childStore.selectedChild?.is_premature && latestMeasurementUsesCorrectedAge" 
               class="text-emerald-200 text-xs bg-emerald-700/30 px-2 py-0.5 rounded">
              Usia Koreksi digunakan
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- No Measurements State -->
    <div v-if="!measurementStore.loading && !measurementStore.hasMeasurements && measurementStore.measurements.length === 0" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="text-6xl mb-4">📊</div>
      <h2 class="text-2xl font-bold text-gray-900 mb-2">Belum Ada Data Pengukuran</h2>
      <p class="text-gray-600 mb-6">Mulai pantau pertumbuhan anak dengan menambahkan pengukuran pertama</p>
      <NuxtLink to="/growth/add" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition">
        Tambah Pengukuran Pertama
      </NuxtLink>
    </div>

    <!-- Measurements List - Always show if there are measurements -->
    <div v-if="childStore.selectedChild && !measurementStore.loading && (measurementStore.hasMeasurements || measurementStore.measurements.length > 0)" class="space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-bold text-gray-900">Riwayat Pengukuran ({{ measurementStore.measurements.length }})</h2>
        <NuxtLink to="/growth/charts" class="text-jurnal-teal-600 hover:text-jurnal-teal-700 font-medium flex items-center space-x-1">
          <span>Lihat Grafik</span>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </NuxtLink>
      </div>

      <div class="grid md:grid-cols-2 gap-6">
        <MeasurementCard 
          v-for="measurement in measurementStore.sortedMeasurements" 
          :key="measurement?.id || Math.random()"
          :measurement="measurement"
          @delete="confirmDelete"
        />
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold text-gray-900 mb-2">Hapus Pengukuran?</h3>
        <p class="text-gray-600 mb-6">Apakah Anda yakin ingin menghapus data pengukuran ini?</p>
        <div class="flex space-x-3">
          <button 
            @click="showDeleteModal = false"
            class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium"
          >
            Batal
          </button>
          <button 
            @click="deleteMeasurement"
            class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition font-medium"
          >
            Hapus
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { isValidNumber, formatZScore, safeFormatDate } from '~/composables/useSafeData'

definePageMeta({
  middleware: 'auth'
})

const childStore = useChildStore()
const measurementStore = useMeasurementStore()
const { calculateCorrectedAge } = useCorrectedAge()
const showDeleteModal = ref(false)
const measurementToDelete = ref(null)

// Check if latest measurement uses corrected age
const latestMeasurementUsesCorrectedAge = computed(() => {
  if (!childStore.selectedChild || !measurementStore.latestMeasurement) {
    return false
  }
  
  const ageInfo = calculateCorrectedAge(
    childStore.selectedChild.dob,
    childStore.selectedChild.is_premature || false,
    childStore.selectedChild.gestational_age,
    measurementStore.latestMeasurement.measurement_date
  )
  
  return ageInfo.useCorrected
})

const confirmDelete = (measurement) => {
  if (!measurement || !measurement.id) {
    console.error('Invalid measurement data:', measurement)
    return
  }
  measurementToDelete.value = measurement
  showDeleteModal.value = true
}

const deleteMeasurement = async () => {
  if (!measurementToDelete.value || !childStore.selectedChild) {
    alert('Data tidak valid untuk dihapus')
    return
  }

  try {
    const result = await measurementStore.deleteMeasurement(
      childStore.selectedChild.id, 
      measurementToDelete.value.id
    )
    
    if (result.success) {
      showDeleteModal.value = false
      measurementToDelete.value = null
    } else {
      alert('Gagal menghapus pengukuran: ' + (result.error || 'Terjadi kesalahan'))
    }
  } catch (error) {
    console.error('Error deleting measurement:', error)
    alert('Terjadi kesalahan saat menghapus pengukuran: ' + (error.message || 'Unknown error'))
  }
}

// Use utility function for date formatting
const formatDate = safeFormatDate

// Track if component is mounted
const isMounted = ref(false)

// Watch for child changes
const stopWatcher = watch(() => childStore.selectedChild, async (newChild) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  
  if (newChild) {
    try {
      await Promise.all([
        measurementStore.fetchMeasurements(newChild.id),
        measurementStore.fetchLatestMeasurement(newChild.id)
      ])
      
      // Guard again after async operation
      if (!isMounted.value) return
    } catch (error) {
      if (isMounted.value) {
        console.error('Error fetching measurements:', error)
      }
    }
  } else {
    // Guard: Don't update if component is unmounted
    if (isMounted.value) {
      measurementStore.clearMeasurements()
    }
  }
}, { immediate: false })

onMounted(async () => {
  isMounted.value = true
  
  if (childStore.selectedChild) {
    try {
      await Promise.all([
        measurementStore.fetchMeasurements(childStore.selectedChild.id),
        measurementStore.fetchLatestMeasurement(childStore.selectedChild.id)
      ])
      
      // Guard: Check if still mounted after async operation
      if (!isMounted.value) return
    } catch (error) {
      if (isMounted.value) {
        console.error('Error fetching initial measurements:', error)
      }
    }
  }
})

onUnmounted(() => {
  isMounted.value = false
  if (stopWatcher) {
    stopWatcher()
  }
})
</script>
