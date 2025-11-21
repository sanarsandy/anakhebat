<template>
  <div class="bg-white rounded-xl shadow-sm p-6 border-l-4" :class="borderColor">
    <div class="flex items-start justify-between mb-4">
      <div class="flex-1">
        <p class="text-sm text-gray-500">{{ formatDate(measurement.measurement_date) }}</p>
        <p class="text-xs text-gray-400 mt-1">Usia: {{ measurement.age_display }}</p>
      </div>
      <button 
        @click.stop="handleDelete"
        class="text-red-600 hover:text-red-700 p-2 rounded-lg hover:bg-red-50 transition flex-shrink-0 ml-2 border border-red-200 hover:border-red-300"
        title="Hapus pengukuran"
        type="button"
        aria-label="Hapus pengukuran"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>

    <div class="grid grid-cols-2 gap-4 mb-4">
      <!-- Weight -->
      <div>
        <p class="text-gray-500 text-sm">Berat Badan</p>
        <p class="text-2xl font-bold text-gray-900">{{ measurement.weight }} <span class="text-sm font-normal">kg</span></p>
        <div v-if="measurement.weight_for_age_zscore !== null" class="flex items-center space-x-2 mt-1">
          <span class="text-xs px-2 py-1 rounded-full" :class="getStatusBadgeClass(measurement.weight_for_age_zscore)">
            Z: {{ measurement.weight_for_age_zscore.toFixed(2) }}
          </span>
        </div>
      </div>

      <!-- Height -->
      <div>
        <p class="text-gray-500 text-sm">Tinggi Badan</p>
        <p class="text-2xl font-bold text-gray-900">{{ measurement.height }} <span class="text-sm font-normal">cm</span></p>
        <div v-if="measurement.height_for_age_zscore !== null" class="flex items-center space-x-2 mt-1">
          <span class="text-xs px-2 py-1 rounded-full" :class="getStatusBadgeClass(measurement.height_for_age_zscore)">
            Z: {{ measurement.height_for_age_zscore.toFixed(2) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Head Circumference -->
    <div v-if="measurement.head_circumference" class="mb-4">
      <p class="text-gray-500 text-sm">Lingkar Kepala</p>
      <p class="text-xl font-bold text-gray-900">{{ measurement.head_circumference }} <span class="text-sm font-normal">cm</span></p>
    </div>

    <!-- Status Badges -->
    <div class="flex flex-wrap gap-2">
      <span v-if="measurement.nutritional_status" class="text-xs px-3 py-1 rounded-full" :class="getStatusTextClass(measurement.weight_for_age_zscore)">
        {{ measurement.nutritional_status }}
      </span>
      <span v-if="measurement.height_status" class="text-xs px-3 py-1 rounded-full" :class="getStatusTextClass(measurement.height_for_age_zscore)">
        Tinggi: {{ measurement.height_status }}
      </span>
      <span v-if="measurement.weight_for_height_status" class="text-xs px-3 py-1 rounded-full" :class="getStatusTextClass(measurement.weight_for_height_zscore)">
        {{ measurement.weight_for_height_status }}
      </span>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  measurement: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['delete'])

const handleDelete = () => {
  console.log('Delete button clicked in MeasurementCard:', props.measurement)
  emit('delete', props.measurement)
}

const borderColor = computed(() => {
  const weightZ = props.measurement.weight_for_age_zscore
  const heightZ = props.measurement.height_for_age_zscore
  
  if (weightZ < -2 || heightZ < -2 || weightZ > 2 || heightZ > 2) {
    return 'border-red-500'
  } else if (weightZ < -1 || heightZ < -1 || weightZ > 1 || heightZ > 1) {
    return 'border-amber-500'
  }
  return 'border-emerald-500'
})

const getStatusBadgeClass = (zscore) => {
  if (zscore < -2 || zscore > 2) {
    return 'bg-red-100 text-red-700'
  } else if (zscore < -1 || zscore > 1) {
    return 'bg-amber-100 text-amber-700'
  }
  return 'bg-emerald-100 text-emerald-700'
}

const getStatusTextClass = (zscore) => {
  if (zscore < -2 || zscore > 2) {
    return 'bg-red-50 text-red-700 border border-red-200'
  } else if (zscore < -1 || zscore > 1) {
    return 'bg-amber-50 text-amber-700 border border-amber-200'
  }
  return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}
</script>
