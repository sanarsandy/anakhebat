<template>
  <div>
    <div class="mb-6">
      <NuxtLink
        to="/admin/users"
        class="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 mb-4"
      >
        <Icon name="mdi:arrow-left" class="w-4 h-4" />
        <span>Back to Users</span>
      </NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900">User Details</h1>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-jurnal-teal-600 mx-auto"></div>
      <p class="mt-4 text-gray-600">Loading...</p>
    </div>

    <div v-else-if="user" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- User Info -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">User Information</h2>
          <dl class="space-y-4">
            <div>
              <dt class="text-sm font-medium text-gray-500">Full Name</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ user.full_name }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Email</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ user.email || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Phone Number</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ user.phone_number || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Role</dt>
              <dd class="mt-1">
                <span
                  class="px-2 py-1 text-xs font-semibold rounded-full"
                  :class="user.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-gray-100 text-gray-700'"
                >
                  {{ user.role }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Auth Provider</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ user.auth_provider || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Phone Verified</dt>
              <dd class="mt-1">
                <span
                  v-if="user.phone_verified"
                  class="px-2 py-1 text-xs font-semibold rounded-full bg-emerald-100 text-emerald-700"
                >
                  Verified
                </span>
                <span
                  v-else
                  class="px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-700"
                >
                  Not Verified
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Created At</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ formatDate(user.created_at) }}</dd>
            </div>
          </dl>
        </div>

        <!-- Statistics -->
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">Statistics</h2>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p class="text-sm text-gray-600">Children</p>
              <p class="text-2xl font-bold text-gray-900">{{ statistics.children_count || 0 }}</p>
            </div>
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

      <!-- Actions -->
      <div class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
          <h2 class="text-xl font-bold text-gray-900 mb-4">Actions</h2>
          <div class="space-y-3">
            <button
              @click="showEditModal = true; editingUser = { ...user }"
              class="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Edit User
            </button>
            <button
              @click="showResetPasswordModal = true"
              class="w-full px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 transition-colors"
            >
              Reset Password
            </button>
            <button
              @click="deleteUser"
              class="w-full px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
            >
              Delete User
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit User Modal -->
    <div v-if="showEditModal && editingUser" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showEditModal = false; editingUser = null">
      <div class="bg-white rounded-xl shadow-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold text-gray-900 mb-4">Edit User</h3>
        <form @submit.prevent="handleUpdateUser" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Full Name *</label>
            <input v-model="editingUser.full_name" type="text" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input v-model="editingUser.email" type="email" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
            <input v-model="editingUser.phone_number" type="tel" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Role *</label>
            <select v-model="editingUser.role" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent">
              <option value="parent">Parent</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Auth Provider</label>
            <select v-model="editingUser.auth_provider" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent">
              <option value="email">Email</option>
              <option value="phone">Phone</option>
              <option value="google">Google</option>
            </select>
          </div>
          <div class="flex gap-3 pt-4">
            <button type="submit" class="flex-1 px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors">
              Update
            </button>
            <button type="button" @click="showEditModal = false; editingUser = null" class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors">
              Batal
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Reset Password Modal -->
    <div v-if="showResetPasswordModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showResetPasswordModal = false; newPassword = ''">
      <div class="bg-white rounded-xl shadow-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold text-gray-900 mb-4">Reset Password</h3>
        <form @submit.prevent="handleResetPassword" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">New Password *</label>
            <input v-model="newPassword" type="password" required minlength="6" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
            <p class="mt-1 text-xs text-gray-500">Minimal 6 karakter</p>
          </div>
          <div class="flex gap-3 pt-4">
            <button type="submit" class="flex-1 px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 transition-colors">
              Reset Password
            </button>
            <button type="button" @click="showResetPasswordModal = false; newPassword = ''" class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors">
              Batal
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

const route = useRoute()
const authStore = useAuthStore()

const user = ref<any>(null)
const statistics = ref<any>({})
const loading = ref(true)
const showEditModal = ref(false)
const showResetPasswordModal = ref(false)
const editingUser = ref<any>(null)
const newPassword = ref('')

const fetchUser = async () => {
  try {
    loading.value = true
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    const response = await $fetch(`${apiBase}/api/admin/users/${route.params.id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    if (response && response.user) {
      user.value = response.user
      statistics.value = response.statistics || {}
    }
  } catch (error: any) {
    console.error('Failed to fetch user:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat data user'
    alert(errorMsg)
  } finally {
    loading.value = false
  }
}

const deleteUser = async () => {
  if (!confirm(`Apakah Anda yakin ingin menghapus user ${user.value?.full_name}?`)) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    await $fetch(`${apiBase}/api/admin/users/${route.params.id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    navigateTo('/admin/users')
  } catch (error: any) {
    console.error('Failed to delete user:', error)
    alert('Gagal menghapus user')
  }
}

const handleUpdateUser = async () => {
  if (!editingUser.value?.id) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    const updateData: any = {}
    if (editingUser.value.full_name) updateData.full_name = editingUser.value.full_name
    if (editingUser.value.email) updateData.email = editingUser.value.email
    if (editingUser.value.phone_number !== undefined) updateData.phone_number = editingUser.value.phone_number
    if (editingUser.value.role) updateData.role = editingUser.value.role
    if (editingUser.value.auth_provider) updateData.auth_provider = editingUser.value.auth_provider

    await $fetch(`${apiBase}/api/admin/users/${editingUser.value.id}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: updateData
    })

    showEditModal.value = false
    editingUser.value = null
    await fetchUser()
  } catch (error: any) {
    console.error('Failed to update user:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal mengupdate user'
    alert(errorMsg)
  }
}

const handleResetPassword = async () => {
  if (!newPassword.value || newPassword.value.length < 6) {
    alert('Password minimal 6 karakter')
    return
  }

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    await $fetch(`${apiBase}/api/admin/users/${route.params.id}/reset-password`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: {
        new_password: newPassword.value
      }
    })

    showResetPasswordModal.value = false
    newPassword.value = ''
    alert('Password berhasil direset')
  } catch (error: any) {
    console.error('Failed to reset password:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal reset password'
    alert(errorMsg)
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
  fetchUser()
})
</script>

