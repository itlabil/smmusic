import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getStreamUrl, recordPlay } from '@/services/songs'

export const usePlayerStore = defineStore('player', () => {
  const queue = ref([])
  const originalQueue = ref([])
  const currentIndex = ref(-1)
  const isPlaying = ref(false)
  const highQuality = ref(false)
  const isShuffled = ref(false)
  const repeatMode = ref('off') // 'off' | 'all' | 'one'
  const isQueueOpen = ref(false)
  const isFullScreenOpen = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const volume = ref(1)
  const audio = new Audio()

  const currentSong = computed(() =>
    currentIndex.value >= 0 ? queue.value[currentIndex.value] : null,
  )

  function playQueue(songs, startIndex = 0) {
    originalQueue.value = songs
    const startSong = songs[startIndex]

    if (isShuffled.value) {
      queue.value = shuffleExcluding(songs, startSong)
      currentIndex.value = 0
    } else {
      queue.value = songs
      currentIndex.value = startIndex
    }

    loadAndPlay()
  }

  function shuffleExcluding(songs, keepFirst) {
    const rest = songs.filter((s) => s.id !== keepFirst.id)
    for (let i = rest.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1))
      ;[rest[i], rest[j]] = [rest[j], rest[i]]
    }
    return [keepFirst, ...rest]
  }

  function toggleShuffle() {
    isShuffled.value = !isShuffled.value
    const current = currentSong.value
    if (!current) return

    if (isShuffled.value) {
      queue.value = shuffleExcluding(originalQueue.value, current)
      currentIndex.value = 0
    } else {
      queue.value = originalQueue.value
      currentIndex.value = originalQueue.value.findIndex((s) => s.id === current.id)
    }
  }

  function loadAndPlay() {
    const song = currentSong.value
    if (!song) return

    const quality = highQuality.value && song.flac_path ? 'high' : 'standard'
    audio.src = getStreamUrl(song.id, quality)
    audio.play()
    isPlaying.value = true
    recordPlay(song.id).catch(() => {}) // fire-and-forget
  }

  function togglePlay() {
    if (isPlaying.value) {
      audio.pause()
      isPlaying.value = false
    } else {
      audio.play()
      isPlaying.value = true
    }
  }

  function next() {
    if (currentIndex.value < queue.value.length - 1) {
      currentIndex.value++
      loadAndPlay()
    } else if (repeatMode.value === 'all') {
      currentIndex.value = 0
      loadAndPlay()
    }
  }

  function cycleRepeatMode() {
    if (repeatMode.value === 'off') repeatMode.value = 'all'
    else if (repeatMode.value === 'all') repeatMode.value = 'one'
    else repeatMode.value = 'off'
  }

  function toggleQueuePanel() {
    isQueueOpen.value = !isQueueOpen.value
  }

  function toggleFullScreen() {
    isFullScreenOpen.value = !isFullScreenOpen.value
  }

  function jumpTo(index) {
    currentIndex.value = index
    loadAndPlay()
  }

  function prev() {
    if (currentIndex.value > 0) {
      currentIndex.value--
      loadAndPlay()
    }
  }

  function toggleHighQuality() {
    highQuality.value = !highQuality.value
    if (currentSong.value) {
      const currentTime = audio.currentTime
      loadAndPlay()
      audio.currentTime = currentTime
    }
  }

  audio.addEventListener('ended', () => {
    if (repeatMode.value === 'one') {
      audio.currentTime = 0
      audio.play()
    } else {
      next()
    }
  })
  audio.addEventListener('timeupdate', () => {
    currentTime.value = audio.currentTime
  })
  audio.addEventListener('loadedmetadata', () => {
    duration.value = audio.duration
  })

  function seek(time) {
    audio.currentTime = time
    currentTime.value = time
  }

  function setVolume(value) {
    volume.value = value
    audio.volume = value
  }

  return {
    queue,
    currentIndex,
    currentSong,
    isPlaying,
    highQuality,
    isShuffled,
    repeatMode,
    isQueueOpen,
    isFullScreenOpen,
    currentTime,
    duration,
    volume,
    audio,
    playQueue,
    togglePlay,
    next,
    prev,
    toggleHighQuality,
    toggleShuffle,
    cycleRepeatMode,
    toggleQueuePanel,
    toggleFullScreen,
    jumpTo,
    seek,
    setVolume,
  }
})