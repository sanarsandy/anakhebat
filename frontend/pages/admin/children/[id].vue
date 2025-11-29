<template>
  <div>
    <div class="mb-6">
      <NuxtLink
        to="/admin/children"
        class="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 mb-4"
      >
        <Icon name="mdi:arrow-left" class="w-4 h-4" />
        <span>Back to Children</span>
      </NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900">Child Details</h1>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
      <p class="mt-4 text-gray-600">Loading...</p>
    </div>

    <div v-else-if="child" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Child Info -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">Child Information</h2>
          <dl class="space-y-4">
            <div>
              <dt class="text-sm font-medium text-gray-500">Name</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ child.name }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Date of Birth</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ formatDate(child.dob) }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Gender</dt>
              <dd class="mt-1">
                <span
                  class="px-2 py-1 text-xs font-semibold rounded-full"
                  :class="child.gender === 'male' ? 'bg-blue-100 text-blue-700' : 'bg-pink-100 text-pink-700'"
                >
                  {{ child.gender === 'male' ? 'Laki-laki' : 'Perempuan' }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Birth Weight</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ child.birth_weight }} kg</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Birth Height</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ child.birth_height }} cm</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Premature</dt>
              <dd class="mt-1">
                <span
                  class="px-2 py-1 text-xs font-semibold rounded-full"
                  :class="child.is_premature ? 'bg-orange-100 text-orange-700' : 'bg-gray-100 text-gray-700'"
                >
                  {{ child.is_premature ? 'Ya' : 'Tidak' }}
                </span>
              </dd>
            </div>
            <div v-if="child.is_premature && child.gestational_age">
              <dt class="text-sm font-medium text-gray-500">Gestational Age</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ child.gestational_age }} minggu</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Created At</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ formatDate(child.created_at) }}</dd>
            </div>
          </dl>
        </div>

        <!-- Parent Info -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">Parent Information</h2>
          <dl class="space-y-4">
            <div>
              <dt class="text-sm font-medium text-gray-500">Parent Name</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ parent.name || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Email</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ parent.email || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Phone Number</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ parent.phone_number || '-' }}</dd>
            </div>
          </dl>
        </div>

        <!-- Statistics -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">Statistics</h2>
          <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
            <div>
              <p class="text-sm text-gray-600">Measurements</p>
              <p class="text-2xl font-bold text-gray-900">{{ statistics.measurements_count || 0 }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Assessments</p>
              <p class="text-2xl font-bold text-gray-900">{{ statistics.assessments_count || 0 }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Immunizations</p>
              <p class="text-2xl font-bold text-gray-900">{{ statistics.immunizations_count || 0 }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
          <div class="space-y-3">
            <NuxtLink
              :to="`/children/${child.id}`"
              class="block w-full px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors text-center"
            >
              View Full Profile
            </NuxtLink>
            <NuxtLink
              :to="`/growth/charts?child=${child.id}`"
              class="block w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-center"
            >
              View Growth Charts
            </NuxtLink>
            <NuxtLink
              :to="`/development/denver-ii?child=${child.id}`"
              class="block w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors text-center"
            >
              View Denver II
            </NuxtLink>
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

const route = useRoute()
const authStore = useAuthStore()
const config = useRuntimeConfig()

const child = ref<any>(null)
const parent = ref<any>({})
const statistics = ref<any>({})
const loading = ref(true)

const fetchChild = async () => {
  try {
    loading.value = true
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const apiBase = config.public.apiBase || 'http://localhost:8080'

    const response = await $fetch(`${apiBase}/api/admin/children/${route.params.id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    if (response) {
      child.value = response.child
      parent.value = response.parent || {}
      statistics.value = response.statistics || {}
    }
  } catch (error: any) {
    console.error('Failed to fetch child:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat data child'
    alert(errorMsg)
  } finally {
    loading.value = false
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchChild()
})
</script>


