<template>
  <div>
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Grafik Pertumbuhan</h1>
        <p class="text-gray-600 mt-2">Visualisasi pertumbuhan anak berdasarkan standar WHO</p>
      </div>
      <NuxtLink to="/growth" class="px-4 py-2 text-gray-700 hover:text-gray-900 font-medium">
        ← Kembali ke Daftar
      </NuxtLink>
    </div>

    <!-- No Child Selected -->
    <div v-if="!childStore.selectedChild" class="bg-amber-50 border border-amber-200 rounded-xl p-8 text-center">
      <div class="text-5xl mb-4">⚠️</div>
      <h2 class="text-xl font-bold text-gray-900 mb-2">Pilih Anak Terlebih Dahulu</h2>
      <p class="text-gray-600 mb-4">Silakan pilih anak dari dropdown di header untuk melihat grafik pertumbuhan</p>
      <NuxtLink to="/children" class="inline-block px-6 py-3 bg-indigo-600 text-white font-semibold rounded-lg hover:bg-indigo-700 transition">
        Kelola Profil Anak
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-xl shadow-sm p-12 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
      <p class="text-gray-600">Memuat data grafik...</p>
    </div>

    <!-- Charts -->
    <div v-else-if="childStore.selectedChild" class="space-y-8">
      <!-- Chart Tabs -->
      <div class="bg-white rounded-xl shadow-sm p-4">
        <div class="flex space-x-2 border-b border-gray-200">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'px-6 py-3 font-medium text-sm transition-colors',
              activeTab === tab.id
                ? 'text-indigo-600 border-b-2 border-indigo-600'
                : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- Weight-for-Age Chart -->
      <div v-if="activeTab === 'wfa'" class="bg-white rounded-xl shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-4">Berat Badan menurut Umur (BB/U)</h2>
        <div class="h-96">
          <canvas ref="wfaChartRef"></canvas>
        </div>
      </div>

      <!-- Height-for-Age Chart -->
      <div v-if="activeTab === 'hfa'" class="bg-white rounded-xl shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-4">Tinggi Badan menurut Umur (TB/U)</h2>
        <div class="h-96">
          <canvas ref="hfaChartRef"></canvas>
        </div>
      </div>


      <!-- Legend -->
      <div class="bg-white rounded-xl shadow-sm p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Keterangan</h3>
        <div class="grid md:grid-cols-4 gap-4">
          <div class="flex items-center space-x-2">
            <div class="w-4 h-4 bg-red-300"></div>
            <span class="text-sm text-gray-600">-3 SD (Sangat Kurang)</span>
          </div>
          <div class="flex items-center space-x-2">
            <div class="w-4 h-4 bg-orange-300"></div>
            <span class="text-sm text-gray-600">-2 SD (Kurang)</span>
          </div>
          <div class="flex items-center space-x-2">
            <div class="w-4 h-4 bg-yellow-300"></div>
            <span class="text-sm text-gray-600">-1 SD (Normal Bawah)</span>
          </div>
          <div class="flex items-center space-x-2">
            <div class="w-4 h-4 bg-green-300"></div>
            <span class="text-sm text-gray-600">Median (Normal)</span>
          </div>
        </div>
        <div class="mt-4 pt-4 border-t border-gray-200">
          <div class="flex items-center space-x-2">
            <div class="w-4 h-4 border-2 border-indigo-600 bg-indigo-50"></div>
            <span class="text-sm font-medium text-gray-900">Data Pengukuran Anak</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

// Chart.js will be imported dynamically in onMounted
let Chart = null

const childStore = useChildStore()
const measurementStore = useMeasurementStore()
const apiBase = useApiUrl()
const authStore = useAuthStore()

const loading = ref(true)
const activeTab = ref('wfa')
const tabs = [
  { id: 'wfa', label: 'Berat Badan / Umur' },
  { id: 'hfa', label: 'Tinggi Badan / Umur' }
]

const wfaChartRef = ref(null)
const hfaChartRef = ref(null)

let wfaChart = null
let hfaChart = null

const fetchWHOStandards = async (indicator, gender) => {
  try {
    const params = {
      indicator,
      gender
    }
    
    params.minAge = 0
    params.maxAge = 60
    
    const response = await $fetch(`${apiBase}/api/who-standards`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      params
    })
    return response.standards || []
  } catch (error) {
    console.error(`Failed to fetch WHO standards for ${indicator}:`, error)
    return []
  }
}

const calculateAgeInMonths = (dob) => {
  const birthDate = new Date(dob)
  const today = new Date()
  const months = (today.getFullYear() - birthDate.getFullYear()) * 12 + 
                 (today.getMonth() - birthDate.getMonth())
  return months
}

const createChart = (canvasRef, chartType, indicator, label, yAxisLabel) => {
  if (!canvasRef || !childStore.selectedChild || !Chart) {
    console.warn('Cannot create chart: missing requirements', { canvasRef: !!canvasRef, child: !!childStore.selectedChild, Chart: !!Chart })
    return
  }

  // Destroy existing chart if any
  if (chartType === 'wfa' && wfaChart) {
    try {
      wfaChart.destroy()
    } catch (e) {
      console.warn('Error destroying wfa chart:', e)
    }
    wfaChart = null
  }
  if (chartType === 'hfa' && hfaChart) {
    try {
      hfaChart.destroy()
    } catch (e) {
      console.warn('Error destroying hfa chart:', e)
    }
    hfaChart = null
  }

  const measurements = measurementStore.measurements || []
  const gender = childStore.selectedChild.gender === 'L' || childStore.selectedChild.gender === 'laki-laki' ? 'male' : 'female'
  const dob = new Date(childStore.selectedChild.dob)

  // Prepare measurement data
  const measurementData = measurements.map(m => {
    const measurementDate = new Date(m.measurement_date)
    const ageMonths = ((measurementDate.getFullYear() - dob.getFullYear()) * 12) + 
                      (measurementDate.getMonth() - dob.getMonth())
    
    if (indicator === 'wfa') {
      return { x: ageMonths, y: m.weight }
    } else if (indicator === 'hfa') {
      return { x: ageMonths, y: m.height }
    }
    return null
  }).filter(Boolean)

  // Fetch WHO standards and create chart
  fetchWHOStandards(indicator, gender).then(standards => {
    // Double-check canvas ref is still available (might have been destroyed)
    if (!canvasRef || !Chart) {
      console.warn('Canvas ref or Chart not available when creating chart after fetch')
      return
    }
    
    // Log standards data for debugging
    console.log(`WHO standards for ${indicator} (${gender}):`, standards.slice(0, 3))
    
    // Filter out standards with zero SD values (shouldn't happen, but just in case)
    const validStandards = standards.filter(s => s.sd0 > 0)
    if (validStandards.length === 0) {
      console.error(`No valid WHO standards found for ${indicator}!`)
      return
    }
    
    const xValues = validStandards.map(s => s.x_value)
    const datasets = [
      // WHO Standard Curves
      {
        label: '-3 SD',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd3neg })),
        borderColor: 'rgb(239, 68, 68)',
        backgroundColor: 'rgba(239, 68, 68, 0.1)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false
      },
      {
        label: '-2 SD',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd2neg })),
        borderColor: 'rgb(251, 146, 60)',
        backgroundColor: 'rgba(251, 146, 60, 0.1)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false
      },
      {
        label: '-1 SD',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd1neg })),
        borderColor: 'rgb(253, 224, 71)',
        backgroundColor: 'rgba(253, 224, 71, 0.1)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false
      },
      {
        label: 'Median',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd0 })),
        borderColor: 'rgb(34, 197, 94)',
        backgroundColor: 'rgba(34, 197, 94, 0.1)',
        borderWidth: 2,
        pointRadius: 0,
        fill: false
      },
      {
        label: '+1 SD',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd1 })),
        borderColor: 'rgb(253, 224, 71)',
        backgroundColor: 'rgba(253, 224, 71, 0.1)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false
      },
      {
        label: '+2 SD',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd2 })),
        borderColor: 'rgb(251, 146, 60)',
        backgroundColor: 'rgba(251, 146, 60, 0.1)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false
      },
      {
        label: '+3 SD',
        data: validStandards.map(s => ({ x: s.x_value, y: s.sd3 })),
        borderColor: 'rgb(239, 68, 68)',
        backgroundColor: 'rgba(239, 68, 68, 0.1)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0,
        fill: false
      },
      // Child's measurements
      {
        label: 'Data Pengukuran',
        data: measurementData,
        borderColor: 'rgb(99, 102, 241)',
        backgroundColor: 'rgba(99, 102, 241, 0.5)',
        borderWidth: 3,
        pointRadius: 6,
        pointHoverRadius: 8,
        pointBackgroundColor: 'rgb(99, 102, 241)',
        pointBorderColor: '#fff',
        pointBorderWidth: 2,
        fill: false,
        tension: 0.4
      }
    ]

    const chartConfig = {
      type: 'line',
      data: { datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: 'index',
          intersect: false
        },
        plugins: {
          legend: {
            display: true,
            position: 'top'
          },
          tooltip: {
            enabled: true
          }
        },
        scales: {
          x: {
            title: {
              display: true,
                  text: 'Umur (bulan)'
            },
            type: 'linear',
            position: 'bottom'
          },
          y: {
            title: {
              display: true,
              text: yAxisLabel
            }
          }
        }
      }
    }

    // Final check before creating chart
    if (!canvasRef || !Chart) {
      console.warn('Canvas ref or Chart not available right before chart creation')
      return
    }
    
    try {
      const chart = new Chart(canvasRef, chartConfig)
      
      if (chartType === 'wfa') wfaChart = chart
      if (chartType === 'hfa') hfaChart = chart
      
      console.log(`Chart created successfully for ${chartType}`)
    } catch (error) {
      console.error(`Error creating chart for ${chartType}:`, error)
    }
  }).catch(error => {
    console.error(`Error fetching WHO standards for ${indicator}:`, error)
  })
}

// Helper function to wait for canvas ref to be available
const waitForCanvasRef = async (ref, maxAttempts = 10) => {
  for (let i = 0; i < maxAttempts; i++) {
    if (ref.value) {
      return ref.value
    }
    await new Promise(resolve => setTimeout(resolve, 50))
  }
  return null
}

onMounted(async () => {
  isMounted.value = true
  
  // Import Chart.js dynamically on client side
  if (process.client) {
    const chartModule = await import('chart.js')
    Chart = chartModule.Chart
    Chart.register(...chartModule.registerables)
  }

  // Guard: Check if still mounted after async import
  if (!isMounted.value) return

  if (childStore.selectedChild) {
    loading.value = true
    try {
      await measurementStore.fetchMeasurements(childStore.selectedChild.id)
      
      // Guard: Check if still mounted after async operation
      if (!isMounted.value) return
      
      // Wait for DOM to update
      await nextTick()
      await nextTick() // Double nextTick to ensure v-if has rendered
      
      // Guard: Check again before creating charts
      if (!isMounted.value) return
      
      // Wait for canvas ref to be available
      let canvasRef = null
      if (activeTab.value === 'wfa') {
        canvasRef = await waitForCanvasRef(wfaChartRef)
        if (canvasRef && isMounted.value) {
          createChart(canvasRef, 'wfa', 'wfa', 'Berat Badan menurut Umur', 'Berat Badan (kg)')
        }
      } else if (activeTab.value === 'hfa') {
        canvasRef = await waitForCanvasRef(hfaChartRef)
        if (canvasRef && isMounted.value) {
          createChart(canvasRef, 'hfa', 'hfa', 'Tinggi Badan menurut Umur', 'Tinggi Badan (cm)')
        }
      }
    } catch (error) {
      if (isMounted.value) {
        console.error('Error loading charts:', error)
      }
    } finally {
      if (isMounted.value) {
        loading.value = false
      }
    }
  } else {
    loading.value = false
  }
})

// Function to initialize chart for active tab
const initializeChartForActiveTab = async () => {
  if (!Chart || !childStore.selectedChild || loading.value) return
  
  const newTab = activeTab.value
  
  // Wait for DOM to update (v-if needs to render the canvas)
  await nextTick()
  await nextTick() // Double nextTick to ensure v-if has rendered
  
  // Wait for canvas ref to be available
  let canvasRef = null
  if (newTab === 'wfa') {
    if (wfaChart) {
      try {
        wfaChart.destroy()
      } catch (e) {
        console.warn('Error destroying wfa chart:', e)
      }
      wfaChart = null
    }
    canvasRef = await waitForCanvasRef(wfaChartRef)
    if (canvasRef && !wfaChart) {
      createChart(canvasRef, 'wfa', 'wfa', 'Berat Badan menurut Umur', 'Berat Badan (kg)')
    }
  } else if (newTab === 'hfa') {
    if (hfaChart) {
      try {
        hfaChart.destroy()
      } catch (e) {
        console.warn('Error destroying hfa chart:', e)
      }
      hfaChart = null
    }
    canvasRef = await waitForCanvasRef(hfaChartRef)
    if (canvasRef && !hfaChart) {
      createChart(canvasRef, 'hfa', 'hfa', 'Tinggi Badan menurut Umur', 'Tinggi Badan (cm)')
    }
  }
}

// Track if component is mounted
const isMounted = ref(false)

// Watch for tab changes - always recreate chart when tab becomes active
const stopTabWatcher = watch(activeTab, async (newTab) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  await initializeChartForActiveTab()
}, { immediate: false })

// Watch for when loading finishes - this ensures chart is created after initial load
const stopLoadingWatcher = watch(() => loading.value, async (isLoading) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  
  if (!isLoading && Chart && childStore.selectedChild) {
    // Small delay to ensure DOM is fully ready
    await nextTick()
    
    // Guard again after async operation
    if (!isMounted.value) return
    
    await initializeChartForActiveTab()
  }
}, { immediate: false })

// Watch for child changes
const stopChildWatcher = watch(() => childStore.selectedChild, async (newChild) => {
  // Guard: Don't update if component is unmounted
  if (!isMounted.value) return
  
  if (newChild) {
    loading.value = true
    try {
      await measurementStore.fetchMeasurements(newChild.id)
      await nextTick()
      
      // Guard again after async operation
      if (!isMounted.value) return
      
      // Recreate all charts
      if (wfaChart) {
        try {
          wfaChart.destroy()
        } catch (e) {
          console.warn('Error destroying wfa chart:', e)
        }
      }
      if (hfaChart) {
        try {
          hfaChart.destroy()
        } catch (e) {
          console.warn('Error destroying hfa chart:', e)
        }
      }
      
      wfaChart = null
      hfaChart = null
      
      // Guard: Check again before creating charts
      if (!isMounted.value) return
      
      if (wfaChartRef.value) {
        createChart(wfaChartRef.value, 'wfa', 'wfa', 'Berat Badan menurut Umur', 'Berat Badan (kg)')
      }
      if (hfaChartRef.value) {
        createChart(hfaChartRef.value, 'hfa', 'hfa', 'Tinggi Badan menurut Umur', 'Tinggi Badan (cm)')
      }
    } catch (error) {
      if (isMounted.value) {
        console.error('Error reloading charts:', error)
      }
    } finally {
      if (isMounted.value) {
        loading.value = false
      }
    }
  }
}, { immediate: false })

onMounted(() => {
  isMounted.value = true
})

onUnmounted(() => {
  isMounted.value = false
  
  // Stop all watchers
  stopTabWatcher()
  stopLoadingWatcher()
  stopChildWatcher()
  
  // Destroy charts
  if (wfaChart) {
    try {
      wfaChart.destroy()
    } catch (e) {
      console.warn('Error destroying wfa chart on unmount:', e)
    }
  }
  if (hfaChart) {
    try {
      hfaChart.destroy()
    } catch (e) {
      console.warn('Error destroying hfa chart on unmount:', e)
    }
  }
})
</script>


