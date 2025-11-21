<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-blue-50">
    <div class="container mx-auto px-4 py-16">
      <div class="text-center mb-16">
        <h1 class="text-5xl md:text-6xl font-bold text-gray-900 mb-4">
          <span class="text-indigo-600">Tukem</span> - Tumbuh Kembang
        </h1>
        <p class="text-xl md:text-2xl text-gray-600 mb-8">
          Aplikasi Pemantauan Tumbuh Kembang Anak Berbasis Standar WHO & CDC
        </p>
      </div>

      <div class="max-w-5xl mx-auto">
        <div class="grid md:grid-cols-2 gap-8 mb-12">
          <div class="bg-white p-8 rounded-2xl shadow-lg border-2 border-indigo-100 hover:border-indigo-300 transition-all hover:shadow-xl">
            <div class="text-4xl mb-4">📊</div>
            <h3 class="text-2xl font-bold text-indigo-600 mb-3">Pertumbuhan Fisik</h3>
            <p class="text-gray-600 mb-4">Pantau berat badan, tinggi badan, dan lingkar kepala anak sesuai standar WHO dengan grafik interaktif</p>
            <ul class="text-sm text-gray-500 space-y-2">
              <li>✓ Z-Score Calculation</li>
              <li>✓ Deteksi Status Gizi</li>
              <li>✓ Grafik Pertumbuhan</li>
            </ul>
          </div>
          
          <div class="bg-white p-8 rounded-2xl shadow-lg border-2 border-emerald-100 hover:border-emerald-300 transition-all hover:shadow-xl">
            <div class="text-4xl mb-4">🧠</div>
            <h3 class="text-2xl font-bold text-emerald-600 mb-3">Perkembangan Mental</h3>
            <p class="text-gray-600 mb-4">Lacak milestone perkembangan motorik, sensorik, bahasa, dan sosial berdasarkan CDC/KPSP</p>
            <ul class="text-sm text-gray-500 space-y-2">
              <li>✓ Checklist Milestone</li>
              <li>✓ Red Flags Detection</li>
              <li>✓ Rekomendasi Stimulasi</li>
            </ul>
          </div>
        </div>

        <div class="text-center bg-white rounded-2xl shadow-xl p-8">
          <div v-if="!authStore.isAuthenticated" class="space-y-6">
            <p class="text-gray-700 text-lg">Mulai pantau tumbuh kembang anak Anda sekarang!</p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
              <NuxtLink to="/register" class="inline-block px-8 py-4 bg-indigo-600 text-white font-semibold rounded-xl hover:bg-indigo-700 transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-0.5">
                Daftar Gratis
              </NuxtLink>
              <NuxtLink to="/login" class="inline-block px-8 py-4 bg-white text-indigo-600 font-semibold rounded-xl border-2 border-indigo-600 hover:bg-indigo-50 transition-all">
                Login
              </NuxtLink>
            </div>
          </div>
          <div v-else class="space-y-4">
            <p class="text-gray-700 text-lg">Selamat datang kembali, <span class="font-bold text-indigo-600">{{ authStore.user?.full_name }}</span>!</p>
            <NuxtLink to="/dashboard" class="inline-block px-8 py-4 bg-indigo-600 text-white font-semibold rounded-xl hover:bg-indigo-700 transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-0.5">
              Buka Dashboard
            </NuxtLink>
          </div>
        </div>

        <div class="mt-12 text-center">
          <p class="text-gray-500 text-sm">Standar WHO Anthropometry • CDC Milestones • Piramida Belajar</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: false
})

const authStore = useAuthStore()

onMounted(() => {
  authStore.initialize()
})
</script>
