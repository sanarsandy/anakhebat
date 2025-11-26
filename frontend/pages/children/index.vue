<template>
  <div>
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Profil Anak</h1>
        <p class="text-gray-600 mt-2">Kelola data anak Anda</p>
      </div>
      <NuxtLink to="/children/add" class="px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition flex items-center space-x-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span>Tambah Anak</span>
      </NuxtLink>
    </div>

    <!-- Empty State -->
    <div v-if="!childStore.hasChildren" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="text-6xl mb-4">👶</div>
      <h2 class="text-2xl font-bold text-gray-900 mb-2">Belum Ada Data Anak</h2>
      <p class="text-gray-600 mb-6">Mulai dengan menambahkan profil anak Anda</p>
      <NuxtLink to="/children/add" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-semibold rounded-lg hover:bg-jurnal-teal-700 transition">
        Tambah Anak Pertama
      </NuxtLink>
    </div>

    <!-- Children Grid -->
    <div v-else class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="child in childStore.children" 
        :key="child.id"
        class="bg-white rounded-xl shadow-sm hover:shadow-lg transition-all p-6 border-2"
        :class="childStore.selectedChild?.id === child.id ? 'border-jurnal-teal-500' : 'border-gray-200'"
      >
        <div class="flex items-start justify-between mb-4">
          <div class="text-5xl">{{ child.gender === 'male' ? '👦' : '👧' }}</div>
          <button 
            @click="selectChild(child.id)"
            class="px-3 py-1 text-sm rounded-lg transition"
            :class="childStore.selectedChild?.id === child.id ? 'bg-jurnal-teal-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
          >
            {{ childStore.selectedChild?.id === child.id ? 'Terpilih' : 'Pilih' }}
          </button>
        </div>

        <h3 class="text-xl font-bold text-gray-900 mb-2">{{ child.name }}</h3>
        <div class="space-y-1 text-sm text-gray-600 mb-4">
          <p>{{ child.gender === 'male' ? 'Laki-laki' : 'Perempuan' }}</p>
          <div class="space-y-0.5">
            <p>{{ getAgeDisplay(child) }}</p>
            <div v-if="getAgeInfo(child).useCorrected" class="flex items-center gap-1">
              <span class="px-1.5 py-0.5 bg-amber-100 text-amber-700 text-xs rounded">Usia Koreksi: {{ getAgeInfo(child).correctedDisplay }}</span>
            </div>
          </div>
          <p class="text-xs text-gray-500">Lahir: {{ formatDate(child.dob) }}</p>
          <p v-if="child.is_premature" class="text-xs text-orange-600 font-medium">
            Prematur{{ child.gestational_age ? ` (${child.gestational_age} minggu)` : '' }}
          </p>
        </div>

        <div class="flex space-x-2">
          <NuxtLink 
            :to="`/children/${child.id}/edit`"
            class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 text-center rounded-lg hover:bg-gray-200 transition text-sm font-medium"
          >
            Edit
          </NuxtLink>
          <button 
            @click="confirmDelete(child)"
            class="flex-1 px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition text-sm font-medium"
          >
            Hapus
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold text-gray-900 mb-2">Hapus Profil Anak?</h3>
        <p class="text-gray-600 mb-6">Apakah Anda yakin ingin menghapus profil <strong>{{ childToDelete?.name }}</strong>? Semua data terkait akan dihapus.</p>
        <div class="flex space-x-3">
          <button 
            @click="showDeleteModal = false"
            class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium"
          >
            Batal
          </button>
          <button 
            @click="deleteChild"
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
definePageMeta({
  middleware: 'auth'
})

const childStore = useChildStore()
const showDeleteModal = ref(false)
const childToDelete = ref(null)

const selectChild = (childId) => {
  childStore.selectChild(childId)
}

const confirmDelete = (child) => {
  childToDelete.value = child
  showDeleteModal.value = true
}

const deleteChild = async () => {
  if (childToDelete.value) {
    await childStore.deleteChild(childToDelete.value.id)
    showDeleteModal.value = false
    childToDelete.value = null
  }
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
  return info.chronologicalDisplay
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  await childStore.fetchChildren()
})
</script>
