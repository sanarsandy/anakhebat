<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-all p-4">
    <div class="flex items-start gap-4">
      <!-- Thumbnail or Icon -->
      <div class="flex-shrink-0">
        <div v-if="recommendation.content.thumbnail_url" class="w-20 h-20 rounded-lg overflow-hidden bg-gray-100">
          <img :src="recommendation.content.thumbnail_url" :alt="recommendation.content.title" class="w-full h-full object-cover" />
        </div>
        <div v-else class="w-20 h-20 rounded-lg bg-gradient-to-br flex items-center justify-center" :class="getCategoryColor(recommendation.content.category)">
          <span class="text-3xl">{{ getCategoryIcon(recommendation.content.category) }}</span>
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 min-w-0">
        <!-- Priority Badge -->
        <div class="flex items-center gap-2 mb-2">
          <span v-if="recommendation.priority === 'high'" class="px-2 py-0.5 bg-red-100 text-red-700 text-xs font-semibold rounded">
            Prioritas Tinggi
          </span>
          <span v-else-if="recommendation.priority === 'medium'" class="px-2 py-0.5 bg-amber-100 text-amber-700 text-xs font-semibold rounded">
            Prioritas Sedang
          </span>
          <span class="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded" :class="getCategoryBadgeColor(recommendation.content.category)">
            {{ getCategoryName(recommendation.content.category) }}
          </span>
          <span class="px-2 py-0.5 bg-blue-100 text-blue-700 text-xs rounded">
            {{ recommendation.content.content_type === 'video' ? 'Video' : 'Artikel' }}
          </span>
        </div>

        <!-- Title -->
        <h3 class="font-bold text-gray-900 text-sm mb-1 line-clamp-2">
          {{ recommendation.content.title }}
        </h3>

        <!-- Description -->
        <p v-if="recommendation.content.description" class="text-xs text-gray-600 mb-2 line-clamp-2">
          {{ recommendation.content.description }}
        </p>

        <!-- Reason -->
        <p class="text-xs text-gray-500 mb-3 italic">
          {{ recommendation.reason }}
        </p>

        <!-- Action -->
        <a 
          :href="recommendation.content.url" 
          target="_blank" 
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 text-jurnal-teal-600 hover:text-jurnal-teal-700 text-xs font-semibold transition"
        >
          <span>{{ recommendation.content.content_type === 'video' ? 'Tonton Video' : 'Baca Artikel' }}</span>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
          </svg>
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  recommendation: {
    type: Object,
    required: true
  }
})

const getCategoryName = (category) => {
  const names = {
    'sensory': 'Sensorik',
    'motor': 'Motorik',
    'perception': 'Persepsi',
    'cognitive': 'Kognitif'
  }
  return names[category] || category
}

const getCategoryIcon = (category) => {
  const icons = {
    'sensory': '👁️',
    'motor': '🤲',
    'perception': '🎨',
    'cognitive': '🧠'
  }
  return icons[category] || '📚'
}

const getCategoryColor = (category) => {
  const colors = {
    'sensory': 'from-purple-200 to-purple-300',
    'motor': 'from-blue-200 to-blue-300',
    'perception': 'from-green-200 to-green-300',
    'cognitive': 'from-orange-200 to-orange-300'
  }
  return colors[category] || 'from-gray-200 to-gray-300'
}

const getCategoryBadgeColor = (category) => {
  const colors = {
    'sensory': 'text-purple-700 bg-purple-50',
    'motor': 'text-blue-700 bg-blue-50',
    'perception': 'text-green-700 bg-green-50',
    'cognitive': 'text-orange-700 bg-orange-50'
  }
  return colors[category] || 'text-gray-700 bg-gray-50'
}
</script>

