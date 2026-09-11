<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Pencil, Trash2, ListMusic } from 'lucide-vue-next'
import {
  listPlaylistSongs,
  getPlaylistDetail,
  removeSongFromPlaylist,
  renamePlaylist,
  deletePlaylist,
  reorderPlaylist,
} from '@/services/playlists'
import { likeSong, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import { formatTotalDuration } from '@/utils/format'
import SongRow from '@/components/song-card/SongRow.vue'
import SongListHeader from '@/components/song-card/SongListHeader.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'

const route = useRoute()
const router = useRouter()
const player = usePlayerStore()

const playlistId = ref(route.params.id)
const songs = ref([])
const detail = ref(null)
const showEditModal = ref(false)
const selectedSong = ref(null)
const draggedIndex = ref(null)

async function loadSongs() {
  const res = await listPlaylistSongs(playlistId.value)
  songs.value = res.data.songs || []
}

async function loadDetail() {
  const res = await getPlaylistDetail(playlistId.value)
  detail.value = res.data.playlist
}

onMounted(() => {
  loadSongs()
  loadDetail()
})

watch(
  () => route.params.id,
  (newId) => {
    playlistId.value = newId
    loadSongs()
    loadDetail()
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
  await loadDetail()
}

async function handleRename() {
  const name = prompt('New playlist name:', detail.value?.name)
  if (!name || !name.trim()) return
  await renamePlaylist(playlistId.value, name.trim())
  await loadDetail()
}

async function handleDelete() {
  if (!confirm('Delete this playlist? This cannot be undone.')) return
  await deletePlaylist(playlistId.value)
  router.push('/')
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
</script>

<template>
  <div>
    <div class="bg-gradient-to-b from-neutral-700 to-neutral-900 p-4 md:p-8 pb-6">
      <div class="flex flex-col sm:flex-row items-start sm:items-end gap-4">
        <div class="w-32 h-32 sm:w-48 sm:h-48 bg-neutral-800 rounded shadow-xl flex items-center justify-center flex-shrink-0">
          <ListMusic :size="48" class="text-neutral-500" />
        </div>
        <div class="min-w-0">
          <p class="text-neutral-300 text-xs font-medium uppercase">Playlist</p>
          <h1 class="text-white text-2xl sm:text-4xl md:text-5xl font-bold mt-1 mb-3 break-words">
            {{ detail?.name }}
          </h1>
          <div v-if="detail" class="flex items-center gap-1 text-neutral-300 text-sm">
            <span class="font-medium text-white">{{ detail.owner_username }}</span>
            <span>•</span>
            <span>{{ detail.song_count }} songs</span>
            <span v-if="detail.total_duration_sec > 0">•</span>
            <span v-if="detail.total_duration_sec > 0">{{ formatTotalDuration(detail.total_duration_sec) }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="p-4 md:p-8">
      <div class="flex items-center justify-end gap-2 mb-6">
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

      <div v-if="songs.length === 0" class="text-neutral-400">
        This playlist is empty. Add songs from your library.
      </div>

      <div v-else>
        <SongListHeader />
        <div class="divide-y divide-neutral-800/60">
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
              :index="index"
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