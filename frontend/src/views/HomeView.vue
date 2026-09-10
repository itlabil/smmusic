<script setup>
import { onMounted, ref } from 'vue'
import { listSongs, likeSong, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import SongRow from '@/components/song-card/SongRow.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'

const songs = ref([])
const player = usePlayerStore()
const showModal = ref(false)
const selectedSongId = ref(null)

onMounted(async () => {
  const res = await listSongs()
  songs.value = res.data.songs || []
})

function playSong(index) {
  player.playQueue(songs.value, index)
}

async function handleToggleLike(song) {
  if (song.is_liked) {
    await unlikeSong(song.id)
  } else {
    await likeSong(song.id)
  }
  song.is_liked = !song.is_liked
}

function openAddToPlaylist(songId) {
  selectedSongId.value = songId
  showModal.value = true
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
        @toggle-like="handleToggleLike(song)"
        @add-to-playlist="openAddToPlaylist(song.id)"
      />
    </div>

    <AddToPlaylistModal
      v-if="showModal"
      :song-id="selectedSongId"
      @close="showModal = false"
    />
  </div>
</template>