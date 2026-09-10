<script setup>
import { usePlayerStore } from '@/stores/player'

const player = usePlayerStore()
</script>

<template>
  <div
    v-if="player.currentSong"
    class="h-20 bg-neutral-900 border-t border-neutral-800 flex items-center px-4 text-white flex-shrink-0"
  >
    <div class="flex items-center gap-3 w-64">
      <div class="w-12 h-12 bg-neutral-700 rounded flex-shrink-0"></div>
      <div class="min-w-0">
        <p class="text-sm font-medium truncate">{{ player.currentSong.title }}</p>
        <p class="text-xs text-neutral-400 truncate">{{ player.currentSong.artist }}</p>
      </div>
    </div>

    <div class="flex-1 flex flex-col items-center justify-center gap-2">
      <div class="flex items-center gap-4">
        <button @click="player.prev" class="text-neutral-400 hover:text-white">⏮</button>
        <button
          @click="player.togglePlay"
          class="w-8 h-8 bg-white rounded-full flex items-center justify-center text-black"
        >
          {{ player.isPlaying ? '⏸' : '▶' }}
        </button>
        <button @click="player.next" class="text-neutral-400 hover:text-white">⏭</button>
      </div>
    </div>

    <div class="w-64 flex justify-end items-center gap-2">
      <button
        @click="player.toggleHighQuality"
        :class="[
          'text-xs px-2 py-1 rounded border',
          player.highQuality
            ? 'border-green-500 text-green-500'
            : 'border-neutral-600 text-neutral-400',
        ]"
      >
        FLAC
      </button>
    </div>
  </div>
</template>