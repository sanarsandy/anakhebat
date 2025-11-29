<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">WHO Standards Management</h1>
        <p class="text-gray-600 mt-1">Kelola WHO Child Growth Standards LMS parameters</p>
      </div>
      <button
        @click="showCreateModal = true; editingStandard = {}"
        class="px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors flex items-center gap-2"
      >
        <Icon name="mdi:plus" class="w-5 h-5" />
        <span>Tambah WHO Standard</span>
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Indicator</label>
          <select
            v-model="filters.indicator"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchStandards"
          >
            <option value="">All Indicators</option>
            <option value="wfa">WFA - Weight for Age</option>
            <option value="hfa">HFA - Height for Age</option>
            <option value="wfh">WFH - Weight for Height</option>
            <option value="hcfa">HCFA - Head Circumference for Age</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Gender</label>
          <select
            v-model="filters.gender"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchStandards"
          >
            <option value="">All Genders</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Standards Table -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Indicator
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Gender
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Age/Height
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                LMS Values
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="standard in standards" :key="standard.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm font-medium text-gray-900 uppercase">{{ standard.indicator }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ standard.gender }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <div v-if="standard.age_months !== null && standard.age_months !== undefined">
                  {{ standard.age_months }} months
                </div>
                <div v-else-if="standard.height_cm !== null && standard.height_cm !== undefined">
                  {{ standard.height_cm }} cm
                </div>
                <div v-else>-</div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">
                <div>L: {{ standard.l_value.toFixed(6) }}</div>
                <div>M: {{ standard.m_value.toFixed(6) }}</div>
                <div>S: {{ standard.s_value.toFixed(6) }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="editStandard(standard)"
                  class="text-jurnal-teal-600 hover:text-jurnal-teal-900 mr-4"
                >
                  Edit
                </button>
                <button
                  @click="deleteStandard(standard)"
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
          {{ Math.min(pagination.page * pagination.limit, pagination.total) }} dari {{ pagination.total }} standards
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
            {{ editingStandard.id ? 'Edit WHO Standard' : 'Tambah WHO Standard' }}
          </h2>
        </div>
        <form @submit.prevent="saveStandard" class="p-6 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Indicator *</label>
              <select
                v-model="editingStandard.indicator"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="wfa">WFA - Weight for Age</option>
                <option value="hfa">HFA - Height for Age</option>
                <option value="wfh">WFH - Weight for Height</option>
                <option value="hcfa">HCFA - Head Circumference for Age</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Gender *</label>
              <select
                v-model="editingStandard.gender"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="male">Male</option>
                <option value="female">Female</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Age (Months)</label>
              <input
                v-model.number="editingStandard.age_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
                placeholder="0-60"
              />
              <p class="text-xs text-gray-500 mt-1">Leave empty for weight-for-height</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Height (cm)</label>
              <input
                v-model.number="editingStandard.height_cm"
                type="number"
                step="0.01"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
                placeholder="For weight-for-height only"
              />
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">L Value *</label>
              <input
                v-model.number="editingStandard.l_value"
                type="number"
                step="0.000001"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">M Value *</label>
              <input
                v-model.number="editingStandard.m_value"
                type="number"
                step="0.000001"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">S Value *</label>
              <input
                v-model.number="editingStandard.s_value"
                type="number"
                step="0.000001"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-4 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD -3</label>
              <input
                v-model.number="editingStandard.sd3neg"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD -2</label>
              <input
                v-model.number="editingStandard.sd2neg"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD -1</label>
              <input
                v-model.number="editingStandard.sd1neg"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD 0 (Median)</label>
              <input
                v-model.number="editingStandard.sd0"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD +1</label>
              <input
                v-model.number="editingStandard.sd1"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD +2</label>
              <input
                v-model.number="editingStandard.sd2"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SD +3</label>
              <input
                v-model.number="editingStandard.sd3"
                type="number"
                step="0.0001"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
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

const standards = ref([])
const pagination = ref({ page: 1, limit: 50, total: 0 })
const filters = ref({
  indicator: '',
  gender: ''
})
const showCreateModal = ref(false)
const editingStandard = ref({})
const saving = ref(false)

const fetchStandards = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.indicator) params.append('indicator', filters.value.indicator)
    if (filters.value.gender) params.append('gender', filters.value.gender)

    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/who-standards?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) throw new Error('Failed to fetch WHO standards')

    const data = await response.json()
    standards.value = data.standards || []
    pagination.value = data.pagination || pagination.value
  } catch (error) {
    console.error('Error fetching WHO standards:', error)
    alert('Gagal memuat WHO standards')
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchStandards()
}

const editStandard = (standard: any) => {
  editingStandard.value = { ...standard }
  showCreateModal.value = true
}

const saveStandard = async () => {
  saving.value = true
  try {
    const url = editingStandard.value.id
      ? `${apiBase}/api/admin/who-standards/${editingStandard.value.id}`
      : `${apiBase}/api/admin/who-standards`

    const method = editingStandard.value.id ? 'PUT' : 'POST'

    const payload: any = {
      indicator: editingStandard.value.indicator,
      gender: editingStandard.value.gender,
      l_value: parseFloat(editingStandard.value.l_value),
      m_value: parseFloat(editingStandard.value.m_value),
      s_value: parseFloat(editingStandard.value.s_value),
    }

    if (editingStandard.value.age_months !== null && editingStandard.value.age_months !== undefined) {
      payload.age_months = parseInt(editingStandard.value.age_months)
    }
    if (editingStandard.value.height_cm !== null && editingStandard.value.height_cm !== undefined) {
      payload.height_cm = parseFloat(editingStandard.value.height_cm)
    }
    if (editingStandard.value.sd3neg !== null && editingStandard.value.sd3neg !== undefined) {
      payload.sd3neg = parseFloat(editingStandard.value.sd3neg)
    }
    if (editingStandard.value.sd2neg !== null && editingStandard.value.sd2neg !== undefined) {
      payload.sd2neg = parseFloat(editingStandard.value.sd2neg)
    }
    if (editingStandard.value.sd1neg !== null && editingStandard.value.sd1neg !== undefined) {
      payload.sd1neg = parseFloat(editingStandard.value.sd1neg)
    }
    if (editingStandard.value.sd0 !== null && editingStandard.value.sd0 !== undefined) {
      payload.sd0 = parseFloat(editingStandard.value.sd0)
    }
    if (editingStandard.value.sd1 !== null && editingStandard.value.sd1 !== undefined) {
      payload.sd1 = parseFloat(editingStandard.value.sd1)
    }
    if (editingStandard.value.sd2 !== null && editingStandard.value.sd2 !== undefined) {
      payload.sd2 = parseFloat(editingStandard.value.sd2)
    }
    if (editingStandard.value.sd3 !== null && editingStandard.value.sd3 !== undefined) {
      payload.sd3 = parseFloat(editingStandard.value.sd3)
    }

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
      throw new Error(error.error || 'Failed to save WHO standard')
    }

    showCreateModal.value = false
    editingStandard.value = {}
    fetchStandards()
  } catch (error: any) {
    console.error('Error saving WHO standard:', error)
    alert(error.message || 'Gagal menyimpan WHO standard')
  } finally {
    saving.value = false
  }
}

const deleteStandard = async (standard: any) => {
  if (!confirm(`Yakin ingin menghapus WHO standard ini?`)) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/who-standards/${standard.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to delete WHO standard')
    }

    fetchStandards()
  } catch (error: any) {
    console.error('Error deleting WHO standard:', error)
    alert(error.message || 'Gagal menghapus WHO standard')
  }
}

onMounted(() => {
  fetchStandards()
})
</script>

