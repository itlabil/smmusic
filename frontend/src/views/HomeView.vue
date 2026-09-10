<script setup>
import { onMounted, ref } from 'vue'
import { listSongs } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import SongRow from '@/components/song-card/SongRow.vue'

const songs = ref([])
const player = usePlayerStore()

onMounted(async () => {
  const res = await listSongs()
  songs.value = res.data.songs || []
})

function playSong(index) {
  player.playQueue(songs.value, index)
}
</script>

<template>
  <div class="p-8">
    <h1 class="text-white text-2xl font-bold mb-6">Your Library</h1>

    <div v-if="songs.length === 0" class="text-neutral-400">
      No songs yet. Upload your first song to get started.
    </div>

    <div v-else class="space-y-1">
      <SongRow
        v-for="(song, index) in songs"
        :key="song.id"
        :song="song"
        @play="playSong(index)"
      />
    </div>
  </div>
</template>