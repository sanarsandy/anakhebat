<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Immunization Schedules Management</h1>
        <p class="text-gray-600 mt-1">Kelola jadwal imunisasi berdasarkan rekomendasi IDAI</p>
      </div>
      <button
        @click="showCreateModal = true; editingSchedule = { category: 'wajib', priority: 'medium', is_required: true, is_active: true, source: 'IDAI' }"
        class="px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors flex items-center gap-2"
      >
        <Icon name="mdi:plus" class="w-5 h-5" />
        <span>Tambah Jadwal</span>
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
            placeholder="Cari jadwal..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @input="debouncedSearch"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select
            v-model="filters.category"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchSchedules"
          >
            <option value="">All Categories</option>
            <option value="wajib">Wajib</option>
            <option value="tambahan">Tambahan</option>
            <option value="catch-up">Catch-up</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
          <select
            v-model="filters.priority"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchSchedules"
          >
            <option value="">All Priorities</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Required</label>
          <select
            v-model="filters.is_required"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchSchedules"
          >
            <option value="">All</option>
            <option value="true">Required</option>
            <option value="false">Optional</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Schedules Table -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Dose
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Age Range
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Category
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Priority
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="schedule in schedules" :key="schedule.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="text-sm font-medium text-gray-900">{{ schedule.name }}</div>
                <div v-if="schedule.name_id" class="text-xs text-gray-500 mt-1">{{ schedule.name_id }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                Dose {{ schedule.dose_number }}
                <span v-if="schedule.total_doses">/ {{ schedule.total_doses }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <div v-if="schedule.age_optimal_days">
                  {{ Math.floor((schedule.age_optimal_days || 0) / 30) }} months
                </div>
                <div v-else-if="schedule.age_optimal_months">
                  {{ schedule.age_optimal_months }} months
                </div>
                <div v-else>-</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full" :class="getCategoryColor(schedule.category)">
                  {{ schedule.category }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full" :class="getPriorityColor(schedule.priority)">
                  {{ schedule.priority }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="editSchedule(schedule)"
                  class="text-jurnal-teal-600 hover:text-jurnal-teal-900 mr-4"
                >
                  Edit
                </button>
                <button
                  @click="deleteSchedule(schedule)"
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
          {{ Math.min(pagination.page * pagination.limit, pagination.total) }} dari {{ pagination.total }} schedules
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
      <div class="bg-white rounded-xl shadow-xl max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="p-6 border-b border-gray-200">
          <h2 class="text-2xl font-bold text-gray-900">
            {{ editingSchedule.id ? 'Edit Immunization Schedule' : 'Tambah Immunization Schedule' }}
          </h2>
        </div>
        <form @submit.prevent="saveSchedule" class="p-6 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Name (EN) *</label>
              <input
                v-model="editingSchedule.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Name (ID)</label>
              <input
                v-model="editingSchedule.name_id"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              v-model="editingSchedule.description"
              rows="2"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            />
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age Min (Days)</label>
              <input
                v-model.number="editingSchedule.age_min_days"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age Optimal (Days)</label>
              <input
                v-model.number="editingSchedule.age_optimal_days"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age Max (Days)</label>
              <input
                v-model.number="editingSchedule.age_max_days"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age Min (Months)</label>
              <input
                v-model.number="editingSchedule.age_min_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age Optimal (Months)</label>
              <input
                v-model.number="editingSchedule.age_optimal_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age Max (Months)</label>
              <input
                v-model.number="editingSchedule.age_max_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Dose Number *</label>
              <input
                v-model.number="editingSchedule.dose_number"
                type="number"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Total Doses</label>
              <input
                v-model.number="editingSchedule.total_doses"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Interval (Days)</label>
              <input
                v-model.number="editingSchedule.interval_from_previous_days"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Interval (Months)</label>
              <input
                v-model.number="editingSchedule.interval_from_previous_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <select
                v-model="editingSchedule.category"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="wajib">Wajib</option>
                <option value="tambahan">Tambahan</option>
                <option value="catch-up">Catch-up</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
              <select
                v-model="editingSchedule.priority"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Source</label>
              <input
                v-model="editingSchedule.source"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
                placeholder="IDAI"
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Notes</label>
            <textarea
              v-model="editingSchedule.notes"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="flex items-center gap-2">
                <input
                  v-model="editingSchedule.is_required"
                  type="checkbox"
                  class="w-4 h-4 text-jurnal-teal-600 border-gray-300 rounded focus:ring-jurnal-teal-500"
                />
                <span class="text-sm font-medium text-gray-700">Required</span>
              </label>
            </div>
            <div>
              <label class="flex items-center gap-2">
                <input
                  v-model="editingSchedule.is_active"
                  type="checkbox"
                  class="w-4 h-4 text-jurnal-teal-600 border-gray-300 rounded focus:ring-jurnal-teal-500"
                />
                <span class="text-sm font-medium text-gray-700">Active</span>
              </label>
            </div>
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

const schedules = ref([])
const pagination = ref({ page: 1, limit: 50, total: 0 })
const filters = ref({
  search: '',
  category: '',
  priority: '',
  is_required: ''
})
const showCreateModal = ref(false)
const editingSchedule = ref({})
const saving = ref(false)

let searchTimeout: NodeJS.Timeout | null = null

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchSchedules()
  }, 500)
}

const fetchSchedules = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.search) params.append('search', filters.value.search)
    if (filters.value.category) params.append('category', filters.value.category)
    if (filters.value.priority) params.append('priority', filters.value.priority)
    if (filters.value.is_required) params.append('is_required', filters.value.is_required)

    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/immunization-schedules?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) throw new Error('Failed to fetch immunization schedules')

    const data = await response.json()
    schedules.value = data.schedules || []
    pagination.value = data.pagination || pagination.value
  } catch (error) {
    console.error('Error fetching immunization schedules:', error)
    alert('Gagal memuat immunization schedules')
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchSchedules()
}

const editSchedule = (schedule: any) => {
  editingSchedule.value = { ...schedule }
  showCreateModal.value = true
}

const saveSchedule = async () => {
  saving.value = true
  try {
    const url = editingSchedule.value.id
      ? `${apiBase}/api/admin/immunization-schedules/${editingSchedule.value.id}`
      : `${apiBase}/api/admin/immunization-schedules`

    const method = editingSchedule.value.id ? 'PUT' : 'POST'

    const payload: any = {
      name: editingSchedule.value.name,
      dose_number: parseInt(editingSchedule.value.dose_number),
      category: editingSchedule.value.category || 'wajib',
      priority: editingSchedule.value.priority || 'medium',
      is_required: editingSchedule.value.is_required !== false,
      is_active: editingSchedule.value.is_active !== false,
      source: editingSchedule.value.source || 'IDAI',
    }

    if (editingSchedule.value.name_id) payload.name_id = editingSchedule.value.name_id
    if (editingSchedule.value.description) payload.description = editingSchedule.value.description
    if (editingSchedule.value.age_min_days) payload.age_min_days = parseInt(editingSchedule.value.age_min_days)
    if (editingSchedule.value.age_optimal_days) payload.age_optimal_days = parseInt(editingSchedule.value.age_optimal_days)
    if (editingSchedule.value.age_max_days) payload.age_max_days = parseInt(editingSchedule.value.age_max_days)
    if (editingSchedule.value.age_min_months) payload.age_min_months = parseInt(editingSchedule.value.age_min_months)
    if (editingSchedule.value.age_optimal_months) payload.age_optimal_months = parseInt(editingSchedule.value.age_optimal_months)
    if (editingSchedule.value.age_max_months) payload.age_max_months = parseInt(editingSchedule.value.age_max_months)
    if (editingSchedule.value.total_doses) payload.total_doses = parseInt(editingSchedule.value.total_doses)
    if (editingSchedule.value.interval_from_previous_days) payload.interval_from_previous_days = parseInt(editingSchedule.value.interval_from_previous_days)
    if (editingSchedule.value.interval_from_previous_months) payload.interval_from_previous_months = parseInt(editingSchedule.value.interval_from_previous_months)
    if (editingSchedule.value.notes) payload.notes = editingSchedule.value.notes

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
      throw new Error(error.error || 'Failed to save immunization schedule')
    }

    showCreateModal.value = false
    editingSchedule.value = {}
    fetchSchedules()
  } catch (error: any) {
    console.error('Error saving immunization schedule:', error)
    alert(error.message || 'Gagal menyimpan immunization schedule')
  } finally {
    saving.value = false
  }
}

const deleteSchedule = async (schedule: any) => {
  if (!confirm(`Yakin ingin menghapus jadwal imunisasi ini?`)) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/immunization-schedules/${schedule.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to delete immunization schedule')
    }

    fetchSchedules()
  } catch (error: any) {
    console.error('Error deleting immunization schedule:', error)
    alert(error.message || 'Gagal menghapus immunization schedule')
  }
}

const getCategoryColor = (category: string) => {
  const colors: Record<string, string> = {
    wajib: 'bg-red-100 text-red-800',
    tambahan: 'bg-blue-100 text-blue-800',
    'catch-up': 'bg-yellow-100 text-yellow-800'
  }
  return colors[category] || 'bg-gray-100 text-gray-800'
}

const getPriorityColor = (priority: string) => {
  const colors: Record<string, string> = {
    high: 'bg-red-100 text-red-800',
    medium: 'bg-yellow-100 text-yellow-800',
    low: 'bg-green-100 text-green-800'
  }
  return colors[priority] || 'bg-gray-100 text-gray-800'
}

onMounted(() => {
  fetchSchedules()
})
</script>

