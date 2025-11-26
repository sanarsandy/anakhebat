<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Tambah Pengukuran</h1>
      <p class="text-gray-600 mt-2">Input data pertumbuhan anak</p>
    </div>

    <!-- No Child Selected -->
    <div v-if="!childStore.selectedChild" class="bg-amber-50 border border-amber-200 rounded-xl p-8 text-center">
      <div class="text-5xl mb-4">⚠️</div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Pilih Anak Terlebih Dahulu</h2>
      <p class="text-gray-600 mb-4">Silakan pilih anak dari dropdown di header</p>
      <NuxtLink to="/children" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition">
        Kelola Profil Anak
      </NuxtLink>
    </div>

    <!-- Form -->
    <div v-else class="max-w-2xl bg-white rounded-xl shadow-sm p-8">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Measurement Date -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Tanggal Pengukuran *</label>
          <input 
            v-model="form.measurement_date"
            type="date" 
            required
            :max="today"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
          />
          <p v-if="calculatedAge" class="text-sm text-gray-600 mt-2">Usia saat pengukuran: {{ calculatedAge }}</p>
        </div>

        <!-- Weight -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Berat Badan (kg) *</label>
          <input 
            v-model.number="form.weight"
            type="number" 
            step="0.01"
            min="0"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
            placeholder="3.5"
          />
        </div>

        <!-- Height -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Tinggi Badan (cm) *</label>
          <input 
            v-model.number="form.height"
            type="number" 
            step="0.1"
            min="0"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
            placeholder="50.5"
          />
        </div>

        <!-- Head Circumference -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Lingkar Kepala (cm) <span class="text-gray-500">(Opsional)</span></label>
          <input 
            v-model.number="form.head_circumference"
            type="number" 
            step="0.1"
            min="0"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
            placeholder="35.0"
          />
        </div>

        <!-- Error Message -->
        <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
          <p class="text-red-600 text-sm">{{ error }}</p>
        </div>

        <!-- Actions -->
        <div class="flex space-x-4 pt-4">
          <button 
            type="button"
            @click="$router.back()"
            class="flex-1 px-6 py-3 bg-gray-100 text-gray-700 font-semibold rounded-lg hover:bg-gray-200 transition"
          >
            Batal
          </button>
          <button 
            type="submit"
            :disabled="loading"
            class="flex-1 px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition disabled:opacity-50"
          >
            {{ loading ? 'Menyimpan...' : 'Simpan' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const childStore = useChildStore()
const measurementStore = useMeasurementStore()
const router = useRouter()

const today = new Date().toISOString().split('T')[0]

const form = ref({
  measurement_date: today,
  weight: null,
  height: null,
  head_circumference: null
})

const loading = ref(false)
const error = ref('')
const calculatedAge = ref('')

// Calculate age when date changes
watch(() => form.value.measurement_date, (newDate) => {
  if (newDate && childStore.selectedChild) {
    const dob = new Date(childStore.selectedChild.dob)
    const measurementDate = new Date(newDate)
    
    const years = measurementDate.getFullYear() - dob.getFullYear()
    const months = measurementDate.getMonth() - dob.getMonth()
    const totalMonths = years * 12 + months
    
    const displayYears = Math.floor(totalMonths / 12)
    const displayMonths = totalMonths % 12
    
    if (displayYears > 0) {
      calculatedAge.value = `${displayYears} tahun ${displayMonths} bulan`
    } else {
      calculatedAge.value = `${totalMonths} bulan`
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  if (!childStore.selectedChild) {
    error.value = 'Silakan pilih anak terlebih dahulu'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const measurementData = {
      measurement_date: form.value.measurement_date,
      weight: form.value.weight,
      height: form.value.height
    }

    // Only include head circumference if provided
    if (form.value.head_circumference) {
      measurementData.head_circumference = form.value.head_circumference
    }

    const result = await measurementStore.addMeasurement(
      childStore.selectedChild.id,
      measurementData
    )

    if (result.success) {
      router.push('/growth')
    } else {
      error.value = result.error || 'Gagal menambahkan pengukuran'
      console.error('Submission failed:', result.error)
    }
  } catch (e) {
    console.error('Submission exception:', e)
    error.value = e.message || 'Terjadi kesalahan. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}

</script>
