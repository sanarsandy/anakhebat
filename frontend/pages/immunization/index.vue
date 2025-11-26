<template>
  <div class="p-6 max-w-7xl mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Jadwal Imunisasi</h1>
      <p class="text-gray-600 mt-2">Pantau jadwal imunisasi si kecil berdasarkan rekomendasi IDAI</p>
    </div>

    <!-- No Child Selected -->
    <div v-if="!childStore.selectedChild" class="bg-amber-50 border border-amber-200 rounded-xl p-8 text-center">
      <div class="text-5xl mb-4">⚠️</div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Pilih Anak Terlebih Dahulu</h2>
      <p class="text-gray-600 mb-4">Silakan pilih anak dari dropdown di header untuk melihat jadwal imunisasi.</p>
    </div>

    <div v-else class="space-y-6">
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
        <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-jurnal-teal-500">
          <p class="text-gray-500 text-xs font-medium uppercase mb-1">Total</p>
          <p class="text-2xl font-bold text-gray-900">{{ immunizationStore.summary?.total || 0 }}</p>
        </div>
        <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-green-500">
          <p class="text-gray-500 text-xs font-medium uppercase mb-1">Selesai</p>
          <p class="text-2xl font-bold text-gray-900">{{ immunizationStore.summary?.completed || 0 }}</p>
        </div>
        <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-yellow-500">
          <p class="text-gray-500 text-xs font-medium uppercase mb-1">Menunggu</p>
          <p class="text-2xl font-bold text-gray-900">{{ immunizationStore.summary?.pending || 0 }}</p>
        </div>
        <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-red-500">
          <p class="text-gray-500 text-xs font-medium uppercase mb-1">Terlambat</p>
          <p class="text-2xl font-bold text-gray-900">{{ immunizationStore.summary?.overdue || 0 }}</p>
        </div>
        <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-blue-500">
          <p class="text-gray-500 text-xs font-medium uppercase mb-1">Akan Datang</p>
          <p class="text-2xl font-bold text-gray-900">{{ immunizationStore.summary?.upcoming || 0 }}</p>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="immunizationStore.loading" class="bg-white rounded-xl shadow-sm p-12 text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="text-gray-600 mt-4">Memuat jadwal imunisasi...</p>
      </div>

      <!-- Schedule List -->
      <div v-else-if="immunizationStore.hasSchedule" class="space-y-4">
        <!-- Filter Tabs -->
        <div class="bg-white rounded-lg shadow-sm p-4">
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              @click="activeTab = tab.key"
              class="px-4 py-2 rounded-lg text-sm font-medium transition"
              :class="activeTab === tab.key 
                ? 'bg-jurnal-teal-600 text-white' 
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
            >
              {{ tab.label }} ({{ tab.count }})
            </button>
          </div>
        </div>

        <!-- Immunization Cards -->
        <div v-for="immunization in filteredImmunizations" :key="immunization.schedule.id" 
             class="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition p-6">
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <!-- Header -->
              <div class="flex items-center gap-3 mb-3">
                <h3 class="text-lg font-bold text-gray-900">
                  {{ immunization.schedule.name_id || immunization.schedule.name }}
                </h3>
                <span v-if="immunization.schedule.dose_number > 1" 
                      class="px-2 py-0.5 bg-gray-100 text-gray-700 text-xs rounded">
                  Dosis {{ immunization.schedule.dose_number }}
                </span>
                <span class="px-2 py-0.5 rounded text-xs font-semibold"
                      :class="getStatusBadgeClass(immunization.status)">
                  {{ getStatusLabel(immunization.status) }}
                </span>
              </div>

              <!-- Info -->
              <div class="space-y-1 text-sm text-gray-600 mb-4">
                <p v-if="immunization.schedule.description" class="text-gray-500">
                  {{ immunization.schedule.description }}
                </p>
                <p>
                  <span class="font-medium">Usia Optimal:</span> 
                  {{ immunization.schedule.age_optimal_months || immunization.schedule.age_optimal_days ? 
                     (immunization.schedule.age_optimal_months ? `${immunization.schedule.age_optimal_months} bulan` : `${Math.floor((immunization.schedule.age_optimal_days || 0) / 30)} bulan`) 
                     : '-' }}
                </p>
                <p v-if="immunization.due_date">
                  <span class="font-medium">Jadwal:</span> 
                  {{ formatDate(immunization.due_date) }}
                  <span v-if="immunization.days_until_due" class="text-blue-600">
                    ({{ immunization.days_until_due }} hari lagi)
                  </span>
                  <span v-if="immunization.days_overdue" class="text-red-600">
                    (Terlambat {{ immunization.days_overdue }} hari)
                  </span>
                </p>
              </div>

              <!-- Completed Info -->
              <div v-if="immunization.status === 'completed' && immunization.record" 
                   class="bg-green-50 border border-green-200 rounded-lg p-3 mb-3">
                <p class="text-sm font-semibold text-green-800 mb-1">✅ Sudah Diberikan</p>
                <div class="text-xs text-green-700 space-y-1">
                  <p><span class="font-medium">Tanggal:</span> {{ formatDate(immunization.record.given_date) }}</p>
                  <p v-if="immunization.record.location">
                    <span class="font-medium">Lokasi:</span> {{ immunization.record.location }}
                  </p>
                  <p v-if="immunization.record.doctor_name">
                    <span class="font-medium">Dokter:</span> {{ immunization.record.doctor_name }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Action Button -->
            <div class="ml-4">
              <button
                v-if="immunization.status !== 'completed'"
                @click="openRecordModal(immunization)"
                class="px-4 py-2 bg-jurnal-teal-600 text-white text-sm font-semibold rounded-lg hover:bg-jurnal-teal-700 transition"
              >
                Tandai Sudah
              </button>
              <span v-else class="text-green-600 text-sm font-semibold">✓ Selesai</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="bg-white rounded-xl shadow-sm p-12 text-center">
        <div class="text-5xl mb-4">💉</div>
        <h3 class="text-xl font-bold text-gray-900 mb-2">Belum Ada Jadwal Imunisasi</h3>
        <p class="text-gray-600">Jadwal imunisasi akan muncul di sini.</p>
      </div>
    </div>

    <!-- Record Modal -->
    <div v-if="showRecordModal && selectedImmunization" 
         class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
         @click.self="closeRecordModal">
      <div class="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-4">
          Catat Imunisasi: {{ selectedImmunization.schedule.name_id || selectedImmunization.schedule.name }}
        </h2>
        
        <form @submit.prevent="submitRecord">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Pemberian *</label>
              <input
                v-model="recordForm.given_date"
                type="date"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi</label>
              <input
                v-model="recordForm.location"
                type="text"
                placeholder="Contoh: RS. Contoh, Puskesmas..."
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Fasilitas Kesehatan</label>
              <input
                v-model="recordForm.healthcare_facility"
                type="text"
                placeholder="Nama rumah sakit/puskesmas"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nama Dokter</label>
              <input
                v-model="recordForm.doctor_name"
                type="text"
                placeholder="Nama dokter yang memberikan"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nomor Batch Vaksin</label>
              <input
                v-model="recordForm.vaccine_batch_number"
                type="text"
                placeholder="Nomor batch vaksin (opsional)"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Catatan</label>
              <textarea
                v-model="recordForm.notes"
                rows="3"
                placeholder="Catatan tambahan (opsional)"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              ></textarea>
            </div>
          </div>
          
          <div class="flex gap-3 mt-6">
            <button
              type="button"
              @click="closeRecordModal"
              class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="immunizationStore.loading"
              class="flex-1 px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition disabled:opacity-50"
            >
              Simpan
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const childStore = useChildStore()
const immunizationStore = useImmunizationStore()

const activeTab = ref('all')
const showRecordModal = ref(false)
const selectedImmunization = ref(null)

const recordForm = ref({
  given_date: '',
  location: '',
  healthcare_facility: '',
  doctor_name: '',
  vaccine_batch_number: '',
  notes: ''
})

const tabs = computed(() => {
  return [
    { key: 'all', label: 'Semua', count: immunizationStore.summary?.total || 0 },
    { key: 'completed', label: 'Selesai', count: immunizationStore.summary?.completed || 0 },
    { key: 'pending', label: 'Menunggu', count: immunizationStore.summary?.pending || 0 },
    { key: 'overdue', label: 'Terlambat', count: immunizationStore.summary?.overdue || 0 },
    { key: 'upcoming', label: 'Akan Datang', count: immunizationStore.summary?.upcoming || 0 },
  ]
})

const filteredImmunizations = computed(() => {
  if (activeTab.value === 'all') {
    return immunizationStore.schedule
  }
  
  // Get current age untuk filter
  const currentAgeMonths = immunizationStore.ageMonths || 0
  
  // Helper: cek apakah sudah terlewat
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  const isTerlewat = (item) => {
    // Sudah selesai tidak terhitung terlewat
    if (item.status === 'completed') return false
    
    // Cek status overdue dari backend
    if (item.status === 'overdue') return true
    
    // Cek due_date sudah lewat
    if (item.due_date) {
      const dueDate = new Date(item.due_date)
      dueDate.setHours(0, 0, 0, 0)
      if (dueDate <= today) return true
    }
    
    // Cek usia optimal sudah terlewat
    const ageOptimal = item.schedule.age_optimal_months || 0
    if (ageOptimal < currentAgeMonths) return true
    
    return false
  }
  
  const isAkanDatang = (item) => {
    // Sudah selesai tidak terhitung akan datang
    if (item.status === 'completed') return false
    
    // Jadwal bulan ini = usia optimal = umur anak saat ini
    const ageOptimal = item.schedule.age_optimal_months || 0
    return ageOptimal === currentAgeMonths
  }
  
  // Filter berdasarkan tab
  if (activeTab.value === 'overdue') {
    // Tab "Terlambat": yang sudah lewat umurnya dan belum ditandai
    return immunizationStore.schedule.filter(i => isTerlewat(i))
  } else if (activeTab.value === 'upcoming') {
    // Tab "Akan Datang": yang jadwalnya bulan ini
    return immunizationStore.schedule.filter(i => isAkanDatang(i))
  } else {
    // Tab lainnya: filter berdasarkan status dari backend
    return immunizationStore.schedule.filter(i => i.status === activeTab.value)
  }
})

const getStatusLabel = (status) => {
  const labels = {
    'completed': 'Selesai',
    'pending': 'Menunggu',
    'overdue': 'Terlambat',
    'upcoming': 'Akan Datang'
  }
  return labels[status] || status
}

const getStatusBadgeClass = (status) => {
  const classes = {
    'completed': 'bg-green-100 text-green-700',
    'pending': 'bg-yellow-100 text-yellow-700',
    'overdue': 'bg-red-100 text-red-700',
    'upcoming': 'bg-blue-100 text-blue-700'
  }
  return classes[status] || 'bg-gray-100 text-gray-700'
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

const openRecordModal = (immunization) => {
  selectedImmunization.value = immunization
  recordForm.value = {
    given_date: new Date().toISOString().split('T')[0],
    location: '',
    healthcare_facility: '',
    doctor_name: '',
    vaccine_batch_number: '',
    notes: ''
  }
  showRecordModal.value = true
}

const closeRecordModal = () => {
  showRecordModal.value = false
  selectedImmunization.value = null
}

const submitRecord = async () => {
  if (!selectedImmunization.value || !childStore.selectedChild) return
  
  const result = await immunizationStore.recordImmunization(childStore.selectedChild.id, {
    immunization_schedule_id: selectedImmunization.value.schedule.id,
    ...recordForm.value
  })
  
  if (result.success) {
    closeRecordModal()
    // Show success message
    alert('Imunisasi berhasil dicatat!')
  } else {
    alert('Gagal mencatat imunisasi: ' + (result.error || 'Unknown error'))
  }
}

// Track if component is mounted
const isMounted = ref(false)

// Watch for child change
const stopWatcher = watch(() => childStore.selectedChild, async (newChild) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  
  if (newChild) {
    try {
      await immunizationStore.fetchSchedule(newChild.id)
      
      // Guard again after async operation
      if (!isMounted.value) return
    } catch (error) {
      if (isMounted.value) {
        console.error('Error fetching immunization schedule:', error)
      }
    }
  } else {
    if (isMounted.value) {
      immunizationStore.clearSchedule()
    }
  }
}, { immediate: false })

onMounted(async () => {
  isMounted.value = true
  
  if (childStore.selectedChild) {
    try {
      await immunizationStore.fetchSchedule(childStore.selectedChild.id)
      
      // Guard: Check if still mounted after async operation
      if (!isMounted.value) return
    } catch (error) {
      if (isMounted.value) {
        console.error('Error fetching initial immunization schedule:', error)
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

