/**
 * Composable untuk lifecycle management yang reusable
 * Mengurangi boilerplate di setiap component
 */

import { ref, onMounted, onUnmounted, type WatchStopHandle } from 'vue'

export const useComponentLifecycle = () => {
  const isMounted = ref(false)
  const isInitialized = ref(false)
  const watchers: WatchStopHandle[] = []

  /**
   * Register watcher untuk auto-cleanup
   */
  const registerWatcher = (stopHandle: WatchStopHandle) => {
    watchers.push(stopHandle)
    return stopHandle
  }

  /**
   * Guard untuk async operations
   */
  const guardAsync = async <T>(
    operation: () => Promise<T>,
    onError?: (error: Error) => void
  ): Promise<T | null> => {
    if (!isMounted.value) return null
    
    try {
      const result = await operation()
      
      // Check again after async operation
      if (!isMounted.value) return null
      
      return result
    } catch (error) {
      if (isMounted.value && onError) {
        onError(error as Error)
      }
      return null
    }
  }

  /**
   * Guard untuk sync operations
   */
  const guardSync = <T>(operation: () => T, fallback: T): T => {
    if (!isMounted.value) return fallback
    try {
      return operation()
    } catch (error) {
      if (isMounted.value) {
        console.error('Error in guarded operation:', error)
      }
      return fallback
    }
  }

  onMounted(() => {
    isMounted.value = true
  })

  onUnmounted(() => {
    isMounted.value = false
    isInitialized.value = false
    
    // Cleanup all watchers
    watchers.forEach(stop => {
      if (stop) {
        try {
          stop()
        } catch (error) {
          console.warn('Error stopping watcher:', error)
        }
      }
    })
    watchers.length = 0
  })

  return {
    isMounted,
    isInitialized,
    registerWatcher,
    guardAsync,
    guardSync,
    setInitialized: (value: boolean) => {
      isInitialized.value = value
    }
  }
}

