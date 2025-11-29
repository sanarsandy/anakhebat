<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">System Settings</h1>
      <p class="text-gray-600 mt-1">Kelola konfigurasi sistem</p>
    </div>

    <!-- Category Tabs -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 mb-6">
      <div class="border-b border-gray-200">
        <nav class="flex -mb-px overflow-x-auto">
          <button
            v-for="category in categories"
            :key="category.id"
            @click="activeCategory = category.id"
            class="px-6 py-4 text-sm font-medium border-b-2 transition-colors whitespace-nowrap"
            :class="activeCategory === category.id ? 'border-jurnal-teal-600 text-jurnal-teal-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
          >
            {{ category.label }}
          </button>
        </nav>
      </div>
    </div>

    <!-- Settings List -->
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
      <p class="mt-4 text-gray-600">Loading...</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="setting in filteredSettings"
        :key="setting.id"
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-200"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <h3 class="text-lg font-semibold text-gray-900">{{ setting.key }}</h3>
            <p class="text-sm text-gray-600 mt-1">{{ setting.description || 'No description' }}</p>
            <div class="mt-3">
              <label class="block text-sm font-medium text-gray-700 mb-2">Value</label>
              <input
                v-if="setting.type === 'boolean'"
                type="checkbox"
                :checked="setting.value === 'true'"
                @change="updateSetting(setting.key, $event.target.checked ? 'true' : 'false')"
                class="w-5 h-5 text-jurnal-teal-600 rounded focus:ring-jurnal-teal-500"
              />
              <input
                v-else
                v-model="editingValues[setting.key]"
                :type="setting.type === 'number' ? 'number' : 'text'"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
                @blur="updateSetting(setting.key, editingValues[setting.key])"
                @keyup.enter="updateSetting(setting.key, editingValues[setting.key])"
              />
            </div>
            <div class="mt-2 flex items-center gap-4 text-xs text-gray-500">
              <span>Type: {{ setting.type }}</span>
              <span v-if="setting.updated_at">Updated: {{ formatDate(setting.updated_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="filteredSettings.length === 0" class="text-center py-12 text-gray-500">
        No settings found in this category
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

const categories = [
  { id: '', label: 'All' },
  { id: 'general', label: 'General' },
  { id: 'security', label: 'Security' },
  { id: 'notifications', label: 'Notifications' },
  { id: 'email', label: 'Email' },
  { id: 'features', label: 'Features' }
]

const activeCategory = ref('')
const settings = ref<any[]>([])
const loading = ref(false)
const editingValues = ref<Record<string, string>>({})

const filteredSettings = computed(() => {
  if (!activeCategory.value) return settings.value
  return settings.value.filter(s => s.category === activeCategory.value)
})

const fetchSettings = async () => {
  try {
    loading.value = true
    const params = activeCategory.value ? `?category=${activeCategory.value}` : ''
    const response = await $fetch(`${apiBase}/api/admin/settings${params}`, {
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    })
    settings.value = response.settings || []
    // Initialize editing values
    settings.value.forEach(s => {
      editingValues.value[s.key] = s.value
    })
  } catch (error: any) {
    console.error('Failed to fetch settings:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat settings'
    alert(errorMsg)
  } finally {
    loading.value = false
  }
}

const updateSetting = async (key: string, value: string) => {
  try {
    await $fetch(`${apiBase}/api/admin/settings/${key}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: { value }
    })
    // Update local state
    const setting = settings.value.find(s => s.key === key)
    if (setting) {
      setting.value = value
      setting.updated_at = new Date().toISOString()
    }
  } catch (error: any) {
    console.error('Failed to update setting:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal update setting'
    alert(errorMsg)
    // Revert editing value
    const setting = settings.value.find(s => s.key === key)
    if (setting) {
      editingValues.value[key] = setting.value
    }
  }
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

watch(activeCategory, () => {
  fetchSettings()
})

onMounted(() => {
  fetchSettings()
})
</script>


