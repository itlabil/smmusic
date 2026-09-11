<script setup>
import {
  ChevronDown, Shuffle, SkipBack, SkipForward, Play, Pause,
  Repeat, Repeat1, Heart, Volume2, Volume1, VolumeX, Music, ListMusic,
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

function openQueueFromFullScreen() {
  player.toggleFullScreen()
  player.toggleQueuePanel()
}
</script>

<template>
  <div
    v-if="player.currentSong"
    class="fixed inset-x-0 top-0 bottom-16 z-40 sm:hidden bg-gradient-to-b from-neutral-700 to-neutral-950 flex flex-col text-white p-4 overflow-hidden"
  >
    <div class="flex items-center justify-between mb-2 flex-shrink-0">
      <button @click="player.toggleFullScreen" class="p-2 -ml-2">
        <ChevronDown :size="22" />
      </button>
      <p class="text-xs text-neutral-300 uppercase tracking-wide">Now Playing</p>
      <div class="w-8"></div>
    </div>

    <div class="flex-1 flex flex-col items-center justify-center min-h-0">
      <div class="w-full max-w-[220px] aspect-square bg-neutral-800 rounded-lg shadow-2xl overflow-hidden flex items-center justify-center mb-4 flex-shrink-0">
        <img
          v-if="getCoverUrl(player.currentSong)"
          :src="getCoverUrl(player.currentSong)"
          class="w-full h-full object-cover"
          alt=""
        />
        <Music v-else :size="56" class="text-neutral-600" />
      </div>

      <div class="w-full max-w-xs flex items-center justify-between mb-3 flex-shrink-0">
        <div class="min-w-0">
          <p class="text-white text-lg font-bold truncate">{{ player.currentSong.title }}</p>
          <p class="text-neutral-400 text-sm truncate">{{ player.currentSong.artist }}</p>
        </div>
        <button
          @click="handleToggleLike"
          :class="['p-2 flex-shrink-0', player.currentSong.is_liked ? 'text-green-500' : 'text-neutral-400']"
        >
          <Heart :size="20" :fill="player.currentSong.is_liked ? 'currentColor' : 'none'" />
        </button>
      </div>

      <div class="w-full max-w-xs mb-3 flex-shrink-0">
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

      <div class="w-full max-w-xs flex items-center justify-between mb-3 flex-shrink-0">
        <button
          @click="player.toggleShuffle"
          :class="player.isShuffled ? 'text-green-500' : 'text-neutral-400'"
        >
          <Shuffle :size="18" />
        </button>
        <button @click="player.prev" class="text-white">
          <SkipBack :size="24" fill="currentColor" />
        </button>
        <button
          @click="player.togglePlay"
          class="w-12 h-12 bg-white rounded-full flex items-center justify-center text-black"
        >
          <Pause v-if="player.isPlaying" :size="20" fill="currentColor" />
          <Play v-else :size="20" fill="currentColor" class="ml-0.5" />
        </button>
        <button @click="player.next" class="text-white">
          <SkipForward :size="24" fill="currentColor" />
        </button>
        <button
          @click="player.cycleRepeatMode"
          :class="player.repeatMode !== 'off' ? 'text-green-500' : 'text-neutral-400'"
        >
          <Repeat1 v-if="player.repeatMode === 'one'" :size="18" />
          <Repeat v-else :size="18" />
        </button>
      </div>

      <div class="w-full max-w-xs flex items-center gap-3 mt-6 flex-shrink-0">
        <button
          @click="player.toggleHighQuality"
          :class="[
            'text-xs px-2 py-1 rounded border transition flex-shrink-0',
            player.highQuality
              ? 'border-green-500 text-green-500'
              : 'border-neutral-600 text-neutral-400',
          ]"
        >
          FLAC
        </button>
        <button
          @click="openQueueFromFullScreen"
          :class="['flex-shrink-0', player.isQueueOpen ? 'text-green-500' : 'text-neutral-400']"
        >
          <ListMusic :size="18" />
        </button>

        <VolumeX v-if="player.volume === 0" :size="16" class="text-neutral-400 flex-shrink-0" />
        <Volume1 v-else-if="player.volume < 0.5" :size="16" class="text-neutral-400 flex-shrink-0" />
        <Volume2 v-else :size="16" class="text-neutral-400 flex-shrink-0" />
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