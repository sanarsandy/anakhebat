<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">User Management</h1>
        <p class="text-gray-600 mt-1">Kelola semua user di sistem</p>
      </div>
      <button
        @click="showCreateModal = true; editingUser = { full_name: '', email: '', password: '', phone_number: '', role: 'parent', auth_provider: 'email' }"
        class="px-4 py-2 bg-jurnal-teal-600 text-white rounded-lg hover:bg-jurnal-teal-700 transition-colors flex items-center gap-2"
      >
        <Icon name="mdi:plus" class="w-5 h-5" />
        <span>Tambah User</span>
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input
            v-model="filters.search"
            type="text"
            placeholder="Cari nama, email, atau phone..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @input="debouncedSearch"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
          <select
            v-model="filters.role"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchUsers"
          >
            <option value="">All Roles</option>
            <option value="parent">Parent</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Auth Provider</label>
          <select
            v-model="filters.auth_provider"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent"
            @change="fetchUsers"
          >
            <option value="">All Providers</option>
            <option value="email">Email</option>
            <option value="phone">Phone</option>
            <option value="google">Google</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Users Table -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                User
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Contact
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Role
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Auth Provider
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Created At
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-10 h-10 rounded-full bg-jurnal-teal-600 flex items-center justify-center mr-3">
                    <span class="text-sm font-semibold text-white">
                      {{ getInitials(user.full_name) }}
                    </span>
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ user.full_name }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ user.email || '-' }}</div>
                <div class="text-sm text-gray-500">{{ user.phone_number || '-' }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  class="px-2 py-1 text-xs font-semibold rounded-full"
                  :class="user.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-gray-100 text-gray-700'"
                >
                  {{ user.role }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ user.auth_provider || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(user.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="viewUser(user.id)"
                  class="text-jurnal-teal-600 hover:text-jurnal-teal-900 mr-3"
                >
                  View
                </button>
                <button
                  @click="editUser(user)"
                  class="text-blue-600 hover:text-blue-900 mr-3"
                >
                  Edit
                </button>
                <button
                  @click="deleteUser(user)"
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

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showCreateModal = false; editingUser = null">
      <div class="bg-white rounded-xl shadow-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold text-gray-900 mb-4">Tambah User Baru</h3>
        <form @submit.prevent="handleCreateUser" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Full Name *</label>
            <input v-model="editingUser.full_name" type="text" required class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input v-model="editingUser.email" type="email" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input v-model="editingUser.password" type="password" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-jurnal-teal-500 focus:border-transparent" />
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
              Simpan
            </button>
            <button type="button" @click="showCreateModal = false; editingUser = null" class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors">
              Batal
            </button>
          </div>
        </form>
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
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const authStore = useAuthStore()

const users = ref<any[]>([])
const filters = ref({
  search: '',
  role: '',
  auth_provider: ''
})
const pagination = ref({
  page: 1,
  limit: 20,
  total: 0
})

let searchTimeout: NodeJS.Timeout | null = null

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    fetchUsers()
  }, 500)
}

const fetchUsers = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString()
    })
    if (filters.value.search) params.append('search', filters.value.search)
    if (filters.value.role) params.append('role', filters.value.role)
    if (filters.value.auth_provider) params.append('auth_provider', filters.value.auth_provider)

    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    const response = await $fetch(`${apiBase}/api/admin/users?${params.toString()}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    if (response && response.users) {
      users.value = response.users
      pagination.value = response.pagination || {
        page: pagination.value.page,
        limit: pagination.value.limit,
        total: 0
      }
    }
  } catch (error: any) {
    console.error('Failed to fetch users:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal memuat data users'
    alert(errorMsg)
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  fetchUsers()
}

const viewUser = (userId: string) => {
  navigateTo(`/admin/users/${userId}`)
}

const editUser = (user: any) => {
  editingUser.value = { ...user }
  showEditModal.value = true
}

const deleteUser = async (user: any) => {
  if (!confirm(`Apakah Anda yakin ingin menghapus user ${user.full_name}?`)) return

  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    await $fetch(`${apiBase}/api/admin/users/${user.id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    await fetchUsers()
  } catch (error: any) {
    console.error('Failed to delete user:', error)
    alert('Gagal menghapus user')
  }
}

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingUser = ref<any>(null)

const handleCreateUser = async () => {
  try {
    const token = authStore.token || useCookie('token').value
    if (!token) {
      throw new Error('No authentication token found')
    }

    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'

    const newUser = {
      full_name: editingUser.value.full_name,
      email: editingUser.value.email || '',
      password: editingUser.value.password || '',
      phone_number: editingUser.value.phone_number || '',
      role: editingUser.value.role || 'parent',
      auth_provider: editingUser.value.auth_provider || 'email'
    }

    await $fetch(`${apiBase}/api/admin/users`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: newUser
    })

    showCreateModal.value = false
    editingUser.value = null
    await fetchUsers()
  } catch (error: any) {
    console.error('Failed to create user:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal membuat user'
    alert(errorMsg)
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
    await fetchUsers()
  } catch (error: any) {
    console.error('Failed to update user:', error)
    const errorMsg = error?.data?.error || error?.message || 'Gagal mengupdate user'
    alert(errorMsg)
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const getInitials = (name: string) => {
  if (!name) return 'A'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

onMounted(() => {
  fetchUsers()
})
</script>

