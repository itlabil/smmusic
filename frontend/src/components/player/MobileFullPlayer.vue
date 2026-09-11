<script setup>
import {
  ChevronDown, Shuffle, SkipBack, SkipForward, Play, Pause,
  Repeat, Repeat1, Heart, Volume2, Volume1, VolumeX, Music,
} from 'lucide-vue-next'
import { usePlayerStore } from '@/stores/player'
import { likeSong, unlikeSong, getCoverUrl } from '@/services/songs'
import { formatDuration } from '@/utils/format'

const player = usePlayerStore()

function handleSeek(e) {
  player.seek(Number(e.target.value))
}

function handleVolumeChange(e) {
  player.setVolume(Number(e.target.value))
}

async function handleToggleLike() {
  const song = player.currentSong
  if (!song) return
  if (song.is_liked) {
    await unlikeSong(song.id)
  } else {
    await likeSong(song.id)
  }
  song.is_liked = !song.is_liked
}
</script>

<template>
  <div
    v-if="player.currentSong"
    class="fixed inset-0 z-50 sm:hidden bg-gradient-to-b from-neutral-700 to-neutral-950 flex flex-col text-white p-6"
  >
    <div class="flex items-center justify-between mb-8 flex-shrink-0">
      <button @click="player.toggleFullScreen" class="p-2 -ml-2">
        <ChevronDown :size="24" />
      </button>
      <p class="text-xs text-neutral-300 uppercase tracking-wide">Now Playing</p>
      <div class="w-8"></div>
    </div>

    <div class="flex-1 flex flex-col items-center justify-center overflow-y-auto">
      <div class="w-full max-w-xs aspect-square bg-neutral-800 rounded-lg shadow-2xl overflow-hidden flex items-center justify-center mb-8 flex-shrink-0">
        <img
          v-if="getCoverUrl(player.currentSong)"
          :src="getCoverUrl(player.currentSong)"
          class="w-full h-full object-cover"
          alt=""
        />
        <Music v-else :size="64" class="text-neutral-600" />
      </div>

      <div class="w-full max-w-xs flex items-center justify-between mb-6">
        <div class="min-w-0">
          <p class="text-white text-xl font-bold truncate">{{ player.currentSong.title }}</p>
          <p class="text-neutral-400 text-sm truncate">{{ player.currentSong.artist }}</p>
        </div>
        <button
          @click="handleToggleLike"
          :class="['p-2 flex-shrink-0', player.currentSong.is_liked ? 'text-green-500' : 'text-neutral-400']"
        >
          <Heart :size="22" :fill="player.currentSong.is_liked ? 'currentColor' : 'none'" />
        </button>
      </div>

      <div class="w-full max-w-xs mb-4">
        <input
          type="range"
          min="0"
          :max="player.duration || 0"
          :value="player.currentTime"
          @input="handleSeek"
          class="w-full h-1 accent-white cursor-pointer"
        />
        <div class="flex justify-between text-xs text-neutral-400 mt-1">
          <span>{{ formatDuration(player.currentTime) }}</span>
          <span>{{ formatDuration(player.duration) }}</span>
        </div>
      </div>

      <div class="w-full max-w-xs flex items-center justify-between mb-8">
        <button
          @click="player.toggleShuffle"
          :class="player.isShuffled ? 'text-green-500' : 'text-neutral-400'"
        >
          <Shuffle :size="20" />
        </button>
        <button @click="player.prev" class="text-white">
          <SkipBack :size="28" fill="currentColor" />
        </button>
        <button
          @click="player.togglePlay"
          class="w-14 h-14 bg-white rounded-full flex items-center justify-center text-black"
        >
          <Pause v-if="player.isPlaying" :size="24" fill="currentColor" />
          <Play v-else :size="24" fill="currentColor" class="ml-1" />
        </button>
        <button @click="player.next" class="text-white">
          <SkipForward :size="28" fill="currentColor" />
        </button>
        <button
          @click="player.cycleRepeatMode"
          :class="player.repeatMode !== 'off' ? 'text-green-500' : 'text-neutral-400'"
        >
          <Repeat1 v-if="player.repeatMode === 'one'" :size="20" />
          <Repeat v-else :size="20" />
        </button>
      </div>

      <div class="w-full max-w-xs flex items-center gap-2">
        <VolumeX v-if="player.volume === 0" :size="18" class="text-neutral-400" />
        <Volume1 v-else-if="player.volume < 0.5" :size="18" class="text-neutral-400" />
        <Volume2 v-else :size="18" class="text-neutral-400" />
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          :value="player.volume"
          @input="handleVolumeChange"
          class="flex-1 h-1 accent-white cursor-pointer"
        />
      </div>
    </div>
  </div>
</template>