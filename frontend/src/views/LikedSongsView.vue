<script setup>
import { onMounted, ref } from 'vue'
import { Heart } from 'lucide-vue-next'
import { listLikedSongs, unlikeSong, getLikedSongsSummary } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import { formatTotalDuration } from '@/utils/format'
import SongRow from '@/components/song-card/SongRow.vue'
import SongListHeader from '@/components/song-card/SongListHeader.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import { useAuthStore } from '@/stores/auth'

const songs = ref([])
const summary = ref(null)
const player = usePlayerStore()
const authStore = useAuthStore()
const showEditModal = ref(false)
const selectedSong = ref(null)
const showAddToPlaylistModal = ref(false)
const selectedSongId = ref(null)

async function loadSongs() {
  const res = await listLikedSongs()
  songs.value = (res.data.songs || []).map((s) => ({ ...s, is_liked: true }))
}

async function loadSummary() {
  const res = await getLikedSongsSummary()
  summary.value = res.data.summary
}

onMounted(() => {
  loadSongs()
  loadSummary()
})

function playSong(index) {
  player.playQueue(songs.value, index)
}

async function handleToggleLike(song) {
  await unlikeSong(song.id)
  songs.value = songs.value.filter((s) => s.id !== song.id)
  await loadSummary()
}

function openEdit(song) {
  selectedSong.value = song
  showEditModal.value = true
}

function openAddToPlaylist(songId) {
  selectedSongId.value = songId
  showAddToPlaylistModal.value = true
}

function handleAdded() {
  const song = songs.value.find((s) => s.id === selectedSongId.value)
  if (song) song.is_in_playlist = true
}
</script>

<template>
  <div>
    <div class="bg-gradient-to-b from-purple-800 to-neutral-900 p-4 md:p-8 pb-6">
      <div class="flex flex-col sm:flex-row items-start sm:items-end gap-4">
        <div class="w-32 h-32 sm:w-48 sm:h-48 bg-gradient-to-br from-indigo-500 to-purple-300 rounded shadow-xl flex items-center justify-center flex-shrink-0">
          <Heart :size="48" class="text-white" fill="white" />
        </div>
        <div class="min-w-0">
          <p class="text-neutral-300 text-xs font-medium uppercase">Playlist</p>
          <h1 class="text-white text-2xl sm:text-4xl md:text-5xl font-bold mt-1 mb-3">
            Liked Songs
          </h1>
          <div v-if="summary" class="flex items-center gap-1 text-neutral-300 text-sm">
            <span class="font-medium text-white">{{ authStore.user?.username }}</span>
            <span>•</span>
            <span>{{ summary.song_count }} songs</span>
            <span v-if="summary.total_duration_sec > 0">•</span>
            <span v-if="summary.total_duration_sec > 0">{{ formatTotalDuration(summary.total_duration_sec) }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="p-4 md:p-8">
      <div v-if="songs.length === 0" class="text-neutral-400">
        You haven't liked any songs yet.
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
            @edit="openEdit(song)"
            @add-to-playlist="openAddToPlaylist(song.id)"
          />
        </div>
      </div>
    </div>

    <EditSongModal
      v-if="showEditModal"
      :song="selectedSong"
      @close="showEditModal = false"
      @updated="loadSongs"
    />

    <AddToPlaylistModal
      v-if="showAddToPlaylistModal"
      :song-id="selectedSongId"
      @added="handleAdded"
      @close="showAddToPlaylistModal = false"
    />
  </div>
</template>