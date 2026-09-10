<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Pencil, Trash2 } from 'lucide-vue-next'
import {
  listPlaylistSongs,
  removeSongFromPlaylist,
  renamePlaylist,
  deletePlaylist,
  reorderPlaylist,
} from '@/services/playlists'
import { likeSong, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import SongRow from '@/components/song-card/SongRow.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'

const route = useRoute()
const router = useRouter()
const player = usePlayerStore()

const playlistId = ref(route.params.id)
const songs = ref([])
const playlistName = ref('')
const showEditModal = ref(false)
const selectedSong = ref(null)
const draggedIndex = ref(null)

async function loadSongs() {
  const res = await listPlaylistSongs(playlistId.value)
  songs.value = res.data.songs || []
}

onMounted(loadSongs)

watch(
  () => route.params.id,
  (newId) => {
    playlistId.value = newId
    loadSongs()
  },
)

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

async function handleRemove(songId) {
  await removeSongFromPlaylist(playlistId.value, songId)
  songs.value = songs.value.filter((s) => s.id !== songId)
}

function openEdit(song) {
  selectedSong.value = song
  showEditModal.value = true
}

function handleDragStart(index) {
  draggedIndex.value = index
}

function handleDragOver(index) {
  if (draggedIndex.value === null || draggedIndex.value === index) return
  const dragged = songs.value[draggedIndex.value]
  songs.value.splice(draggedIndex.value, 1)
  songs.value.splice(index, 0, dragged)
  draggedIndex.value = index
}

async function handleDragEnd() {
  draggedIndex.value = null
  const songIds = songs.value.map((s) => s.id)
  await reorderPlaylist(playlistId.value, songIds)
}

async function handleRename() {
  const name = prompt('New playlist name:', playlistName.value)
  if (!name || !name.trim()) return
  await renamePlaylist(playlistId.value, name.trim())
  playlistName.value = name.trim()
}

async function handleDelete() {
  if (!confirm('Delete this playlist? This cannot be undone.')) return
  await deletePlaylist(playlistId.value)
  router.push('/')
}
</script>

<template>
  <div class="p-4 md:p-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-white text-xl md:text-2xl font-bold">Playlist</h1>
      <div class="flex gap-2">
        <button
          @click="handleRename"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm text-neutral-300 bg-neutral-800 hover:bg-neutral-700 hover:text-white transition"
        >
          <Pencil :size="14" />
          Rename
        </button>
        <button
          @click="handleDelete"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm text-red-400 bg-neutral-800 hover:bg-red-500/10 hover:text-red-300 transition"
        >
          <Trash2 :size="14" />
          Delete
        </button>
      </div>
    </div>

    <div v-if="songs.length === 0" class="text-neutral-400">
      This playlist is empty. Add songs from your library.
    </div>

    <div v-else class="divide-y divide-neutral-800/60">
      <div
        v-for="(song, index) in songs"
        :key="song.id"
        draggable="true"
        @dragstart="handleDragStart(index)"
        @dragover.prevent="handleDragOver(index)"
        @dragend="handleDragEnd"
      >
        <SongRow
          :song="song"
          show-remove
          :show-add-to-playlist="false"
          draggable
          @play="playSong(index)"
          @toggle-like="handleToggleLike(song)"
          @remove="handleRemove(song.id)"
          @edit="openEdit(song)"
        />
      </div>
    </div>

    <EditSongModal
      v-if="showEditModal"
      :song="selectedSong"
      @close="showEditModal = false"
      @updated="loadSongs"
    />
  </div>
</template>