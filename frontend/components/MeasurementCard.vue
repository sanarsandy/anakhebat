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
        <div v-if="isValidNumber(measurement.weight_for_age_zscore)" class="flex items-center space-x-2 mt-1">
          <span class="text-xs px-2 py-1 rounded-full" :class="getZScoreBadgeClass(measurement.weight_for_age_zscore)">
            Z: {{ formatZScore(measurement.weight_for_age_zscore) }}
          </span>
        </div>
      </div>

      <!-- Height -->
      <div>
        <p class="text-gray-500 text-sm">Tinggi Badan</p>
        <p class="text-2xl font-bold text-gray-900">{{ measurement.height }} <span class="text-sm font-normal">cm</span></p>
        <div v-if="isValidNumber(measurement.height_for_age_zscore)" class="flex items-center space-x-2 mt-1">
          <span class="text-xs px-2 py-1 rounded-full" :class="getZScoreBadgeClass(measurement.height_for_age_zscore)">
            Z: {{ formatZScore(measurement.height_for_age_zscore) }}
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
      <span v-if="measurement.nutritional_status" class="text-xs px-3 py-1 rounded-full" :class="getZScoreTextClass(measurement.weight_for_age_zscore)">
        {{ measurement.nutritional_status }}
      </span>
      <span v-if="measurement.height_status" class="text-xs px-3 py-1 rounded-full" :class="getZScoreTextClass(measurement.height_for_age_zscore)">
        Tinggi: {{ measurement.height_status }}
      </span>
      <span v-if="measurement.weight_for_height_status" class="text-xs px-3 py-1 rounded-full" :class="getZScoreTextClass(measurement.weight_for_height_zscore)">
        {{ measurement.weight_for_height_status }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { isValidNumber, formatZScore, safeFormatDate } from '~/composables/useSafeData'
import { getZScoreBadgeClass, getZScoreTextClass, getZScoreBorderColor } from '~/utils/validation'

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

// Use utility function for border color
const borderColor = computed(() => {
  if (!props.measurement) return 'border-gray-500'
  return getZScoreBorderColor(
    props.measurement.weight_for_age_zscore,
    props.measurement.height_for_age_zscore
  )
})

// Use utility function for date formatting
const formatDate = safeFormatDate
</script>
