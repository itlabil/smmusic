<script setup>
import { onMounted, ref } from 'vue'
import { listSongs, listRecentlyPlayed, likeSong, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import SongRow from '@/components/song-card/SongRow.vue'
import SongListHeader from '@/components/song-card/SongListHeader.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'

const songs = ref([])
const recentSongs = ref([])
const player = usePlayerStore()
const showAddToPlaylistModal = ref(false)
const showEditModal = ref(false)
const selectedSongId = ref(null)
const selectedSong = ref(null)

async function loadSongs() {
  const res = await listSongs()
  songs.value = res.data.songs || []
}

async function loadRecentSongs() {
  const res = await listRecentlyPlayed()
  recentSongs.value = res.data.songs || []
}

onMounted(() => {
  loadSongs()
  loadRecentSongs()
})

function playSong(index) {
  player.playQueue(songs.value, index)
}

function playRecentSong(index) {
  player.playQueue(recentSongs.value, index)
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

function handleAdded() {
  const song =
    songs.value.find((s) => s.id === selectedSongId.value) ||
    recentSongs.value.find((s) => s.id === selectedSongId.value)
  if (song) song.is_in_playlist = true
}

function openEdit(song) {
  selectedSong.value = song
  showEditModal.value = true
}

async function handleUpdated() {
  await loadSongs()
  await loadRecentSongs()
}
</script>

<template>
  <div class="p-4 md:p-8">
    <div v-if="recentSongs.length > 0" class="mb-10">
      <h2 class="text-white text-lg md:text-xl font-bold mb-4">Recently Played</h2>
      <SongListHeader :show-index="false" />
      <div class="divide-y divide-neutral-800/60">
        <SongRow
          v-for="(song, index) in recentSongs"
          :key="'recent-' + song.id"
          :song="song"
          @play="playRecentSong(index)"
          @toggle-like="handleToggleLike(song)"
          @add-to-playlist="openAddToPlaylist(song.id)"
          @edit="openEdit(song)"
        />
      </div>
    </div>

    <h1 class="text-white text-xl md:text-2xl font-bold mb-6">Your Library</h1>

    <div v-if="songs.length === 0" class="text-neutral-400">
      No songs yet. Upload your first song to get started.
    </div>

    <div v-else>
      <SongListHeader />
      <div class="divide-y divide-neutral-800/60">
        <SongRow
          v-for="(song, index) in songs"
          :key="song.id"
          :song="song"
          :index="index"
          @play="playSong(index)"
          @toggle-like="handleToggleLike(song)"
          @add-to-playlist="openAddToPlaylist(song.id)"
          @edit="openEdit(song)"
        />
      </div>
    </div>

    <AddToPlaylistModal
      v-if="showAddToPlaylistModal"
      :song-id="selectedSongId"
      @added="handleAdded"
      @close="showAddToPlaylistModal = false"
    />

    <EditSongModal
      v-if="showEditModal"
      :song="selectedSong"
      @close="showEditModal = false"
      @updated="handleUpdated"
    />
  </div>
</template>