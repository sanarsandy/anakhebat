<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Overlay -->
    <div
      v-if="isSidebarOpen"
      class="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
      @click="isSidebarOpen = false"
    ></div>

    <!-- Admin Sidebar -->
    <aside
      class="fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-gray-200 transform transition-transform duration-300 ease-in-out md:translate-x-0"
      :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
    >
      <div class="flex flex-col h-full">
        <!-- Logo -->
        <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-200">
          <img 
            v-if="logoUrl"
            :src="logoUrl" 
            alt="Jurnal Si Kecil" 
            class="h-8 w-auto object-contain"
          />
          <div>
            <h1 class="text-lg font-bold text-jurnal-teal-600">Jurnal Si Kecil</h1>
            <p class="text-xs text-gray-500">Admin Panel</p>
          </div>
          <!-- Close button for mobile -->
          <button
            @click="isSidebarOpen = false"
            class="ml-auto md:hidden p-2 text-gray-500 hover:text-gray-700"
          >
            <Icon name="mdi:close" class="w-6 h-6" />
          </button>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 px-4 py-4 space-y-1 overflow-y-auto">
          <NuxtLink
            to="/admin/dashboard"
            @click="closeSidebarOnMobile"
            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="isActive('/admin/dashboard') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <Icon name="mdi:view-dashboard" class="w-5 h-5" />
            <span>Dashboard</span>
          </NuxtLink>

          <NuxtLink
            to="/admin/users"
            @click="closeSidebarOnMobile"
            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="isActive('/admin/users') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <Icon name="mdi:account-group" class="w-5 h-5" />
            <span>Users</span>
          </NuxtLink>

          <NuxtLink
            to="/admin/children"
            @click="closeSidebarOnMobile"
            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="isActive('/admin/children') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <Icon name="mdi:baby-face-outline" class="w-5 h-5" />
            <span>Children</span>
          </NuxtLink>

          <NuxtLink
            to="/admin/analytics"
            @click="closeSidebarOnMobile"
            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="isActive('/admin/analytics') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <Icon name="mdi:chart-bar" class="w-5 h-5" />
            <span>Analytics</span>
          </NuxtLink>

          <NuxtLink
            to="/admin/reports"
            @click="closeSidebarOnMobile"
            class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="isActive('/admin/reports') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <Icon name="mdi:file-document" class="w-5 h-5" />
            <span>Reports</span>
          </NuxtLink>

          <div class="pt-4 mt-4 border-t border-gray-200">
            <p class="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              System
            </p>
            <NuxtLink
              to="/admin/settings"
              @click="closeSidebarOnMobile"
              class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="isActive('/admin/settings') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
            >
              <Icon name="mdi:cog" class="w-5 h-5" />
              <span>Settings</span>
            </NuxtLink>
            <NuxtLink
              to="/admin/audit-logs"
              @click="closeSidebarOnMobile"
              class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="isActive('/admin/audit-logs') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
            >
              <Icon name="mdi:file-document-outline" class="w-5 h-5" />
              <span>Audit Logs</span>
            </NuxtLink>
          </div>

          <div class="pt-4 mt-4 border-t border-gray-200">
            <p class="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              Master Data
            </p>
            <NuxtLink
              to="/admin/milestones"
              @click="closeSidebarOnMobile"
              class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="isActive('/admin/milestones') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
            >
              <Icon name="mdi:flag" class="w-5 h-5" />
              <span>Milestones</span>
            </NuxtLink>

            <NuxtLink
              to="/admin/who-standards"
              @click="closeSidebarOnMobile"
              class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="isActive('/admin/who-standards') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
            >
              <Icon name="mdi:chart-line" class="w-5 h-5" />
              <span>WHO Standards</span>
            </NuxtLink>

            <NuxtLink
              to="/admin/stimulation-content"
              @click="closeSidebarOnMobile"
              class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="isActive('/admin/stimulation-content') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
            >
              <Icon name="mdi:book-open-variant" class="w-5 h-5" />
              <span>Stimulation Content</span>
            </NuxtLink>

            <NuxtLink
              to="/admin/immunization-schedules"
              @click="closeSidebarOnMobile"
              class="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="isActive('/admin/immunization-schedules') ? 'bg-jurnal-teal-50 text-jurnal-teal-700' : 'text-gray-700 hover:bg-gray-100'"
            >
              <Icon name="mdi:needle" class="w-5 h-5" />
              <span>Immunization Schedules</span>
            </NuxtLink>
          </div>
        </nav>

        <!-- User Section -->
        <div class="p-4 border-t border-gray-200">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-10 h-10 rounded-full bg-jurnal-teal-600 flex items-center justify-center">
              <span class="text-sm font-semibold text-white">
                {{ getInitials(authStore.user?.full_name || '') }}
              </span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate">
                {{ authStore.user?.full_name }}
              </p>
              <p class="text-xs text-gray-500 truncate">Admin</p>
            </div>
          </div>
          <button
            @click="handleLogout"
            class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <Icon name="mdi:logout" class="w-5 h-5" />
            <span>Logout</span>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="md:ml-64">
      <!-- Header -->
      <header class="bg-white border-b border-gray-200 sticky top-0 z-30">
        <div class="px-4 md:px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <!-- Hamburger Menu Button (Mobile Only) -->
              <button
                @click="isSidebarOpen = true"
                class="md:hidden p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <Icon name="mdi:menu" class="w-6 h-6" />
              </button>
              <h2 class="text-xl md:text-2xl font-bold text-gray-900">{{ pageTitle }}</h2>
            </div>
            <NuxtLink
              to="/dashboard"
              class="flex items-center gap-2 px-3 md:px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <Icon name="mdi:arrow-left" class="w-5 h-5" />
              <span class="hidden sm:inline">Back to App</span>
            </NuxtLink>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="p-4 md:p-6">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const route = useRoute()

const logoUrl = ref<string | null>(null)
const isSidebarOpen = ref(false)

// Close sidebar when route changes on mobile
if (process.client) {
  watch(() => route.path, () => {
    if (window.innerWidth < 768) {
      isSidebarOpen.value = false
    }
  })
}

onMounted(() => {
  try {
    const logoImage = require('~/assets/images/logo.png')
    logoUrl.value = logoImage
  } catch (e) {
    // Logo not found, use fallback
  }
})

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/admin/dashboard': 'Dashboard',
    '/admin/users': 'User Management',
    '/admin/children': 'Children Management',
    '/admin/analytics': 'Analytics',
    '/admin/reports': 'Reports',
    '/admin/settings': 'System Settings',
    '/admin/audit-logs': 'Audit Logs',
    '/admin/milestones': 'Milestones Management',
    '/admin/who-standards': 'WHO Standards Management',
    '/admin/stimulation-content': 'Stimulation Content Management',
    '/admin/immunization-schedules': 'Immunization Schedules Management',
  }
  return titles[route.path] || 'Admin Panel'
})

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/')
}

const getInitials = (name: string) => {
  if (!name) return 'A'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

const handleLogout = () => {
  authStore.logout()
  navigateTo('/admin/login')
}

const closeSidebarOnMobile = () => {
  if (process.client && window.innerWidth < 768) {
    isSidebarOpen.value = false
  }
}
</script>

