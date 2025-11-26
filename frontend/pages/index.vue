<template>
  <div class="min-h-screen bg-white">
    <!-- Hero Section -->
    <section class="relative overflow-hidden">
      <!-- Background decorative elements -->
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute top-0 right-0 w-96 h-96 bg-jurnal-teal-50 rounded-full blur-3xl opacity-30 transform translate-x-1/2 -translate-y-1/2"></div>
        <div class="absolute bottom-0 left-0 w-96 h-96 bg-jurnal-coral-50 rounded-full blur-3xl opacity-20 transform -translate-x-1/2 translate-y-1/2"></div>
      </div>

      <div class="relative container mx-auto px-4 py-20 md:py-28">
        <div class="max-w-4xl mx-auto text-center">
          <!-- Logo -->
          <div class="flex justify-center mb-8 animate-fade-in">
            <img 
              v-if="logoUrl"
              :src="logoUrl" 
              alt="Jurnal Si Kecil" 
              class="h-28 md:h-36 w-auto object-contain drop-shadow-sm"
            />
          </div>

          <!-- Main Heading -->
          <h1 class="text-5xl md:text-6xl lg:text-7xl font-bold text-jurnal-charcoal-800 mb-6 leading-tight">
            <span class="text-jurnal-teal-600">Jurnal Si Kecil</span>
          </h1>
          
          <p class="text-xl md:text-2xl text-jurnal-charcoal-600 mb-4 font-medium">
            Tumbuh Kembang
          </p>
          
          <p class="text-lg md:text-xl text-jurnal-charcoal-500 mb-12 max-w-2xl mx-auto leading-relaxed">
            Aplikasi pemantauan tumbuh kembang anak berbasis standar WHO & CDC. 
            Pantau pertumbuhan fisik dan perkembangan mental si kecil dengan mudah dan akurat.
          </p>

          <!-- CTA Buttons -->
          <div v-if="!authStore.isAuthenticated" class="flex flex-col sm:flex-row gap-4 justify-center items-center mb-16">
            <NuxtLink 
              to="/register" 
              class="group inline-flex items-center gap-2 px-8 py-4 bg-jurnal-teal-600 text-white font-medium rounded-pill hover:bg-jurnal-teal-700 transition-all shadow-sm hover:shadow-md transform hover:-translate-y-0.5"
            >
              <span>Mulai Gratis</span>
              <svg class="w-5 h-5 transition-transform group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </NuxtLink>
            <NuxtLink 
              to="/login" 
              class="inline-flex items-center gap-2 px-8 py-4 bg-white text-jurnal-teal-600 font-medium rounded-pill border-2 border-jurnal-teal-600 hover:bg-jurnal-teal-50 transition-all"
            >
              <span>Masuk</span>
            </NuxtLink>
          </div>

          <div v-else class="mb-16">
            <p class="text-jurnal-charcoal-600 mb-6">
              Selamat datang kembali, <span class="font-semibold text-jurnal-teal-600">{{ authStore.user?.full_name }}</span>!
            </p>
            <NuxtLink 
              to="/dashboard" 
              class="group inline-flex items-center gap-2 px-8 py-4 bg-jurnal-teal-600 text-white font-medium rounded-pill hover:bg-jurnal-teal-700 transition-all shadow-sm hover:shadow-md transform hover:-translate-y-0.5"
            >
              <span>Buka Dashboard</span>
              <svg class="w-5 h-5 transition-transform group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </NuxtLink>
          </div>

          <!-- Trust Badges -->
          <div class="flex flex-wrap justify-center items-center gap-6 text-sm text-jurnal-charcoal-500">
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-jurnal-teal-600" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              <span>Standar WHO</span>
            </div>
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-jurnal-teal-600" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              <span>CDC Milestones</span>
            </div>
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-jurnal-teal-600" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              <span>Akurat & Terpercaya</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-20 md:py-28 bg-jurnal-charcoal-50">
      <div class="container mx-auto px-4">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-jurnal-charcoal-800 mb-4">
            Fitur Utama
          </h2>
          <p class="text-lg text-jurnal-charcoal-600 max-w-2xl mx-auto">
            Semua yang Anda butuhkan untuk memantau tumbuh kembang si kecil
          </p>
        </div>

        <div class="max-w-6xl mx-auto grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <!-- Feature 1: Pertumbuhan Fisik -->
          <div class="group bg-white p-8 rounded-soft-lg shadow-sm border border-jurnal-charcoal-200 hover:border-jurnal-teal-300 transition-all hover:shadow-md">
            <div class="w-14 h-14 rounded-soft-lg bg-jurnal-teal-100 flex items-center justify-center mb-6 group-hover:bg-jurnal-teal-200 transition-colors">
              <svg class="w-7 h-7 text-jurnal-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-jurnal-charcoal-800 mb-3">Pertumbuhan Fisik</h3>
            <p class="text-jurnal-charcoal-600 mb-4 leading-relaxed">
              Pantau berat badan, tinggi badan, dan lingkar kepala anak sesuai standar WHO dengan grafik interaktif
            </p>
            <ul class="space-y-2 text-sm text-jurnal-charcoal-600">
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-teal-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Z-Score Calculation</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-teal-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Deteksi Status Gizi</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-teal-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Grafik Pertumbuhan</span>
              </li>
            </ul>
          </div>

          <!-- Feature 2: Perkembangan Mental -->
          <div class="group bg-white p-8 rounded-soft-lg shadow-sm border border-jurnal-charcoal-200 hover:border-jurnal-coral-300 transition-all hover:shadow-md">
            <div class="w-14 h-14 rounded-soft-lg bg-jurnal-coral-100 flex items-center justify-center mb-6 group-hover:bg-jurnal-coral-200 transition-colors">
              <svg class="w-7 h-7 text-jurnal-coral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-jurnal-charcoal-800 mb-3">Perkembangan Mental</h3>
            <p class="text-jurnal-charcoal-600 mb-4 leading-relaxed">
              Lacak milestone perkembangan motorik, sensorik, bahasa, dan sosial berdasarkan CDC/KPSP
            </p>
            <ul class="space-y-2 text-sm text-jurnal-charcoal-600">
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-coral-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Checklist Milestone</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-coral-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Red Flags Detection</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-coral-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Rekomendasi Stimulasi</span>
              </li>
            </ul>
          </div>

          <!-- Feature 3: Imunisasi -->
          <div class="group bg-white p-8 rounded-soft-lg shadow-sm border border-jurnal-charcoal-200 hover:border-jurnal-gold-300 transition-all hover:shadow-md">
            <div class="w-14 h-14 rounded-soft-lg bg-jurnal-gold-100 flex items-center justify-center mb-6 group-hover:bg-jurnal-gold-200 transition-colors">
              <svg class="w-7 h-7 text-jurnal-gold-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-jurnal-charcoal-800 mb-3">Jadwal Imunisasi</h3>
            <p class="text-jurnal-charcoal-600 mb-4 leading-relaxed">
              Pantau jadwal imunisasi si kecil sesuai rekomendasi IDAI dengan notifikasi pengingat
            </p>
            <ul class="space-y-2 text-sm text-jurnal-charcoal-600">
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-gold-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Jadwal Lengkap</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-gold-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Notifikasi Pengingat</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-5 h-5 text-jurnal-gold-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span>Riwayat Vaksinasi</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 md:py-28 bg-white">
      <div class="container mx-auto px-4">
        <div class="max-w-3xl mx-auto text-center">
          <div class="bg-gradient-to-br from-jurnal-teal-50 to-jurnal-teal-100 rounded-soft-lg p-12 border border-jurnal-teal-200">
            <h2 class="text-3xl md:text-4xl font-bold text-jurnal-charcoal-800 mb-4">
              Siap Memulai?
            </h2>
            <p class="text-lg text-jurnal-charcoal-600 mb-8">
              Daftar sekarang dan mulai pantau tumbuh kembang si kecil dengan mudah
            </p>
            <div v-if="!authStore.isAuthenticated" class="flex flex-col sm:flex-row gap-4 justify-center">
              <NuxtLink 
                to="/register" 
                class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-jurnal-teal-600 text-white font-medium rounded-pill hover:bg-jurnal-teal-700 transition-all shadow-sm hover:shadow-md"
              >
                <span>Daftar Gratis</span>
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
              </NuxtLink>
              <NuxtLink 
                to="/login" 
                class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-white text-jurnal-teal-600 font-medium rounded-pill border-2 border-jurnal-teal-600 hover:bg-jurnal-teal-50 transition-all"
              >
                <span>Masuk</span>
              </NuxtLink>
            </div>
            <div v-else>
              <NuxtLink 
                to="/dashboard" 
                class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-jurnal-teal-600 text-white font-medium rounded-pill hover:bg-jurnal-teal-700 transition-all shadow-sm hover:shadow-md"
              >
                <span>Buka Dashboard</span>
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-jurnal-charcoal-50 border-t border-jurnal-charcoal-200 py-12">
      <div class="container mx-auto px-4">
        <div class="text-center">
          <p class="text-jurnal-charcoal-500 text-sm mb-2">
            © 2025 Jurnal Si Kecil. Semua hak dilindungi.
          </p>
          <p class="text-jurnal-charcoal-400 text-xs">
            Standar WHO Anthropometry • CDC Milestones • Piramida Belajar
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import logoImage from '~/assets/images/logo.png'

definePageMeta({
  layout: false
})

const authStore = useAuthStore()
const logoUrl = ref(null)

onMounted(() => {
  authStore.initialize()
  // Set logo URL only on client to avoid hydration mismatch
  logoUrl.value = logoImage
})
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.6s ease-out;
}
</style>
