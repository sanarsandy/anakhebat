<template>
  <div class="pyramid-container relative w-full max-w-md mx-auto h-80 flex flex-col-reverse items-center justify-center p-4">
    <!-- Level 1: Sensorik (Base) -->
    <div 
      class="w-full h-20 flex items-center justify-center text-white font-bold transition-all duration-700 relative clip-trapezoid-bottom group hover:scale-105 cursor-pointer shadow-lg"
      :class="getGradient(1)"
    >
      <div class="text-center z-10 transform group-hover:-translate-y-1 transition-transform">
        <div class="text-sm font-bold tracking-wider drop-shadow-md">SENSORIK</div>
        <div class="text-xs font-medium opacity-90 bg-white/20 px-2 py-0.5 rounded-full mt-1 backdrop-blur-sm">{{ getProgress(1) }}%</div>
      </div>
      <!-- Shine effect -->
      <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shine"></div>
    </div>

    <!-- Level 2: Motorik -->
    <div 
      class="w-3/4 h-20 flex items-center justify-center text-white font-bold transition-all duration-700 delay-100 relative clip-trapezoid group hover:scale-105 cursor-pointer shadow-lg"
      :class="getGradient(2)"
    >
      <div class="text-center z-10 transform group-hover:-translate-y-1 transition-transform">
        <div class="text-sm font-bold tracking-wider drop-shadow-md">MOTORIK</div>
        <div class="text-xs font-medium opacity-90 bg-white/20 px-2 py-0.5 rounded-full mt-1 backdrop-blur-sm">{{ getProgress(2) }}%</div>
      </div>
      <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shine"></div>
    </div>

    <!-- Level 3: Persepsi -->
    <div 
      class="w-1/2 h-20 flex items-center justify-center text-white font-bold transition-all duration-700 delay-200 relative clip-trapezoid group hover:scale-105 cursor-pointer shadow-lg"
      :class="getGradient(3)"
    >
      <div class="text-center z-10 transform group-hover:-translate-y-1 transition-transform">
        <div class="text-sm font-bold tracking-wider drop-shadow-md">PERSEPSI</div>
        <div class="text-xs font-medium opacity-90 bg-white/20 px-2 py-0.5 rounded-full mt-1 backdrop-blur-sm">{{ getProgress(3) }}%</div>
      </div>
      <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shine"></div>
    </div>

    <!-- Level 4: Kognitif (Top) -->
    <div 
      class="w-1/4 h-20 flex items-center justify-center text-white font-bold transition-all duration-700 delay-300 relative clip-triangle-top group hover:scale-105 cursor-pointer shadow-lg"
      :class="getGradient(4)"
    >
      <div class="text-center z-10 transform group-hover:-translate-y-1 transition-transform">
        <div class="text-sm font-bold tracking-wider drop-shadow-md">KOGNITIF</div>
        <div class="text-xs font-medium opacity-90 bg-white/20 px-2 py-0.5 rounded-full mt-1 backdrop-blur-sm">{{ getProgress(4) }}%</div>
      </div>
      <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shine"></div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  progress: {
    type: Object,
    default: () => ({
      sensory: 0,
      motor: 0,
      perception: 0,
      cognitive: 0
    })
  }
})

const getProgress = (level) => {
  const map = { 1: 'sensory', 2: 'motor', 3: 'perception', 4: 'cognitive' }
  return Math.round(props.progress[map[level]] || 0)
}

const getGradient = (level) => {
  const p = getProgress(level)
  if (p >= 80) return 'bg-gradient-to-br from-emerald-400 to-emerald-600'
  if (p >= 50) return 'bg-gradient-to-br from-amber-400 to-amber-600'
  return 'bg-gradient-to-br from-gray-300 to-gray-400'
}
</script>

<style scoped>
.clip-trapezoid-bottom {
  clip-path: polygon(0 100%, 100% 100%, 90% 0, 10% 0);
  margin-top: -4px;
  filter: drop-shadow(0 4px 6px rgba(0,0,0,0.1));
}
.clip-trapezoid {
  clip-path: polygon(0 100%, 100% 100%, 85% 0, 15% 0);
  margin-top: -4px;
  filter: drop-shadow(0 4px 6px rgba(0,0,0,0.1));
}
.clip-triangle-top {
  clip-path: polygon(0 100%, 100% 100%, 50% 0);
  margin-top: -4px;
  filter: drop-shadow(0 4px 6px rgba(0,0,0,0.1));
}

@keyframes shine {
  100% {
    transform: translateX(100%);
  }
}

.animate-shine {
  animation: shine 1.5s infinite;
}
</style>
