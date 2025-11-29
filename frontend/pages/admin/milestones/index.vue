<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Milestones Management</h1>
        <p class="text-gray-600 mt-1">Kelola milestones untuk perkembangan anak</p>
      </div>
      <button
        @click="showCreateModal = true; editingMilestone = {}"
        class="px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors flex items-center gap-2"
      >
        <Icon name="mdi:plus" class="w-5 h-5" />
        <span>Tambah Milestone</span>
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input
            v-model="filters.search"
            type="text"
            placeholder="Cari milestone..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @input="debouncedSearch"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Source</label>
          <select
            v-model="filters.source"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchMilestones"
          >
            <option value="">All Sources</option>
            <option value="KPSP">KPSP</option>
            <option value="CDC">CDC</option>
            <option value="DENVER">DENVER</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select
            v-model="filters.category"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchMilestones"
          >
            <option value="">All Categories</option>
            <option value="sensory">Sensory</option>
            <option value="motor">Motor</option>
            <option value="perception">Perception</option>
            <option value="cognitive">Cognitive</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Red Flag</label>
          <select
            v-model="filters.red_flag"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchMilestones"
          >
            <option value="">All</option>
            <option value="true">Red Flag</option>
            <option value="false">Normal</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Milestones Table -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Age (Months)
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Question
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Category
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Source
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Level
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Red Flag
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="milestone in milestones" :key="milestone.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ milestone.age_months }}</div>
                <div v-if="milestone.min_age_range || milestone.max_age_range" class="text-xs text-gray-500">
                  {{ milestone.min_age_range }}-{{ milestone.max_age_range }} months
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ milestone.question }}</div>
                <div v-if="milestone.question_en" class="text-xs text-gray-500 mt-1">{{ milestone.question_en }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full" :class="getCategoryColor(milestone.category)">
                  {{ milestone.category }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ milestone.source }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                Level {{ milestone.pyramid_level }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="milestone.is_red_flag" class="px-2 py-1 text-xs font-medium rounded-full bg-red-100 text-red-800">
                  Red Flag
                </span>
                <span v-else class="text-sm text-gray-400">-</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="editMilestone(milestone)"
                  class="text-jurnal-teal-600 hover:text-jurnal-teal-900 mr-4"
                >
                  Edit
                </button>
                <button
                  @click="deleteMilestone(milestone)"
                  class="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total > 0" class="bg-gray-50 px-6 py-4 border-t border-gray-200 flex items-center justify-between">
        <div class="text-sm text-gray-700">
          Menampilkan {{ (pagination.page - 1) * pagination.limit + 1 }} sampai
          {{ Math.min(pagination.page * pagination.limit, pagination.total) }} dari {{ pagination.total }} milestones
        </div>
        <div class="flex gap-2">
          <button
            @click="changePage(pagination.page - 1)"
            :disabled="pagination.page === 1"
            class="px-3 py-1 border border-gray-300 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100"
          >
            Previous
          </button>
          <button
            @click="changePage(pagination.page + 1)"
            :disabled="pagination.page * pagination.limit >= pagination.total"
            class="px-3 py-1 border border-gray-300 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100"
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showCreateModal = false">
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="p-6 border-b border-gray-200">
          <h2 class="text-2xl font-bold text-gray-900">
            {{ editingMilestone.id ? 'Edit Milestone' : 'Tambah Milestone' }}
          </h2>
        </div>
        <form @submit.prevent="saveMilestone" class="p-6 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age (Months) *</label>
              <input
                v-model.number="editingMilestone.age_months"
                type="number"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Pyramid Level *</label>
              <select
                v-model.number="editingMilestone.pyramid_level"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="1">Level 1 - Sensory</option>
                <option value="2">Level 2 - Motor</option>
                <option value="3">Level 3 - Perception</option>
                <option value="4">Level 4 - Cognitive</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Min Age Range</label>
              <input
                v-model.number="editingMilestone.min_age_range"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Max Age Range</label>
              <input
                v-model.number="editingMilestone.max_age_range"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Category *</label>
              <select
                v-model="editingMilestone.category"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="sensory">Sensory</option>
                <option value="motor">Motor</option>
                <option value="perception">Perception</option>
                <option value="cognitive">Cognitive</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Source *</label>
              <select
                v-model="editingMilestone.source"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="KPSP">KPSP</option>
                <option value="CDC">CDC</option>
                <option value="DENVER">DENVER</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Denver Domain</label>
            <select
              v-model="editingMilestone.denver_domain"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            >
              <option :value="null">None</option>
              <option value="PS">PS - Personal-Social</option>
              <option value="FM">FM - Fine Motor-Adaptive</option>
              <option value="L">L - Language</option>
              <option value="GM">GM - Gross Motor</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Question (ID) *</label>
            <textarea
              v-model="editingMilestone.question"
              required
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Question (EN)</label>
            <textarea
              v-model="editingMilestone.question_en"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="flex items-center gap-2">
              <input
                v-model="editingMilestone.is_red_flag"
                type="checkbox"
                class="w-4 h-4 text-jurnal-teal-600 border-gray-300 rounded focus:ring-jurnal-teal-500"
              />
              <span class="text-sm font-medium text-gray-700">Red Flag</span>
            </label>
          </div>

          <div class="flex justify-end gap-3 pt-4 border-t border-gray-200">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 disabled:opacity-50"
            >
              {{ saving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const apiBase = useApiUrl()
const authStore = useAuthStore()

const milestones = ref([])
const pagination = ref({ page: 1, limit: 50, total: 0 })
const filters = ref({
  search: '',
  source: '',
  category: '',
  red_flag: ''
})
const showCreateModal = ref(false)
const editingMilestone = ref({})
const saving = ref(false)

let searchTimeout: NodeJS.Timeout | null = null

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchMilestones()
  }, 500)
}

const fetchMilestones = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.search) params.append('search', filters.value.search)
    if (filters.value.source) params.append('source', filters.value.source)
    if (filters.value.category) params.append('category', filters.value.category)
    if (filters.value.red_flag) params.append('red_flag', filters.value.red_flag)

    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/milestones?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) throw new Error('Failed to fetch milestones')

    const data = await response.json()
    milestones.value = data.milestones || []
    pagination.value = data.pagination || pagination.value
  } catch (error) {
    console.error('Error fetching milestones:', error)
    alert('Gagal memuat milestones')
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchMilestones()
}

const editMilestone = (milestone: any) => {
  editingMilestone.value = { ...milestone }
  showCreateModal.value = true
}

const saveMilestone = async () => {
  saving.value = true
  try {
    const url = editingMilestone.value.id
      ? `${apiBase}/api/admin/milestones/${editingMilestone.value.id}`
      : `${apiBase}/api/admin/milestones`

    const method = editingMilestone.value.id ? 'PUT' : 'POST'

    const payload: any = {
      age_months: editingMilestone.value.age_months,
      pyramid_level: editingMilestone.value.pyramid_level,
      category: editingMilestone.value.category,
      question: editingMilestone.value.question,
      source: editingMilestone.value.source,
      is_red_flag: editingMilestone.value.is_red_flag || false,
    }

    if (editingMilestone.value.min_age_range) payload.min_age_range = editingMilestone.value.min_age_range
    if (editingMilestone.value.max_age_range) payload.max_age_range = editingMilestone.value.max_age_range
    if (editingMilestone.value.question_en) payload.question_en = editingMilestone.value.question_en
    if (editingMilestone.value.denver_domain) payload.denver_domain = editingMilestone.value.denver_domain

    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to save milestone')
    }

    showCreateModal.value = false
    editingMilestone.value = {}
    fetchMilestones()
  } catch (error: any) {
    console.error('Error saving milestone:', error)
    alert(error.message || 'Gagal menyimpan milestone')
  } finally {
    saving.value = false
  }
}

const deleteMilestone = async (milestone: any) => {
  if (!confirm(`Yakin ingin menghapus milestone ini?`)) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/milestones/${milestone.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to delete milestone')
    }

    fetchMilestones()
  } catch (error: any) {
    console.error('Error deleting milestone:', error)
    alert(error.message || 'Gagal menghapus milestone')
  }
}

const getCategoryColor = (category: string) => {
  const colors: Record<string, string> = {
    sensory: 'bg-blue-100 text-blue-800',
    motor: 'bg-green-100 text-green-800',
    perception: 'bg-yellow-100 text-yellow-800',
    cognitive: 'bg-purple-100 text-purple-800'
  }
  return colors[category] || 'bg-gray-100 text-gray-800'
}

onMounted(() => {
  fetchMilestones()
})
</script>

