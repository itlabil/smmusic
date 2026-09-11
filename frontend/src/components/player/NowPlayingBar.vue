<script setup>
import { SkipBack, SkipForward, Play, Pause, Volume2, Volume1, VolumeX, Music, Shuffle, Repeat, Repeat1, ListMusic } from 'lucide-vue-next'
import { usePlayerStore } from '@/stores/player'
import { formatDuration } from '@/utils/format'
import { getCoverUrl } from '@/services/songs'

const player = usePlayerStore()

function handleSeek(e) {
  player.seek(Number(e.target.value))
}

function handleVolumeChange(e) {
  player.setVolume(Number(e.target.value))
}
</script>

<template>
  <div
    v-if="player.currentSong"
    class="bg-neutral-900 border-t border-neutral-800 flex flex-col px-4 py-2 text-white flex-shrink-0"
  >
    <div class="flex items-center">
      <div class="flex items-center gap-3 w-32 sm:w-64">
        <div class="w-10 h-10 sm:w-12 sm:h-12 bg-neutral-700 rounded flex-shrink-0 overflow-hidden flex items-center justify-center">
          <img
            v-if="getCoverUrl(player.currentSong)"
            :src="getCoverUrl(player.currentSong)"
            class="w-full h-full object-cover"
            alt=""
          />
          <Music v-else :size="18" class="text-neutral-500" />
        </div>
        <div class="min-w-0 hidden sm:block">
          <p class="text-sm font-medium truncate">{{ player.currentSong.title }}</p>
          <p class="text-xs text-neutral-400 truncate">{{ player.currentSong.artist }}</p>
        </div>
      </div>

      <div class="flex-1 flex flex-col items-center justify-center gap-1 px-2">
        <div class="flex items-center gap-4">
          <button
            @click="player.toggleShuffle"
            :class="['transition', player.isShuffled ? 'text-green-500' : 'text-neutral-400 hover:text-white']"
            title="Shuffle"
          >
            <Shuffle :size="18" />
          </button>
          <button @click="player.prev" class="text-neutral-400 hover:text-white transition">
            <SkipBack :size="20" fill="currentColor" />
          </button>
          <button
            @click="player.togglePlay"
            class="w-8 h-8 bg-white rounded-full flex items-center justify-center text-black hover:scale-105 transition"
          >
            <Pause v-if="player.isPlaying" :size="16" fill="currentColor" />
            <Play v-else :size="16" fill="currentColor" class="ml-0.5" />
          </button>
          <button @click="player.next" class="text-neutral-400 hover:text-white transition">
            <SkipForward :size="20" fill="currentColor" />
          </button>
          <button
            @click="player.cycleRepeatMode"
            :class="['transition', player.repeatMode !== 'off' ? 'text-green-500' : 'text-neutral-400 hover:text-white']"
            title="Repeat"
          >
            <Repeat1 v-if="player.repeatMode === 'one'" :size="18" />
            <Repeat v-else :size="18" />
          </button>
        </div>

        <div class="w-full max-w-md flex items-center gap-2">
          <span class="text-xs text-neutral-400 w-9 text-right">{{ formatDuration(player.currentTime) }}</span>
          <input
            type="range"
            min="0"
            :max="player.duration || 0"
            :value="player.currentTime"
            @input="handleSeek"
            class="flex-1 h-1 accent-white cursor-pointer"
          />
          <span class="text-xs text-neutral-400 w-9">{{ formatDuration(player.duration) }}</span>
        </div>
      </div>

      <div class="w-16 sm:w-64 flex justify-end items-center gap-3">
        <button
          @click="player.toggleHighQuality"
          :class="[
            'text-xs px-2 py-1 rounded border transition flex-shrink-0',
            player.highQuality
              ? 'border-green-500 text-green-500'
              : 'border-neutral-600 text-neutral-400 hover:text-white hover:border-neutral-400',
          ]"
        >
          FLAC
        </button>

        <button
          @click="player.toggleQueuePanel"
          :class="['transition p-1 hidden sm:block', player.isQueueOpen ? 'text-green-500' : 'text-neutral-400 hover:text-white']"
          title="Queue"
        >
          <ListMusic :size="18" />
        </button>

        <div class="hidden sm:flex items-center gap-1.5">
          <VolumeX v-if="player.volume === 0" :size="16" class="text-neutral-400" />
          <Volume1 v-else-if="player.volume < 0.5" :size="16" class="text-neutral-400" />
          <Volume2 v-else :size="16" class="text-neutral-400" />
          <input
            type="range"
            min="0"
            max="1"
            step="0.01"
            :value="player.volume"
            @input="handleVolumeChange"
            class="w-20 h-1 accent-white cursor-pointer"
          />
        </div>
      </div>
    </div>
  </div>
</template>