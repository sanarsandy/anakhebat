<template>
  <div class="min-h-screen bg-gray-50 pb-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="mb-8 pt-6">
        <h1 class="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-gray-600 mt-2">Selamat datang di Tukem - Tumbuh Kembang Anak</p>
      </div>

      <!-- No Child State -->
      <div v-if="!childStore.hasChildren" class="bg-white rounded-xl shadow-sm p-12 text-center">
        <div class="text-6xl mb-4">👶</div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">Belum Ada Data Anak</h2>
        <p class="text-gray-600 mb-6">Mulai dengan menambahkan profil anak Anda</p>
        <NuxtLink to="/children/add" class="inline-block px-6 py-3 bg-indigo-600 text-white font-semibold rounded-lg hover:bg-indigo-700 transition">
          Tambah Anak Pertama
        </NuxtLink>
      </div>

      <!-- Has Children State -->
      <div v-else-if="isInitialized" class="space-y-6" :key="`dashboard-${childStore.selectedChild?.id || 'no-child'}`">
        <!-- Selected Child Info Card -->
        <div class="bg-gradient-to-r from-indigo-500 to-indigo-600 rounded-xl shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <h2 class="text-2xl font-bold mb-1">{{ childStore.selectedChild?.name }}</h2>
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <p class="text-indigo-100 text-sm">{{ ageInfo.chronologicalDisplay }}</p>
                  <span v-if="ageInfo.useCorrected" class="px-2 py-0.5 bg-amber-400/20 text-amber-100 text-xs rounded-full border border-amber-300/30">
                    Usia Koreksi: {{ ageInfo.correctedDisplay }}
                  </span>
                </div>
                <p v-if="childStore.selectedChild?.is_premature && ageInfo.useCorrected" class="text-indigo-200 text-xs">
                  Prematur ({{ childStore.selectedChild?.gestational_age }} minggu) - Menggunakan usia koreksi
                </p>
                <p class="text-indigo-100 text-sm">{{ childStore.selectedChild?.gender === 'male' ? 'Laki-laki' : 'Perempuan' }}</p>
              </div>
            </div>
            <div class="text-5xl ml-4">{{ childStore.selectedChild?.gender === 'male' ? '👦' : '👧' }}</div>
          </div>
        </div>

        <!-- Red Flag Alert (Conditional) - Show at top if exists -->
        <div v-if="milestoneStore.summary?.red_flags_detected?.length" class="bg-red-50 border-l-4 border-red-500 rounded-lg p-4 flex items-start shadow-sm">
          <div class="text-2xl mr-3 flex-shrink-0">🚨</div>
          <div class="flex-1">
            <h3 class="font-bold text-red-800 mb-1">Perhatian Diperlukan!</h3>
            <p class="text-red-700 text-sm mb-2">Terdeteksi {{ milestoneStore.summary.red_flags_detected.length }} tanda bahaya perkembangan. Segera cek detail di menu Perkembangan.</p>
            <NuxtLink to="/development" class="text-red-600 font-semibold text-sm hover:underline inline-flex items-center">
              Lihat Detail 
              <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </NuxtLink>
          </div>
        </div>

        <!-- Quick Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Status Pertumbuhan Card -->
          <div class="bg-white rounded-lg shadow-sm p-5 border-l-4 border-emerald-500">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <p class="text-gray-500 text-xs font-medium uppercase tracking-wide mb-2">Status Pertumbuhan</p>
                <p class="text-xl font-bold text-gray-900 mb-1">
                  {{ measurementStore.latestMeasurement?.nutritional_status || 'Belum ada data' }}
                </p>
                <p class="text-xs text-gray-600">
                  Tinggi: {{ measurementStore.latestMeasurement?.height_status || '-' }}
                </p>
              </div>
              <div class="text-2xl ml-3">📊</div>
            </div>
          </div>

          <!-- Milestone Card -->
          <div class="bg-white rounded-lg shadow-sm p-5 border-l-4 border-indigo-500">
            <div class="flex items-start justify-between">
              <div class="flex-1 min-w-0">
                <p class="text-gray-500 text-xs font-medium uppercase tracking-wide mb-2">Milestone Tercapai</p>
                <p class="text-xl font-bold text-gray-900 mb-1">
                  {{ milestoneStore.summary ? `${milestoneStore.summary.completed_milestones} / ${milestoneStore.summary.total_milestones}` : 'Belum ada data' }}
                </p>
                <p v-if="milestoneStore.summary" class="text-xs text-indigo-600 font-medium mb-2">
                  {{ Math.round((milestoneStore.summary.completed_milestones / milestoneStore.summary.total_milestones) * 100) || 0 }}% Selesai
                </p>
                <!-- Progress Bar -->
                <div v-if="milestoneStore.summary && milestoneStore.summary.total_milestones > 0" class="w-full bg-gray-200 rounded-full h-1.5">
                  <div 
                    class="bg-indigo-600 h-1.5 rounded-full transition-all duration-300" 
                    :style="{ width: `${Math.round((milestoneStore.summary.completed_milestones / milestoneStore.summary.total_milestones) * 100)}%` }"
                  ></div>
                </div>
              </div>
              <div class="text-2xl ml-3 flex-shrink-0">✅</div>
            </div>
          </div>

          <!-- Pengukuran Terakhir Card -->
          <div class="bg-white rounded-lg shadow-sm p-5 border-l-4 border-amber-500">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <p class="text-gray-500 text-xs font-medium uppercase tracking-wide mb-2">Pengukuran Terakhir</p>
                <div class="space-y-1">
                  <p class="text-lg font-bold text-gray-900">
                    {{ measurementStore.latestMeasurement?.weight ? measurementStore.latestMeasurement.weight + ' kg' : '-' }}
                  </p>
                  <p class="text-lg font-bold text-gray-900">
                    {{ measurementStore.latestMeasurement?.height ? measurementStore.latestMeasurement.height + ' cm' : '-' }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{ measurementStore.latestMeasurement?.measurement_date ? new Date(measurementStore.latestMeasurement.measurement_date).toLocaleDateString('id-ID') : '' }}
                  </p>
                </div>
              </div>
              <div class="text-2xl ml-3 flex-shrink-0">📏</div>
            </div>
          </div>
        </div>


        <!-- Next Immunization Widget -->
        <div v-if="!immunizationStore.loading" class="bg-white rounded-xl shadow-sm p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold text-gray-900 flex items-center gap-2">
              <span class="text-2xl">💉</span>
              Jadwal Imunisasi Berikutnya
            </h3>
            <NuxtLink to="/immunization" class="text-sm text-indigo-600 hover:text-indigo-700 font-medium">
              Lihat Jadwal Lengkap →
            </NuxtLink>
          </div>
          
          <div v-if="nextImmunizations.length > 0" class="space-y-3">
            <!-- Next Immunization Cards -->
            <div 
              v-for="imm in nextImmunizations" 
              :key="imm.schedule.id"
              class="p-4 rounded-lg border-2 transition hover:shadow-md"
              :class="getNextImmunizationCardClass(imm)"
            >
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-2">
                    <h4 class="font-bold text-gray-900">
                      {{ imm.schedule.name_id || imm.schedule.name }}
                      <span v-if="imm.schedule.dose_number > 1" class="text-sm font-normal text-gray-600">
                        - Dosis {{ imm.schedule.dose_number }}
                      </span>
                    </h4>
                    <span class="px-2 py-0.5 rounded text-xs font-semibold"
                          :class="getStatusBadgeClass(imm.status)">
                      {{ getStatusLabel(imm.status) }}
                    </span>
                  </div>
                  
                  <div class="text-sm text-gray-600 space-y-1">
                    <p>
                      <span class="font-medium">Usia Optimal:</span> 
                      {{ imm.schedule.age_optimal_months ? `${imm.schedule.age_optimal_months} bulan` : '-' }}
                    </p>
                    <p v-if="imm.due_date" class="flex items-center gap-2 flex-wrap">
                      <span class="font-medium">Jadwal:</span> 
                      <span>{{ formatDate(imm.due_date) }}</span>
                      <span v-if="imm.days_until_due !== undefined && imm.days_until_due !== null && imm.days_until_due >= 0" 
                            class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs font-semibold">
                        {{ imm.days_until_due }} hari lagi
                      </span>
                      <span v-if="imm.days_overdue !== undefined && imm.days_overdue !== null && imm.days_overdue > 0" 
                            class="px-2 py-0.5 bg-red-100 text-red-700 rounded text-xs font-semibold">
                        Terlambat {{ imm.days_overdue }} hari
                      </span>
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Summary Footer -->
          <div class="mt-4 pt-4 border-t border-gray-200">
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">
                <span v-if="(immunizationStore.summary?.overdue || 0) > 0" class="text-red-600 font-semibold">
                  ⚠️ {{ immunizationStore.summary.overdue }} imunisasi terlambat
                </span>
                <span v-else-if="nextImmunizations.length > 0" class="text-gray-500">
                  Semua jadwal sesuai waktu
                </span>
                <span v-else class="text-green-600 font-semibold">
                  ✓ Semua imunisasi selesai
                </span>
              </span>
              <span class="text-gray-500">
                {{ immunizationStore.summary?.completed || 0 }} / {{ immunizationStore.summary?.total || 0 }} selesai
              </span>
            </div>
          </div>
          
          <!-- Empty State - All Completed -->
          <div v-if="nextImmunizations.length === 0" class="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg text-center">
            <p class="text-sm text-green-800 font-medium">
              ✅ Semua imunisasi untuk usia saat ini telah selesai
            </p>
            <p class="text-xs text-green-700 mt-1">
              Akan ada imunisasi baru ketika anak bertambah usia
            </p>
          </div>
        </div>

        <!-- Action Cards Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Aksi Cepat Section -->
          <div class="bg-white rounded-lg shadow-sm p-6">
            <h3 class="text-base font-bold text-gray-900 mb-4 flex items-center">
              <svg class="w-5 h-5 mr-2 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              Aksi Cepat
            </h3>
            <div class="space-y-3">
              <NuxtLink to="/growth/add" class="flex items-center space-x-3 p-3 border border-gray-200 rounded-lg hover:border-indigo-300 hover:bg-indigo-50 transition group">
                <div class="text-2xl group-hover:scale-110 transition-transform">📊</div>
                <div class="flex-1 min-w-0">
                  <p class="font-semibold text-gray-900 text-sm">Tambah Pengukuran</p>
                  <p class="text-xs text-gray-500">Input data berat & tinggi badan</p>
                </div>
                <svg class="w-5 h-5 text-gray-400 group-hover:text-indigo-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </NuxtLink>

              <NuxtLink to="/development" class="flex items-center space-x-3 p-3 border border-gray-200 rounded-lg hover:border-emerald-300 hover:bg-emerald-50 transition group">
                <div class="text-2xl group-hover:scale-110 transition-transform">✅</div>
                <div class="flex-1 min-w-0">
                  <p class="font-semibold text-gray-900 text-sm">Cek Milestone</p>
                  <p class="text-xs text-gray-500">Evaluasi perkembangan anak</p>
                </div>
                <svg class="w-5 h-5 text-gray-400 group-hover:text-emerald-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </NuxtLink>
            </div>
          </div>

          <!-- Lihat Grafik Section -->
          <div class="bg-white rounded-lg shadow-sm p-6">
            <h3 class="text-base font-bold text-gray-900 mb-4 flex items-center">
              <svg class="w-5 h-5 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              Lihat Grafik
            </h3>
            <div class="space-y-3">
              <NuxtLink 
                to="/growth/charts" 
                class="flex items-center space-x-3 p-3 border-2 border-indigo-200 rounded-lg hover:border-indigo-400 hover:bg-indigo-50 transition group"
                :class="{ 'opacity-50 cursor-not-allowed pointer-events-none': !measurementStore.hasMeasurements }"
              >
                <div class="text-2xl group-hover:scale-110 transition-transform">📈</div>
                <div class="flex-1 min-w-0">
                  <p class="font-semibold text-gray-900 text-sm">Grafik Pertumbuhan</p>
                  <p class="text-xs text-gray-500">BB/U dan TB/U berdasarkan standar WHO</p>
                  <p v-if="!measurementStore.hasMeasurements" class="text-xs text-amber-600 mt-0.5">Belum ada data</p>
                </div>
                <svg class="w-5 h-5 text-indigo-600 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </NuxtLink>

              <NuxtLink 
                to="/development/denver-ii" 
                class="flex items-center space-x-3 p-3 border-2 border-purple-200 rounded-lg hover:border-purple-400 hover:bg-purple-50 transition group"
              >
                <div class="text-2xl group-hover:scale-110 transition-transform">📉</div>
                <div class="flex-1 min-w-0">
                  <p class="font-semibold text-gray-900 text-sm">Grafik Denver II</p>
                  <p class="text-xs text-gray-500">4 domain perkembangan Denver II</p>
                </div>
                <svg class="w-5 h-5 text-purple-600 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </NuxtLink>

              <button 
                @click="downloadPDF"
                :disabled="downloadingPDF"
                class="flex items-center space-x-3 p-3 border-2 border-red-200 rounded-lg hover:border-red-400 hover:bg-red-50 transition group w-full disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <div class="text-2xl group-hover:scale-110 transition-transform">📄</div>
                <div class="flex-1 min-w-0 text-left">
                  <p class="font-semibold text-gray-900 text-sm">Export PDF</p>
                  <p class="text-xs text-gray-500">Unduh laporan lengkap untuk dokter</p>
                </div>
                <svg v-if="!downloadingPDF" class="w-5 h-5 text-red-600 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                <svg v-else class="animate-spin w-5 h-5 text-red-600" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </button>
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

const authStore = useAuthStore()
const childStore = useChildStore()
const measurementStore = useMeasurementStore()
const milestoneStore = useMilestoneStore()
const immunizationStore = useImmunizationStore()

const downloadingPDF = ref(false)

const downloadPDF = async () => {
  if (!childStore.selectedChild) {
    alert('Pilih anak terlebih dahulu')
    return
  }

  downloadingPDF.value = true
  try {
    const apiBase = useApiUrl()
    const childId = childStore.selectedChild.id
    
    const response = await fetch(`${apiBase}/api/children/${childId}/export-pdf`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })

    if (!response.ok) {
      throw new Error('Gagal mengunduh PDF')
    }

    // Get blob from response
    const blob = await response.blob()
    
    // Create download link
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `tukem_report_${childStore.selectedChild.name}_${new Date().toISOString().split('T')[0]}.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Error downloading PDF:', error)
    alert('Gagal mengunduh PDF. Silakan coba lagi.')
  } finally {
    downloadingPDF.value = false
  }
}

const { calculateCorrectedAge } = useCorrectedAge()

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

// Get next immunizations: yang terlewat, bulan ini, atau bulan depan
const nextImmunizations = computed(() => {
  // Get current age untuk filter berdasarkan periode waktu
  const currentAgeMonths = immunizationStore.ageMonths || 
    (ageInfo.value?.chronologicalDisplay ? 
      parseInt(ageInfo.value.chronologicalDisplay.replace(/[^\d]/g, '')) || 0 : 0)
  
  // Gabungkan semua yang belum selesai
  const all = [
    ...immunizationStore.overdueImmunizations, // Yang terlambat
    ...immunizationStore.upcomingImmunizations, // Yang akan datang
    ...immunizationStore.pendingImmunizations // Yang pending
  ]
  
  // Remove duplicates (same schedule.id)
  const seen = new Set()
  const unique = all.filter(item => {
    if (seen.has(item.schedule.id)) {
      return false
    }
    seen.add(item.schedule.id)
    return true
  })
  
  // Filter berdasarkan periode waktu yang relevan
  const filtered = unique.filter(item => {
    const ageOptimal = item.schedule.age_optimal_months || 0
    
    // 1. Tampilkan semua yang overdue (terlewat) - status overdue dari backend
    if (item.status === 'overdue') {
      return true
    }
    
    // Helper: cek apakah sudah terlewat waktunya berdasarkan due_date
    const today = new Date()
    today.setHours(0, 0, 0, 0) // Reset time untuk perbandingan tanggal saja
    
    // 2. Tampilkan yang sudah terlewat waktunya (due date sudah lewat tapi belum selesai)
    if (item.due_date) {
      const dueDate = new Date(item.due_date)
      dueDate.setHours(0, 0, 0, 0)
      
      // Jika due date sudah lewat (hari ini atau sebelumnya) dan belum selesai
      if (dueDate <= today && item.status !== 'completed') {
        return true
      }
    }
    
    // 3. Tampilkan yang usia optimalnya sudah terlewat (usia optimal < umur anak saat ini)
    // Contoh: Hepatitis B-0 usia optimal 0 bulan, jika anak sudah > 0 bulan = terlewat
    if (ageOptimal < currentAgeMonths && item.status !== 'completed') {
      return true
    }
    
    // 4. Tampilkan yang usia optimalnya di bulan ini (currentAgeMonths)
    if (ageOptimal === currentAgeMonths) {
      return true
    }
    
    // 5. Jika tidak ada yang untuk bulan ini (dan tidak terlewat), tampilkan yang untuk bulan depan
    // (cek apakah ada yang untuk bulan ini terlebih dulu yang tidak terlewat)
    
    const hasCurrentMonthNonPast = unique.some(i => {
      const isOverdue = i.status === 'overdue'
      const isPastDue = i.due_date && new Date(i.due_date) <= today && i.status !== 'completed'
      const ageOpt = i.schedule.age_optimal_months || 0
      const isPastAge = ageOpt < currentAgeMonths
      const isCurrentMonth = ageOpt === currentAgeMonths
      return isCurrentMonth && !isOverdue && !isPastDue && !isPastAge
    })
    
    if (!hasCurrentMonthNonPast && ageOptimal === currentAgeMonths + 1) {
      return true
    }
    
    return false
  })
  
  // Sort by: 
  // 1. Overdue / terlewat waktunya first (yang terlambat paling penting)
  // 2. Usia optimal (yang sudah lewat, kemudian bulan ini, kemudian bulan depan)
  // 3. Due date (earliest first)
  // 4. Priority (high first)
  return filtered.sort((a, b) => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    
    // Helper: cek apakah imunisasi sudah terlewat berdasarkan due_date
    const isAPastDue = a.due_date && new Date(a.due_date) <= today && a.status !== 'completed'
    const isBPastDue = b.due_date && new Date(b.due_date) <= today && b.status !== 'completed'
    
    // Helper: cek apakah usia optimal sudah terlewat
    const ageA = a.schedule.age_optimal_months || 0
    const ageB = b.schedule.age_optimal_months || 0
    const isAPastAge = ageA < currentAgeMonths && a.status !== 'completed'
    const isBPastAge = ageB < currentAgeMonths && b.status !== 'completed'
    
    // 1. Prioritaskan overdue atau yang sudah terlewat waktunya (due_date atau usia optimal)
    const isATerlewat = a.status === 'overdue' || isAPastDue || isAPastAge
    const isBTerlewat = b.status === 'overdue' || isBPastDue || isBPastAge
    
    if (isATerlewat && !isBTerlewat) return -1
    if (!isATerlewat && isBTerlewat) return 1
    
    // 2. Sort by usia optimal: yang sudah lewat → bulan ini → bulan depan
    // (ageA dan ageB sudah didefinisikan di atas, tidak perlu didefinisikan lagi)
    const isAPast = ageA < currentAgeMonths
    const isBPast = ageB < currentAgeMonths
    const isACurrent = ageA === currentAgeMonths
    const isBCurrent = ageB === currentAgeMonths
    const isANext = ageA === currentAgeMonths + 1
    const isBNext = ageB === currentAgeMonths + 1
    
    // Prioritas: sudah lewat → bulan ini → bulan depan
    if (isAPast && !isBPast) return -1
    if (!isAPast && isBPast) return 1
    
    if (isACurrent && !isBCurrent && !isBPast) return -1
    if (!isACurrent && !isAPast && isBCurrent) return 1
    
    if (isANext && !isBNext && !isBCurrent && !isBPast) return -1
    if (!isANext && !isACurrent && !isAPast && isBNext) return 1
    
    // Jika sudah lewat keduanya, prioritaskan yang paling lama lewatnya
    if (isAPast && isBPast) {
      return (currentAgeMonths - ageA) - (currentAgeMonths - ageB)
    }
    
    // 3. Jika periode sama, sort by due date (earliest first)
    if (a.due_date && b.due_date) {
      const dateA = new Date(a.due_date).getTime()
      const dateB = new Date(b.due_date).getTime()
      if (dateA !== dateB) {
        return dateA - dateB
      }
    } else if (a.due_date && !b.due_date) return -1
    else if (!a.due_date && b.due_date) return 1
    
    // 4. Jika semua sama, sort by priority (high first)
    const priorityOrder = { 'high': 0, 'medium': 1, 'low': 2 }
    return (priorityOrder[a.schedule.priority] || 1) - (priorityOrder[b.schedule.priority] || 1)
  }).slice(0, 3) // Maksimal 3 item di dashboard
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

const getNextImmunizationCardClass = (imm) => {
  if (imm.status === 'overdue') {
    return 'border-red-300 bg-red-50'
  }
  if (imm.status === 'upcoming' && imm.days_until_due !== undefined && imm.days_until_due <= 7) {
    return 'border-blue-300 bg-blue-50'
  }
  return 'border-yellow-300 bg-yellow-50'
}

// Compute age info for selected child
const ageInfo = computed(() => {
  if (!childStore.selectedChild?.dob) {
    return {
      chronologicalDisplay: '',
      correctedDisplay: null,
      useCorrected: false
    }
  }
  
  return calculateCorrectedAge(
    childStore.selectedChild.dob,
    childStore.selectedChild.is_premature || false,
    childStore.selectedChild.gestational_age,
    undefined // Use current date
  )
})

// Track if component is mounted to prevent updates after unmount
const isMounted = ref(false)
const isInitialized = ref(false)

// Watch for selected child changes to fetch measurements and milestones
// Don't use immediate: true to avoid triggering before component is ready
const stopWatcher = watch(() => childStore.selectedChild, async (newChild) => {
  // Guard: Don't update if component is unmounted or not initialized
  if (!isMounted.value || !isInitialized.value) return
  
  if (newChild) {
    try {
      await Promise.all([
        measurementStore.fetchLatestMeasurement(newChild.id),
        measurementStore.fetchMeasurements(newChild.id), // Also fetch all measurements to check hasMeasurements
        milestoneStore.fetchSummary(newChild.id),
        immunizationStore.fetchSchedule(newChild.id)
      ])
      
      // Guard again after async operation
      if (!isMounted.value) return
    } catch (error) {
      // Guard: Don't log if component is unmounted
      if (isMounted.value) {
        console.error('Error fetching dashboard data:', error)
      }
    }
  } else {
    // Guard: Don't update if component is unmounted
    if (isMounted.value) {
      measurementStore.latestMeasurement = null
      milestoneStore.summary = null
      immunizationStore.clearSchedule()
    }
  }
}, { immediate: false })

onUnmounted(() => {
  isMounted.value = false
  isInitialized.value = false
  stopWatcher()
})

onMounted(async () => {
  isMounted.value = true
  
  try {
    console.log('Dashboard mounted')
    console.log('Current auth token:', authStore.token ? 'exists' : 'missing')
    
    // Fetch children list (this will also restore selected child from localStorage)
    await childStore.fetchChildren()
    
    // Guard: Check if still mounted after async operation
    if (!isMounted.value) return
    
    console.log('Children fetched:', childStore.children.length, 'children')
    console.log('Selected child:', childStore.selectedChild?.id || 'none')
    
    // If there's a selected child, fetch their data
    if (childStore.selectedChild) {
      console.log('Fetching data for selected child:', childStore.selectedChild.id)
      try {
        await Promise.all([
          measurementStore.fetchLatestMeasurement(childStore.selectedChild.id),
          measurementStore.fetchMeasurements(childStore.selectedChild.id), // Also fetch all measurements to check hasMeasurements
          milestoneStore.fetchSummary(childStore.selectedChild.id),
          immunizationStore.fetchSchedule(childStore.selectedChild.id)
        ])
        
        // Guard: Check if still mounted after async operation
        if (!isMounted.value) return
        
        console.log('Latest measurement:', measurementStore.latestMeasurement)
        console.log('Milestone summary:', milestoneStore.summary)
      } catch (fetchError) {
        // Guard: Don't update if component is unmounted
        if (!isMounted.value) return
        
        console.error('Error fetching child data:', fetchError)
        // If fetch fails, it might be because child doesn't belong to user
        // Clear selection and try again
        childStore.clearState()
        await childStore.fetchChildren()
      }
    }
    
    // Mark as initialized after all initial data is loaded
    // This ensures watcher won't trigger until component is ready
    await nextTick()
    isInitialized.value = true
  } catch (error) {
    // Guard: Don't log if component is unmounted
    if (isMounted.value) {
      console.error('Error initializing dashboard:', error)
      // Still mark as initialized even on error to prevent hanging
      isInitialized.value = true
    }
  }
})
</script>
