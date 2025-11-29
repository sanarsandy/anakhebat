<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Stimulation Content Management</h1>
        <p class="text-gray-600 mt-1">Kelola konten stimulasi untuk intervensi perkembangan anak</p>
      </div>
      <button
        @click="showCreateModal = true; editingContent = { is_active: true, content_type: 'article' }"
        class="px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors flex items-center gap-2"
      >
        <Icon name="mdi:plus" class="w-5 h-5" />
        <span>Tambah Konten</span>
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
            placeholder="Cari konten..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @input="debouncedSearch"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select
            v-model="filters.category"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchContents"
          >
            <option value="">All Categories</option>
            <option value="sensory">Sensory</option>
            <option value="motor">Motor</option>
            <option value="perception">Perception</option>
            <option value="cognitive">Cognitive</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Content Type</label>
          <select
            v-model="filters.content_type"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchContents"
          >
            <option value="">All Types</option>
            <option value="article">Article</option>
            <option value="video">Video</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="filters.is_active"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchContents"
          >
            <option value="">All</option>
            <option value="true">Active</option>
            <option value="false">Inactive</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Contents Table -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Title
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Category
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Type
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Age Range
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="content in contents" :key="content.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="text-sm font-medium text-gray-900">{{ content.title }}</div>
                <div v-if="content.description" class="text-xs text-gray-500 mt-1 line-clamp-2">
                  {{ content.description }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full" :class="getCategoryColor(content.category)">
                  {{ content.category }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full" :class="content.content_type === 'video' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800'">
                  {{ content.content_type }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <div v-if="content.age_min_months || content.age_max_months">
                  {{ content.age_min_months || '0' }}-{{ content.age_max_months || '∞' }} months
                </div>
                <div v-else>-</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="content.is_active" class="px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800">
                  Active
                </span>
                <span v-else class="px-2 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-800">
                  Inactive
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="editContent(content)"
                  class="text-jurnal-teal-600 hover:text-jurnal-teal-900 mr-4"
                >
                  Edit
                </button>
                <button
                  @click="deleteContent(content)"
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
          {{ Math.min(pagination.page * pagination.limit, pagination.total) }} dari {{ pagination.total }} contents
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
            {{ editingContent.id ? 'Edit Stimulation Content' : 'Tambah Stimulation Content' }}
          </h2>
        </div>
        <form @submit.prevent="saveContent" class="p-6 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Category *</label>
              <select
                v-model="editingContent.category"
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
              <label class="block text-sm font-medium text-gray-700 mb-1">Content Type *</label>
              <select
                v-model="editingContent.content_type"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              >
                <option value="article">Article</option>
                <option value="video">Video</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Title *</label>
            <input
              v-model="editingContent.title"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              v-model="editingContent.description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">URL *</label>
            <input
              v-model="editingContent.url"
              type="url"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              placeholder="https://..."
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Thumbnail URL</label>
            <input
              v-model="editingContent.thumbnail_url"
              type="url"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              placeholder="https://..."
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Min Age (Months)</label>
              <input
                v-model.number="editingContent.age_min_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Max Age (Months)</label>
              <input
                v-model.number="editingContent.age_max_months"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Milestone ID (Optional)</label>
            <input
              v-model="editingContent.milestone_id"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
              placeholder="UUID of linked milestone"
            />
          </div>

          <div>
            <label class="flex items-center gap-2">
              <input
                v-model="editingContent.is_active"
                type="checkbox"
                class="w-4 h-4 text-jurnal-teal-600 border-gray-300 rounded focus:ring-jurnal-teal-500"
              />
              <span class="text-sm font-medium text-gray-700">Active</span>
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

const contents = ref([])
const pagination = ref({ page: 1, limit: 50, total: 0 })
const filters = ref({
  search: '',
  category: '',
  content_type: '',
  is_active: ''
})
const showCreateModal = ref(false)
const editingContent = ref({})
const saving = ref(false)

let searchTimeout: NodeJS.Timeout | null = null

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchContents()
  }, 500)
}

const fetchContents = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.search) params.append('search', filters.value.search)
    if (filters.value.category) params.append('category', filters.value.category)
    if (filters.value.content_type) params.append('content_type', filters.value.content_type)
    if (filters.value.is_active) params.append('is_active', filters.value.is_active)

    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/stimulation-content?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) throw new Error('Failed to fetch stimulation content')

    const data = await response.json()
    contents.value = data.contents || []
    pagination.value = data.pagination || pagination.value
  } catch (error) {
    console.error('Error fetching stimulation content:', error)
    alert('Gagal memuat stimulation content')
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchContents()
}

const editContent = (content: any) => {
  editingContent.value = { ...content }
  showCreateModal.value = true
}

const saveContent = async () => {
  saving.value = true
  try {
    const url = editingContent.value.id
      ? `${apiBase}/api/admin/stimulation-content/${editingContent.value.id}`
      : `${apiBase}/api/admin/stimulation-content`

    const method = editingContent.value.id ? 'PUT' : 'POST'

    const payload: any = {
      category: editingContent.value.category,
      title: editingContent.value.title,
      content_type: editingContent.value.content_type,
      url: editingContent.value.url,
      is_active: editingContent.value.is_active !== false,
    }

    if (editingContent.value.description) payload.description = editingContent.value.description
    if (editingContent.value.thumbnail_url) payload.thumbnail_url = editingContent.value.thumbnail_url
    if (editingContent.value.age_min_months) payload.age_min_months = parseInt(editingContent.value.age_min_months)
    if (editingContent.value.age_max_months) payload.age_max_months = parseInt(editingContent.value.age_max_months)
    if (editingContent.value.milestone_id) payload.milestone_id = editingContent.value.milestone_id

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
      throw new Error(error.error || 'Failed to save stimulation content')
    }

    showCreateModal.value = false
    editingContent.value = {}
    fetchContents()
  } catch (error: any) {
    console.error('Error saving stimulation content:', error)
    alert(error.message || 'Gagal menyimpan stimulation content')
  } finally {
    saving.value = false
  }
}

const deleteContent = async (content: any) => {
  if (!confirm(`Yakin ingin menghapus konten ini?`)) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${apiBase}/api/admin/stimulation-content/${content.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to delete stimulation content')
    }

    fetchContents()
  } catch (error: any) {
    console.error('Error deleting stimulation content:', error)
    alert(error.message || 'Gagal menghapus stimulation content')
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
  fetchContents()
})
</script>

