<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Analytics</h1>
      <p class="text-gray-600 mt-1">Analisis data sistem Jurnal Si Kecil</p>
    </div>

    <!-- Tabs -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 mb-6">
      <div class="border-b border-gray-200">
        <nav class="flex -mb-px overflow-x-auto">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="px-6 py-4 text-sm font-medium border-b-2 transition-colors whitespace-nowrap"
            :class="activeTab === tab.id ? 'border-jurnal-teal-600 text-jurnal-teal-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>
    </div>

    <!-- Users Analytics -->
    <div v-if="activeTab === 'users'" class="space-y-6">
      <div v-if="usersLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading...</p>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Users by Month</h3>
          <div class="space-y-2">
            <div v-for="item in usersStats.users_by_month" :key="item.month" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.month }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Top Users by Children</h3>
          <div class="space-y-3">
            <div v-for="user in usersStats.top_users_by_children" :key="user.user_id" class="flex items-center justify-between">
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900 truncate">{{ user.full_name }}</p>
                <p class="text-xs text-gray-500 truncate">{{ user.email }}</p>
              </div>
              <span class="text-sm font-semibold text-jurnal-teal-600 ml-2">{{ user.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Average</h3>
          <div class="space-y-4">
            <div>
              <p class="text-sm text-gray-600">Average Children per User</p>
              <p class="text-2xl font-bold text-gray-900 mt-1">{{ usersStats.avg_children_per_user?.toFixed(2) || 0 }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Children Analytics -->
    <div v-if="activeTab === 'children'" class="space-y-6">
      <div v-if="childrenLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading...</p>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Children by Age Group</h3>
          <div class="space-y-2">
            <div v-for="item in childrenStats.children_by_age_group" :key="item.age_group" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.age_group }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Children Registration Trend</h3>
          <div class="space-y-2">
            <div v-for="item in childrenStats.children_by_month" :key="item.month" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.month }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Measurements Analytics -->
    <div v-if="activeTab === 'measurements'" class="space-y-6">
      <div v-if="measurementsLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading...</p>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Measurements by Month</h3>
          <div class="space-y-2">
            <div v-for="item in measurementsStats.measurements_by_month" :key="item.month" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.month }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Growth Status Distribution</h3>
          <div class="space-y-2">
            <div v-for="item in measurementsStats.growth_status_distribution" :key="item.status" class="flex items-center justify-between">
              <span class="text-sm text-gray-600 capitalize">{{ item.status }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Average</h3>
          <div>
            <p class="text-sm text-gray-600">Average Measurements per Child</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ measurementsStats.avg_measurements_per_child?.toFixed(2) || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Assessments Analytics -->
    <div v-if="activeTab === 'assessments'" class="space-y-6">
      <div v-if="assessmentsLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading...</p>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Assessments by Month</h3>
          <div class="space-y-2">
            <div v-for="item in assessmentsStats.assessments_by_month" :key="item.month" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.month }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Red Flags</h3>
          <div>
            <p class="text-3xl font-bold text-red-600">{{ assessmentsStats.red_flags_detected || 0 }}</p>
            <p class="text-sm text-gray-600 mt-1">Red flags detected</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Immunizations Analytics -->
    <div v-if="activeTab === 'immunizations'" class="space-y-6">
      <div v-if="immunizationsLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Loading...</p>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Immunizations by Month</h3>
          <div class="space-y-2">
            <div v-for="item in immunizationsStats.immunizations_by_month" :key="item.month" class="flex items-center justify-between">
              <span class="text-sm text-gray-600">{{ item.month }}</span>
              <span class="text-sm font-semibold text-gray-900">{{ item.count }}</span>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">On Schedule</h3>
          <div>
            <p class="text-3xl font-bold text-green-600">{{ immunizationsStats.on_schedule_count || 0 }}</p>
            <p class="text-sm text-gray-600 mt-1">On schedule immunizations</p>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Catch Up</h3>
          <div>
            <p class="text-3xl font-bold text-orange-600">{{ immunizationsStats.catch_up_count || 0 }}</p>
            <p class="text-sm text-gray-600 mt-1">Catch-up immunizations</p>
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

const activeTab = ref('users')
const tabs = [
  { id: 'users', label: 'Users' },
  { id: 'children', label: 'Children' },
  { id: 'measurements', label: 'Measurements' },
  { id: 'assessments', label: 'Assessments' },
  { id: 'immunizations', label: 'Immunizations' }
]

const usersStats = ref<any>({})
const childrenStats = ref<any>({})
const measurementsStats = ref<any>({})
const assessmentsStats = ref<any>({})
const immunizationsStats = ref<any>({})

const usersLoading = ref(false)
const childrenLoading = ref(false)
const measurementsLoading = ref(false)
const assessmentsLoading = ref(false)
const immunizationsLoading = ref(false)

const fetchUsersAnalytics = async () => {
  try {
    usersLoading.value = true
    const response = await $fetch(`${apiBase}/api/admin/analytics/users`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    usersStats.value = response || {}
  } catch (error: any) {
    console.error('Failed to fetch users analytics:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat analytics users'
    alert(errorMsg)
    usersStats.value = {}
  } finally {
    usersLoading.value = false
  }
}

const fetchChildrenAnalytics = async () => {
  try {
    childrenLoading.value = true
    const response = await $fetch(`${apiBase}/api/admin/analytics/children`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    childrenStats.value = response || {}
  } catch (error: any) {
    console.error('Failed to fetch children analytics:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat analytics children'
    alert(errorMsg)
    childrenStats.value = {}
  } finally {
    childrenLoading.value = false
  }
}

const fetchMeasurementsAnalytics = async () => {
  try {
    measurementsLoading.value = true
    const response = await $fetch(`${apiBase}/api/admin/analytics/measurements`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    measurementsStats.value = response || {}
  } catch (error: any) {
    console.error('Failed to fetch measurements analytics:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat analytics measurements'
    alert(errorMsg)
    measurementsStats.value = {}
  } finally {
    measurementsLoading.value = false
  }
}

const fetchAssessmentsAnalytics = async () => {
  try {
    assessmentsLoading.value = true
    const response = await $fetch(`${apiBase}/api/admin/analytics/assessments`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    assessmentsStats.value = response || {}
  } catch (error: any) {
    console.error('Failed to fetch assessments analytics:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat analytics assessments'
    alert(errorMsg)
    assessmentsStats.value = {}
  } finally {
    assessmentsLoading.value = false
  }
}

const fetchImmunizationsAnalytics = async () => {
  try {
    immunizationsLoading.value = true
    const response = await $fetch(`${apiBase}/api/admin/analytics/immunizations`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    immunizationsStats.value = response || {}
  } catch (error: any) {
    console.error('Failed to fetch immunizations analytics:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat analytics immunizations'
    alert(errorMsg)
    immunizationsStats.value = {}
  } finally {
    immunizationsLoading.value = false
  }
}

watch(activeTab, (newTab) => {
  if (newTab === 'users' && !usersStats.value.users_by_month) {
    fetchUsersAnalytics()
  } else if (newTab === 'children' && !childrenStats.value.children_by_age_group) {
    fetchChildrenAnalytics()
  } else if (newTab === 'measurements' && !measurementsStats.value.measurements_by_month) {
    fetchMeasurementsAnalytics()
  } else if (newTab === 'assessments' && !assessmentsStats.value.assessments_by_month) {
    fetchAssessmentsAnalytics()
  } else if (newTab === 'immunizations' && !immunizationsStats.value.immunizations_by_month) {
    fetchImmunizationsAnalytics()
  }
})

onMounted(() => {
  fetchUsersAnalytics()
})
</script>

