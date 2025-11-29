<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Reports</h1>
      <p class="text-gray-600 mt-1">Generate dan export laporan data</p>
    </div>

    <!-- Report Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Users Report -->
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center gap-4 mb-4">
          <div class="w-12 h-12 rounded-lg bg-jurnal-teal-100 flex items-center justify-center">
            <Icon name="mdi:account-group" class="w-6 h-6 text-jurnal-teal-600" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Users Report</h3>
            <p class="text-sm text-gray-600">Laporan semua users</p>
          </div>
        </div>
        <div class="space-y-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select v-model="usersFilters.role" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm">
              <option value="">All Roles</option>
              <option value="parent">Parent</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Auth Provider</label>
            <select v-model="usersFilters.auth_provider" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm">
              <option value="">All Providers</option>
              <option value="email">Email</option>
              <option value="phone">Phone</option>
              <option value="google">Google</option>
            </select>
          </div>
          <div class="flex gap-2 pt-2">
            <button
              @click="generateUsersReport('json')"
              :disabled="usersLoading"
              class="flex-1 px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors disabled:opacity-50 text-sm"
            >
              View JSON
            </button>
            <button
              @click="generateUsersReport('csv')"
              :disabled="usersLoading"
              class="flex-1 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50 text-sm"
            >
              Export CSV
            </button>
          </div>
        </div>
      </div>

      <!-- Children Report -->
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center gap-4 mb-4">
          <div class="w-12 h-12 rounded-lg bg-jurnal-coral-100 flex items-center justify-center">
            <Icon name="mdi:baby-face-outline" class="w-6 h-6 text-jurnal-coral-600" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Children Report</h3>
            <p class="text-sm text-gray-600">Laporan semua children</p>
          </div>
        </div>
        <div class="space-y-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Gender</label>
            <select v-model="childrenFilters.gender" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm">
              <option value="">All Genders</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>
          </div>
          <div class="flex gap-2 pt-2">
            <button
              @click="generateChildrenReport('json')"
              :disabled="childrenLoading"
              class="flex-1 px-4 py-2 bg-jurnal-coral-600 text-white rounded-lg hover:bg-jurnal-coral-700 transition-colors disabled:opacity-50 text-sm"
            >
              View JSON
            </button>
            <button
              @click="generateChildrenReport('csv')"
              :disabled="childrenLoading"
              class="flex-1 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50 text-sm"
            >
              Export CSV
            </button>
          </div>
        </div>
      </div>

      <!-- Growth Report -->
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center gap-4 mb-4">
          <div class="w-12 h-12 rounded-lg bg-jurnal-gold-100 flex items-center justify-center">
            <Icon name="mdi:chart-line" class="w-6 h-6 text-jurnal-gold-600" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Growth Report</h3>
            <p class="text-sm text-gray-600">Laporan measurements</p>
          </div>
        </div>
        <div class="space-y-3">
          <div class="flex gap-2 pt-2">
            <button
              @click="generateGrowthReport('json')"
              :disabled="growthLoading"
              class="flex-1 px-4 py-2 bg-jurnal-gold-600 text-white rounded-lg hover:bg-jurnal-gold-700 transition-colors disabled:opacity-50 text-sm"
            >
              View JSON
            </button>
            <button
              @click="generateGrowthReport('csv')"
              :disabled="growthLoading"
              class="flex-1 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50 text-sm"
            >
              Export CSV
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Report Results -->
    <div v-if="reportData" class="mt-6 bg-white rounded-xl shadow-sm p-6 border border-gray-200">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-gray-900">Report Results</h3>
        <button
          @click="reportData = null"
          class="text-gray-500 hover:text-gray-700"
        >
          <Icon name="mdi:close" class="w-5 h-5" />
        </button>
      </div>
      <div class="overflow-x-auto">
        <pre class="text-xs bg-gray-50 p-4 rounded-lg overflow-auto max-h-96">{{ JSON.stringify(reportData, null, 2) }}</pre>
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

const usersFilters = ref({
  role: '',
  auth_provider: ''
})

const childrenFilters = ref({
  gender: ''
})

const usersLoading = ref(false)
const childrenLoading = ref(false)
const growthLoading = ref(false)

const reportData = ref<any>(null)

const generateUsersReport = async (format: string) => {
  try {
    usersLoading.value = true
    const params = new URLSearchParams()
    if (usersFilters.value.role) params.append('role', usersFilters.value.role)
    if (usersFilters.value.auth_provider) params.append('auth_provider', usersFilters.value.auth_provider)
    params.append('format', format)

    if (format === 'csv') {
      // Download CSV
      const url = `${apiBase}/api/admin/reports/users?${params.toString()}`
      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || 'Failed to generate CSV report')
      }
      const blob = await response.blob()
      const downloadUrl = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = downloadUrl
      a.download = 'users_report.csv'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(downloadUrl)
    } else {
      // Show JSON
      const response = await $fetch(`${apiBase}/api/admin/reports/users?${params.toString()}`, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      reportData.value = response
    }
  } catch (error: any) {
    console.error('Failed to generate users report:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal generate report'
    alert(errorMsg)
  } finally {
    usersLoading.value = false
  }
}

const generateChildrenReport = async (format: string) => {
  try {
    childrenLoading.value = true
    const params = new URLSearchParams()
    if (childrenFilters.value.gender) params.append('gender', childrenFilters.value.gender)
    params.append('format', format)

    if (format === 'csv') {
      const url = `${apiBase}/api/admin/reports/children?${params.toString()}`
      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || 'Failed to generate CSV report')
      }
      const blob = await response.blob()
      const downloadUrl = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = downloadUrl
      a.download = 'children_report.csv'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(downloadUrl)
    } else {
      const response = await $fetch(`${apiBase}/api/admin/reports/children?${params.toString()}`, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      reportData.value = response
    }
  } catch (error: any) {
    console.error('Failed to generate children report:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal generate report'
    alert(errorMsg)
  } finally {
    childrenLoading.value = false
  }
}

const generateGrowthReport = async (format: string) => {
  try {
    growthLoading.value = true
    const params = new URLSearchParams()
    params.append('format', format)

    if (format === 'csv') {
      const url = `${apiBase}/api/admin/reports/growth?${params.toString()}`
      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || 'Failed to generate CSV report')
      }
      const blob = await response.blob()
      const downloadUrl = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = downloadUrl
      a.download = 'growth_report.csv'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(downloadUrl)
    } else {
      const response = await $fetch(`${apiBase}/api/admin/reports/growth?${params.toString()}`, {
        headers: {
          Authorization: `Bearer ${authStore.token}`
        }
      })
      reportData.value = response
    }
  } catch (error: any) {
    console.error('Failed to generate growth report:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal generate report'
    alert(errorMsg)
  } finally {
    growthLoading.value = false
  }
}
</script>

