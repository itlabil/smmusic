<script setup>
import { X, Music } from 'lucide-vue-next'
import { usePlayerStore } from '@/stores/player'
import { getCoverUrl } from '@/services/songs'

const player = usePlayerStore()

const upcoming = () => player.queue.slice(player.currentIndex + 1)
</script>

<template>
  <aside class="fixed inset-0 z-40 sm:z-auto sm:static w-full sm:w-80 bg-neutral-950 sm:border-l border-neutral-800 h-full flex flex-col text-white flex-shrink-0">
    <div class="flex items-center justify-between p-4 border-b border-neutral-800">
      <h2 class="font-semibold">Queue</h2>
      <button @click="player.toggleQueuePanel" class="text-neutral-400 hover:text-white transition p-1">
        <X :size="18" />
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4">
      <div v-if="player.currentSong">
        <h3 class="text-neutral-400 text-xs font-semibold uppercase mb-2">Now Playing</h3>
        <div class="flex items-center gap-3 mb-6">
          <div class="w-10 h-10 bg-neutral-700 rounded flex-shrink-0 overflow-hidden flex items-center justify-center">
            <img
              v-if="getCoverUrl(player.currentSong)"
              :src="getCoverUrl(player.currentSong)"
              class="w-full h-full object-cover"
              alt=""
            />
            <Music v-else :size="16" class="text-neutral-500" />
          </div>
          <div class="min-w-0">
            <p class="text-green-500 text-sm font-medium truncate">{{ player.currentSong.title }}</p>
            <p class="text-neutral-400 text-xs truncate">{{ player.currentSong.artist }}</p>
          </div>
        </div>
      </div>

      <div v-if="upcoming().length > 0">
        <h3 class="text-neutral-400 text-xs font-semibold uppercase mb-2">Next Up</h3>
        <div class="space-y-1">
          <div
            v-for="(song, i) in upcoming()"
            :key="song.id + '-' + i"
            @click="player.jumpTo(player.currentIndex + 1 + i)"
            class="flex items-center gap-3 px-2 py-2 rounded hover:bg-neutral-800 cursor-pointer"
          >
            <div class="w-10 h-10 bg-neutral-700 rounded flex-shrink-0 overflow-hidden flex items-center justify-center">
              <img v-if="getCoverUrl(song)" :src="getCoverUrl(song)" class="w-full h-full object-cover" alt="" />
              <Music v-else :size="16" class="text-neutral-500" />
            </div>
            <div class="min-w-0">
              <p class="text-sm font-medium truncate">{{ song.title }}</p>
              <p class="text-neutral-400 text-xs truncate">{{ song.artist }}</p>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!player.currentSong" class="text-neutral-400 text-sm">
        Nothing is playing right now.
      </div>
    </div>
  </aside>
</template>