<script setup>
import { onMounted, ref } from 'vue'
import { listLikedSongs, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import SongRow from '@/components/song-card/SongRow.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'

const songs = ref([])
const player = usePlayerStore()
const showEditModal = ref(false)
const selectedSong = ref(null)

function openEdit(song) {
  selectedSong.value = song
  showEditModal.value = true
}

async function loadSongs() {
  const res = await listLikedSongs()
  songs.value = (res.data.songs || []).map((s) => ({ ...s, is_liked: true }))
}

onMounted(loadSongs)

function playSong(index) {
  player.playQueue(songs.value, index)
}

async function handleToggleLike(song) {
  await unlikeSong(song.id)
  songs.value = songs.value.filter((s) => s.id !== song.id)
}
</script>

<template>
  <div class="p-4 md:p-8">
    <h1 class="text-white text-xl md:text-2xl font-bold mb-6">Liked Songs</h1>

    <div v-if="songs.length === 0" class="text-neutral-400">
      You haven't liked any songs yet.
    </div>

    <div v-else class="space-y-1">
      <SongRow
        v-for="(song, index) in songs"
        :key="song.id"
        :song="song"
        @play="playSong(index)"
        @toggle-like="handleToggleLike(song)"
        @edit="openEdit(song)"
      />
    </div>

    <EditSongModal
      v-if="showEditModal"
      :song="selectedSong"
      @close="showEditModal = false"
      @updated="loadSongs"
    />
  </div>
</template>