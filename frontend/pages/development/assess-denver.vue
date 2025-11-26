<template>
  <div class="p-6 max-w-3xl mx-auto pb-32 lg:pb-24">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Checklist Denver II</h1>
        <p class="text-gray-600">Usia: {{ currentAge }} Bulan</p>
      </div>
      <div class="text-sm text-gray-500">
        <span v-if="saving" class="flex items-center text-jurnal-teal-600">
          <Icon name="mdi:loading" class="animate-spin mr-1" /> Menyimpan...
        </span>
        <span v-else-if="lastSaved">Disimpan {{ lastSaved }}</span>
      </div>
    </div>

    <!-- Info Box -->
    <div class="bg-blue-50 border border-blue-200 rounded-xl p-4 mb-6">
      <p class="text-sm text-blue-800">
        <strong>Denver II</strong> adalah alat skrining perkembangan yang mengevaluasi 4 domain: 
        <strong>Personal-Social (PS)</strong>, <strong>Fine Motor-Adaptive (FM)</strong>, 
        <strong>Language (L)</strong>, dan <strong>Gross Motor (GM)</strong>.
      </p>
    </div>

    <div v-if="milestoneStore.loading" class="py-12 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
      <p class="mt-4 text-gray-500">Memuat checklist Denver II...</p>
    </div>

    <div v-else-if="Object.keys(milestoneStore.milestonesByDenverDomain).length === 0" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="text-6xl mb-4">📋</div>
      <h3 class="text-xl font-bold text-gray-900 mb-2">Tidak Ada Milestone</h3>
      <p class="text-gray-600 mb-6">Tidak ada milestone Denver II yang tersedia untuk usia {{ currentAge }} bulan.</p>
      <NuxtLink to="/development" class="inline-block px-6 py-3 bg-jurnal-teal-600 text-white font-bold rounded-lg hover:bg-jurnal-teal-700">
        Kembali ke Dashboard
      </NuxtLink>
    </div>

    <div v-else class="space-y-8">
      <!-- Denver II Domains -->
      <div v-for="(group, domain) in milestoneStore.milestonesByDenverDomain" :key="domain" class="bg-white rounded-xl shadow-sm overflow-hidden">
        <div class="px-6 py-4 bg-gradient-to-r from-jurnal-teal-500 to-purple-600 text-white flex items-center justify-between">
          <h3 class="font-bold text-lg flex items-center">
            <span class="w-8 h-8 rounded-full bg-white bg-opacity-20 text-white flex items-center justify-center text-sm font-bold mr-3">{{ domain }}</span>
            {{ group.name }}
          </h3>
          <span class="text-sm font-medium px-3 py-1 rounded-full bg-white bg-opacity-20">{{ group.items.length }} Item</span>
        </div>

        <div class="divide-y divide-gray-100">
          <div v-for="item in group.items" :key="item.id" class="p-6 hover:bg-gray-50 transition">
            <div class="flex items-start justify-between mb-4">
              <div class="pr-4 flex-1">
                <p class="text-gray-900 font-medium">{{ item.question }}</p>
                <p v-if="item.question_en" class="text-sm text-gray-500 mt-1 italic">{{ item.question_en }}</p>
                <div class="flex items-center mt-2 space-x-2">
                  <span v-if="item.is_red_flag" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                    <Icon name="mdi:alert-circle" class="mr-1" /> Tanda Bahaya
                  </span>
                  <span class="text-xs text-gray-500">Usia target: {{ item.age_months }} bulan</span>
                </div>
              </div>
            </div>

            <!-- Options -->
            <div class="flex space-x-4">
              <label class="flex-1 cursor-pointer group">
                <input type="radio" :name="item.id" value="yes" 
                       :checked="getDraftStatus(item.id) === 'yes'"
                       @change="updateStatus(item.id, 'yes')"
                       class="peer sr-only">
                <div class="text-center py-3 rounded-xl border-2 border-gray-100 peer-checked:bg-green-50 peer-checked:border-green-500 peer-checked:text-green-700 group-hover:border-green-200 transition-all duration-200">
                  <Icon name="mdi:check-circle" class="w-6 h-6 mx-auto mb-1 text-gray-300 peer-checked:text-green-500 group-hover:text-green-400 transition-colors" />
                  <div class="font-bold text-sm">Ya / Bisa</div>
                </div>
              </label>

              <label class="flex-1 cursor-pointer group">
                <input type="radio" :name="item.id" value="sometimes" 
                       :checked="getDraftStatus(item.id) === 'sometimes'"
                       @change="updateStatus(item.id, 'sometimes')"
                       class="peer sr-only">
                <div class="text-center py-3 rounded-xl border-2 border-gray-100 peer-checked:bg-yellow-50 peer-checked:border-yellow-500 peer-checked:text-yellow-700 group-hover:border-yellow-200 transition-all duration-200">
                  <Icon name="mdi:minus-circle" class="w-6 h-6 mx-auto mb-1 text-gray-300 peer-checked:text-yellow-500 group-hover:text-yellow-400 transition-colors" />
                  <div class="font-bold text-sm">Kadang</div>
                </div>
              </label>

              <label class="flex-1 cursor-pointer group">
                <input type="radio" :name="item.id" value="no" 
                       :checked="getDraftStatus(item.id) === 'no'"
                       @change="updateStatus(item.id, 'no')"
                       class="peer sr-only">
                <div class="text-center py-3 rounded-xl border-2 border-gray-100 peer-checked:bg-red-50 peer-checked:border-red-500 peer-checked:text-red-700 group-hover:border-red-200 transition-all duration-200">
                  <Icon name="mdi:close-circle" class="w-6 h-6 mx-auto mb-1 text-gray-300 peer-checked:text-red-500 group-hover:text-red-400 transition-colors" />
                  <div class="font-bold text-sm">Belum</div>
                </div>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Action Bar - Positioned above BottomNav on mobile -->
    <div class="fixed bottom-16 left-0 right-0 bg-white border-t border-gray-200 shadow-lg z-[60] lg:bottom-0 lg:pl-64">
      <div class="max-w-3xl mx-auto px-4 py-3 flex items-center justify-between">
        <div class="text-sm text-gray-600 hidden sm:block">
          {{ draftCount }} item terjawab
        </div>
        <div class="flex space-x-3 w-full sm:w-auto">
          <button @click="$router.back()" class="px-4 py-2.5 text-gray-600 font-medium hover:bg-gray-100 rounded-lg flex-1 sm:flex-none transition">
            Batal
          </button>
          <button 
            @click="saveAll" 
            :disabled="saving || draftCount === 0"
            class="px-6 py-2.5 bg-jurnal-teal-600 text-white font-bold rounded-lg hover:bg-jurnal-teal-700 disabled:opacity-50 disabled:cursor-not-allowed shadow-md flex-1 sm:flex-none transition"
          >
            {{ saving ? 'Menyimpan...' : 'Simpan' }}
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
const milestoneStore = useMilestoneStore()
const router = useRouter()

const saving = ref(false)
const lastSaved = ref('')
const currentAge = ref(0)

const draftCount = computed(() => {
  if (!childStore.selectedChild) return 0
  const drafts = milestoneStore.draftAssessments[childStore.selectedChild.id]
  return drafts ? drafts.length : 0
})

onMounted(async () => {
  milestoneStore.loadDraftsFromStorage()
  
  if (childStore.selectedChild) {
    // Calculate age in months
    const dob = new Date(childStore.selectedChild.dob)
    const now = new Date()
    const diffMonths = (now.getFullYear() - dob.getFullYear()) * 12 + (now.getMonth() - dob.getMonth())
    currentAge.value = diffMonths
    
    // Fetch Denver II milestones for this age
    await milestoneStore.fetchDenverIIMilestones(diffMonths)
  }
})

const getDraftStatus = (milestoneId) => {
  if (!childStore.selectedChild) return null
  const drafts = milestoneStore.draftAssessments[childStore.selectedChild.id] || []
  const draft = drafts.find(d => d.milestone_id === milestoneId)
  return draft ? draft.status : null
}

const updateStatus = (milestoneId, status) => {
  if (!childStore.selectedChild) return
  
  milestoneStore.saveDraft(childStore.selectedChild.id, milestoneId, status)
  lastSaved.value = 'Draft (Local)'
}

const saveAll = async () => {
  if (!childStore.selectedChild) {
    alert('Silakan pilih anak terlebih dahulu')
    return
  }
  
  if (draftCount.value === 0) {
    alert('Tidak ada data penilaian untuk disimpan. Silakan isi checklist terlebih dahulu.')
    return
  }
  
  saving.value = true
  const today = new Date().toISOString().split('T')[0]
  
  try {
    const result = await milestoneStore.syncAssessments(childStore.selectedChild.id, today)
    
    if (result.success) {
      lastSaved.value = 'Tersimpan'
      // Wait a bit before redirecting to show success message
      setTimeout(() => {
        router.push('/development/denver-ii')
      }, 500)
    } else {
      alert('Gagal menyimpan penilaian: ' + (result.error || 'Terjadi kesalahan'))
    }
  } catch (error) {
    console.error('Error saving assessments:', error)
    alert('Gagal menyimpan penilaian: ' + (error.message || 'Terjadi kesalahan'))
  } finally {
    saving.value = false
  }
}
</script>

