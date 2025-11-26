<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Tambah Anak Baru</h1>
      <p class="text-gray-600 mt-2">Lengkapi data profil anak Anda</p>
    </div>

    <div class="max-w-2xl bg-white rounded-xl shadow-sm p-8">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Name -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Nama Lengkap *</label>
          <input 
            v-model="form.name"
            type="text" 
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
            placeholder="Masukkan nama anak"
          />
        </div>

        <!-- Date of Birth -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Tanggal Lahir *</label>
          <input 
            v-model="form.dob"
            type="date" 
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
          />
        </div>

        <!-- Gender -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Jenis Kelamin *</label>
          <div class="grid grid-cols-2 gap-4">
            <label class="flex items-center justify-center p-4 border-2 rounded-lg cursor-pointer transition"
              :class="form.gender === 'male' ? 'border-jurnal-teal-500 bg-jurnal-teal-50' : 'border-gray-300 hover:border-gray-400'">
              <input type="radio" v-model="form.gender" value="male" class="sr-only" required />
              <span class="text-4xl mr-3">👦</span>
              <span class="font-medium">Laki-laki</span>
            </label>
            <label class="flex items-center justify-center p-4 border-2 rounded-lg cursor-pointer transition"
              :class="form.gender === 'female' ? 'border-jurnal-teal-500 bg-jurnal-teal-50' : 'border-gray-300 hover:border-gray-400'">
              <input type="radio" v-model="form.gender" value="female" class="sr-only" required />
              <span class="text-4xl mr-3">👧</span>
              <span class="font-medium">Perempuan</span>
            </label>
          </div>
        </div>

        <!-- Birth Weight & Height -->
        <div class="grid md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Berat Lahir (kg) *</label>
            <input 
              v-model.number="form.birth_weight"
              type="number" 
              step="0.01"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              placeholder="3.2"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Panjang Lahir (cm) *</label>
            <input 
              v-model.number="form.birth_height"
              type="number" 
              step="0.1"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-jurnal-teal-500"
              placeholder="50"
            />
          </div>
        </div>

        <!-- Premature Checkbox -->
        <div class="bg-amber-50 border border-amber-200 rounded-lg p-4">
          <label class="flex items-start cursor-pointer">
            <input 
              v-model="form.is_premature"
              type="checkbox" 
              class="mt-1 mr-3 h-5 w-5 text-jurnal-teal-600 border-gray-300 rounded focus:ring-jurnal-teal-500"
            />
            <div>
              <span class="font-medium text-gray-900">Anak Lahir Prematur</span>
              <p class="text-sm text-gray-600 mt-1">Centang jika anak lahir sebelum usia kehamilan 37 minggu</p>
            </div>
          </label>
        </div>

        <!-- Gestational Age (if premature) -->
        <div v-if="form.is_premature" class="bg-amber-50 border border-amber-200 rounded-lg p-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Usia Kehamilan Saat Lahir (minggu) *</label>
          <input 
            v-model.number="form.gestational_age"
            type="number" 
            min="20"
            max="42"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
            placeholder="34"
          />
          <p class="text-xs text-gray-600 mt-2">Usia kehamilan normal: 37-42 minggu. Prematur: &lt; 37 minggu</p>
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
const router = useRouter()

const form = ref({
  name: '',
  dob: '',
  gender: '',
  birth_weight: null,
  birth_height: null,
  is_premature: false,
  gestational_age: null
})

const loading = ref(false)
const error = ref('')

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    const result = await childStore.addChild(form.value)
    
    if (result.success) {
      router.push('/children')
    } else {
      error.value = result.error || 'Gagal menambahkan data anak. Silakan coba lagi.'
    }
  } catch (e) {
    console.error('Submit error:', e)
    error.value = e.message || 'Terjadi kesalahan. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}
</script>
