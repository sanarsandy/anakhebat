<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Audit Logs</h1>
        <p class="text-gray-600 mt-1">Log aktivitas admin dan perubahan sistem</p>
      </div>
      <button
        @click="exportLogs"
        :disabled="exporting"
        class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50 flex items-center gap-2"
      >
        <Icon name="mdi:download" class="w-5 h-5" />
        <span>Export CSV</span>
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Action</label>
          <select
            v-model="filters.action"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchLogs"
          >
            <option value="">All Actions</option>
            <option value="create">Create</option>
            <option value="update">Update</option>
            <option value="delete">Delete</option>
            <option value="view">View</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Resource Type</label>
          <select
            v-model="filters.resource_type"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchLogs"
          >
            <option value="">All Resources</option>
            <option value="user">User</option>
            <option value="child">Child</option>
            <option value="milestone">Milestone</option>
            <option value="system_setting">System Setting</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date From</label>
          <input
            v-model="filters.date_from"
            type="date"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchLogs"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date To</label>
          <input
            v-model="filters.date_to"
            type="date"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchLogs"
          />
        </div>
      </div>
    </div>

    <!-- Logs Table -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div v-if="loading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading...</p>
      </div>

      <div v-else>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Resource</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP Address</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="log in logs" :key="log.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ log.user_name || 'System' }}</div>
                    <div class="text-sm text-gray-500">{{ log.user_email || '-' }}</div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="{
                      'bg-green-100 text-green-700': log.action === 'create',
                      'bg-blue-100 text-blue-700': log.action === 'update',
                      'bg-red-100 text-red-700': log.action === 'delete',
                      'bg-gray-100 text-gray-700': log.action === 'view'
                    }"
                  >
                    {{ log.action }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div>
                    <div class="text-sm text-gray-900">{{ log.resource_type }}</div>
                    <div v-if="log.resource_id" class="text-xs text-gray-500">{{ log.resource_id.substring(0, 8) }}...</div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ log.ip_address || '-' }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(log.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    @click="viewLog(log.id)"
                    class="text-jurnal-teal-600 hover:text-jurnal-teal-900"
                  >
                    View
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="pagination.total > 0" class="bg-gray-50 px-6 py-4 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-700">
              Showing {{ (pagination.page - 1) * pagination.limit + 1 }} to
              {{ Math.min(pagination.page * pagination.limit, pagination.total) }} of
              {{ pagination.total }} results
            </div>
            <div class="flex gap-2">
              <button
                @click="changePage(pagination.page - 1)"
                :disabled="pagination.page === 1"
                class="px-3 py-1 border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100"
              >
                Previous
              </button>
              <button
                @click="changePage(pagination.page + 1)"
                :disabled="pagination.page * pagination.limit >= pagination.total"
                class="px-3 py-1 border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Log Detail Modal -->
    <div v-if="selectedLog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="selectedLog = null">
      <div class="bg-white rounded-xl shadow-xl p-6 max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold text-gray-900">Audit Log Details</h3>
          <button
            @click="selectedLog = null"
            class="text-gray-500 hover:text-gray-700"
          >
            <Icon name="mdi:close" class="w-6 h-6" />
          </button>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">User</label>
            <p class="text-sm text-gray-900">{{ selectedLog.user_name || 'System' }} ({{ selectedLog.user_email || '-' }})</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Action</label>
            <p class="text-sm text-gray-900">{{ selectedLog.action }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Resource Type</label>
            <p class="text-sm text-gray-900">{{ selectedLog.resource_type }}</p>
          </div>
          <div v-if="selectedLog.resource_id">
            <label class="block text-sm font-medium text-gray-700 mb-1">Resource ID</label>
            <p class="text-sm text-gray-900 font-mono">{{ selectedLog.resource_id }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">IP Address</label>
            <p class="text-sm text-gray-900">{{ selectedLog.ip_address || '-' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">User Agent</label>
            <p class="text-sm text-gray-900">{{ selectedLog.user_agent || '-' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Created At</label>
            <p class="text-sm text-gray-900">{{ formatDate(selectedLog.created_at) }}</p>
          </div>
          <div v-if="selectedLog.before_data && Object.keys(selectedLog.before_data).length > 0">
            <label class="block text-sm font-medium text-gray-700 mb-1">Before Data</label>
            <pre class="text-xs bg-gray-50 p-3 rounded-lg overflow-auto">{{ JSON.stringify(selectedLog.before_data, null, 2) }}</pre>
          </div>
          <div v-if="selectedLog.after_data && Object.keys(selectedLog.after_data).length > 0">
            <label class="block text-sm font-medium text-gray-700 mb-1">After Data</label>
            <pre class="text-xs bg-gray-50 p-3 rounded-lg overflow-auto">{{ JSON.stringify(selectedLog.after_data, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const authStore = useAuthStore()
const apiBase = useApiUrl()

const logs = ref<any[]>([])
const loading = ref(false)
const exporting = ref(false)
const selectedLog = ref<any>(null)

const filters = ref({
  action: '',
  resource_type: '',
  date_from: '',
  date_to: ''
})

const pagination = ref({
  page: 1,
  limit: 50,
  total: 0
})

const fetchLogs = async () => {
  try {
    loading.value = true
    const params = new URLSearchParams()
    params.append('page', pagination.value.page.toString())
    params.append('limit', pagination.value.limit.toString())
    if (filters.value.action) params.append('action', filters.value.action)
    if (filters.value.resource_type) params.append('resource_type', filters.value.resource_type)
    if (filters.value.date_from) params.append('date_from', filters.value.date_from)
    if (filters.value.date_to) params.append('date_to', filters.value.date_to)

    const response = await $fetch(`${apiBase}/api/admin/audit-logs?${params.toString()}`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    logs.value = response.logs || []
    if (response.pagination) {
      pagination.value.total = response.pagination.total
    }
  } catch (error: any) {
    console.error('Failed to fetch audit logs:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat audit logs'
    alert(errorMsg)
  } finally {
    loading.value = false
  }
}

const viewLog = async (logId: string) => {
  try {
    const response = await $fetch(`${apiBase}/api/admin/audit-logs/${logId}`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    selectedLog.value = response
  } catch (error: any) {
    console.error('Failed to fetch log details:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat detail log'
    alert(errorMsg)
  }
}

const exportLogs = async () => {
  try {
    exporting.value = true
    const params = new URLSearchParams()
    if (filters.value.action) params.append('action', filters.value.action)
    if (filters.value.resource_type) params.append('resource_type', filters.value.resource_type)
    if (filters.value.date_from) params.append('date_from', filters.value.date_from)
    if (filters.value.date_to) params.append('date_to', filters.value.date_to)

    const url = `${apiBase}/api/admin/audit-logs/export?${params.toString()}`
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.error || 'Failed to export audit logs')
    }
    const blob = await response.blob()
    const downloadUrl = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = downloadUrl
    a.download = 'audit_logs.csv'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(downloadUrl)
  } catch (error: any) {
    console.error('Failed to export audit logs:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal export audit logs'
    alert(errorMsg)
  } finally {
    exporting.value = false
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchLogs()
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchLogs()
})
</script>

