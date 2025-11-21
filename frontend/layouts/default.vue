<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Desktop Sidebar -->
    <Sidebar class="hidden lg:block" />
    
    <!-- Main Content Area -->
    <div class="lg:ml-64">
      <!-- Header -->
      <header class="bg-white shadow-sm sticky top-0 z-10">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">Tukem</h1>
              <p class="text-sm text-gray-500">Tumbuh Kembang Anak</p>
            </div>
            <div class="flex items-center space-x-4">
              <!-- Child Selector (will be implemented) -->
              <ChildSelector v-if="authStore.isAuthenticated" />
              
              <!-- User Menu -->
              <div v-if="authStore.isAuthenticated" class="flex items-center space-x-3">
                <span class="text-sm text-gray-700">{{ authStore.user?.full_name }}</span>
                <button @click="handleLogout" class="text-sm text-red-600 hover:text-red-700">
                  Logout
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <slot />
      </main>
    </div>

    <!-- Mobile Bottom Navigation -->
    <BottomNav class="lg:hidden" />
  </div>
</template>

<script setup>
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
}
</script>
