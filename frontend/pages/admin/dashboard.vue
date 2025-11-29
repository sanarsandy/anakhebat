<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
      <p class="text-gray-600 mt-1">Overview sistem Jurnal Si Kecil</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Users</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.totalUsers || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-jurnal-teal-100 flex items-center justify-center">
            <Icon name="mdi:account-group" class="w-6 h-6 text-jurnal-teal-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Children</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.totalChildren || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-jurnal-coral-100 flex items-center justify-center">
            <Icon name="mdi:baby-face-outline" class="w-6 h-6 text-jurnal-coral-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Measurements</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.totalMeasurements || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-jurnal-gold-100 flex items-center justify-center">
            <Icon name="mdi:chart-line" class="w-6 h-6 text-jurnal-gold-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Assessments</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.totalAssessments || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center">
            <Icon name="mdi:clipboard-check" class="w-6 h-6 text-purple-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Immunizations</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.totalImmunizations || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
            <Icon name="mdi:needle" class="w-6 h-6 text-green-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Active Users</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.activeUsers || 0 }}</p>
            <p class="text-xs text-gray-500 mt-1">Users with children</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
            <Icon name="mdi:account-check" class="w-6 h-6 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">New Users (30d)</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.newUsersLast30Days || 0 }}</p>
            <p class="text-xs text-gray-500 mt-1">Last 30 days</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-indigo-100 flex items-center justify-center">
            <Icon name="mdi:account-plus" class="w-6 h-6 text-indigo-600" />
          </div>
        </div>
      </div>
    </div>

    <!-- Additional Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Admin Users</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.adminUsers || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center">
            <Icon name="mdi:shield-account" class="w-6 h-6 text-purple-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Parent Users</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.parentUsers || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-teal-100 flex items-center justify-center">
            <Icon name="mdi:account-group" class="w-6 h-6 text-teal-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Male Children</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.maleChildren || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
            <Icon name="mdi:gender-male" class="w-6 h-6 text-blue-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Female Children</p>
            <p class="text-3xl font-bold text-gray-900 mt-2">{{ stats.femaleChildren || 0 }}</p>
          </div>
          <div class="w-12 h-12 rounded-full bg-pink-100 flex items-center justify-center">
            <Icon name="mdi:gender-female" class="w-6 h-6 text-pink-600" />
          </div>
        </div>
      </div>
    </div>

    <!-- Detailed Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Children Statistics</h3>
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Premature</span>
            <span class="text-sm font-semibold text-gray-900">{{ stats.prematureChildren || 0 }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Full Term</span>
            <span class="text-sm font-semibold text-gray-900">{{ stats.fullTermChildren || 0 }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">New Children (30d)</span>
            <span class="text-sm font-semibold text-gray-900">{{ stats.newChildrenLast30Days || 0 }}</span>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Users by Auth Provider</h3>
        <div class="space-y-3">
          <div v-for="(count, provider) in stats.usersByAuthProvider" :key="provider" class="flex items-center justify-between">
            <span class="text-sm text-gray-600 capitalize">{{ provider || 'Unknown' }}</span>
            <span class="text-sm font-semibold text-gray-900">{{ count }}</span>
          </div>
          <div v-if="Object.keys(stats.usersByAuthProvider).length === 0" class="text-sm text-gray-500">
            No data available
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <NuxtLink
        to="/admin/users"
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-200 hover:border-jurnal-teal-300 hover:shadow-md transition-all"
      >
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-lg bg-jurnal-teal-100 flex items-center justify-center">
            <Icon name="mdi:account-group" class="w-6 h-6 text-jurnal-teal-600" />
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">User Management</h3>
            <p class="text-sm text-gray-600">Kelola semua user</p>
          </div>
        </div>
      </NuxtLink>

      <NuxtLink
        to="/admin/children"
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-200 hover:border-jurnal-coral-300 hover:shadow-md transition-all"
      >
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-lg bg-jurnal-coral-100 flex items-center justify-center">
            <Icon name="mdi:baby-face-outline" class="w-6 h-6 text-jurnal-coral-600" />
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">Children Management</h3>
            <p class="text-sm text-gray-600">Lihat semua children</p>
          </div>
        </div>
      </NuxtLink>

      <NuxtLink
        to="/admin/milestones"
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-200 hover:border-jurnal-gold-300 hover:shadow-md transition-all"
      >
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-lg bg-jurnal-gold-100 flex items-center justify-center">
            <Icon name="mdi:flag" class="w-6 h-6 text-jurnal-gold-600" />
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">Milestones</h3>
            <p class="text-sm text-gray-600">Kelola milestones</p>
          </div>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const stats = ref({
  totalUsers: 0,
  totalChildren: 0,
  totalMeasurements: 0,
  totalAssessments: 0,
  totalImmunizations: 0,
  activeUsers: 0,
  adminUsers: 0,
  parentUsers: 0,
  maleChildren: 0,
  femaleChildren: 0,
  prematureChildren: 0,
  fullTermChildren: 0,
  newUsersLast30Days: 0,
  newChildrenLast30Days: 0,
  usersByAuthProvider: {} as Record<string, number>
})

const authStore = useAuthStore()
const config = useRuntimeConfig()

const fetchStats = async () => {
  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      console.error('No authentication token found')
      return
    }

    // Get API base URL with fallback
    const apiBase = config.public.apiBase || 'http://localhost:8080'
    
    const response = await $fetch(`${apiBase}/api/admin/analytics/overview`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    if (response) {
      stats.value = {
        totalUsers: response.total_users || 0,
        totalChildren: response.total_children || 0,
        totalMeasurements: response.total_measurements || 0,
        totalAssessments: response.total_assessments || 0,
        totalImmunizations: response.total_immunizations || 0,
        activeUsers: response.active_users || 0,
        adminUsers: response.admin_users || 0,
        parentUsers: response.parent_users || 0,
        maleChildren: response.male_children || 0,
        femaleChildren: response.female_children || 0,
        prematureChildren: response.premature_children || 0,
        fullTermChildren: response.full_term_children || 0,
        newUsersLast30Days: response.new_users_last_30_days || 0,
        newChildrenLast30Days: response.new_children_last_30_days || 0,
        usersByAuthProvider: response.users_by_auth_provider || {}
      }
    }
  } catch (error: any) {
    console.error('Failed to fetch admin stats:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat statistik'
    console.error('Error details:', errorMsg)
    // Keep default values on error
  }
}

onMounted(() => {
  fetchStats()
})
</script>

