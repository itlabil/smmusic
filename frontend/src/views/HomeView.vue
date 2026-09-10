<script setup>
import { onMounted, ref } from 'vue'
import { listSongs, likeSong, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import SongRow from '@/components/song-card/SongRow.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'

const songs = ref([])
const player = usePlayerStore()
const showAddToPlaylistModal = ref(false)
const showEditModal = ref(false)
const selectedSongId = ref(null)
const selectedSong = ref(null)

async function loadSongs() {
  const res = await listSongs()
  songs.value = res.data.songs || []
}

onMounted(loadSongs)

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
  showAddToPlaylistModal.value = true
}

function openEdit(song) {
  selectedSong.value = song
  showEditModal.value = true
}
</script>

<template>
  <div class="p-4 md:p-8">
    <h1 class="text-white text-xl md:text-2xl font-bold mb-6">Your Library</h1>

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
        @edit="openEdit(song)"
      />
    </div>

    <AddToPlaylistModal
      v-if="showAddToPlaylistModal"
      :song-id="selectedSongId"
      @close="showAddToPlaylistModal = false"
    />

    <EditSongModal
      v-if="showEditModal"
      :song="selectedSong"
      @close="showEditModal = false"
      @updated="loadSongs"
    />
  </div>
</template>