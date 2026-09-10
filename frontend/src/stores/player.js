import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getStreamUrl, recordPlay } from '@/services/songs'

export const usePlayerStore = defineStore('player', () => {
  const queue = ref([])
  const currentIndex = ref(-1)
  const isPlaying = ref(false)
  const highQuality = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const volume = ref(1)
  const audio = new Audio()

  const currentSong = computed(() =>
    currentIndex.value >= 0 ? queue.value[currentIndex.value] : null,
  )

  function playQueue(songs, startIndex = 0) {
    queue.value = songs
    currentIndex.value = startIndex
    loadAndPlay()
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
    }
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

  audio.addEventListener('ended', next)
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
    currentSong,
    isPlaying,
    highQuality,
    currentTime,
    duration,
    volume,
    audio,
    playQueue,
    togglePlay,
    next,
    prev,
    toggleHighQuality,
    seek,
    setVolume,
  }
})