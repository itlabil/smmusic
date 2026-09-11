<script setup>
import { ref, watch } from 'vue'
import { searchSongs, likeSong, unlikeSong } from '@/services/songs'
import { usePlayerStore } from '@/stores/player'
import { Search } from 'lucide-vue-next'
import SongRow from '@/components/song-card/SongRow.vue'
import SongListHeader from '@/components/song-card/SongListHeader.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import EditSongModal from '@/components/song-card/EditSongModal.vue'

const query = ref('')
const songs = ref([])
const isSearching = ref(false)
const player = usePlayerStore()
const showAddToPlaylistModal = ref(false)
const showEditModal = ref(false)
const selectedSongId = ref(null)
const selectedSong = ref(null)

let debounceTimer = null

watch(query, (newQuery) => {
  clearTimeout(debounceTimer)
  if (!newQuery.trim()) {
    songs.value = []
    return
  }
  debounceTimer = setTimeout(async () => {
    isSearching.value = true
    const res = await searchSongs(newQuery.trim())
    songs.value = res.data.songs || []
    isSearching.value = false
  }, 300)
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
  showAddToPlaylistModal.value = true
}

function handleAdded() {
  const song = songs.value.find((s) => s.id === selectedSongId.value)
  if (song) song.is_in_playlist = true
}

function openEdit(song) {
  selectedSong.value = song
  showEditModal.value = true
}
</script>

<template>
  <div class="p-4 md:p-8">
    <div class="relative max-w-md mb-6">
      <Search :size="18" class="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500" />
      <input
        v-model="query"
        type="text"
        placeholder="Search songs or artists..."
        class="w-full bg-neutral-800 text-white rounded-full pl-10 pr-4 py-2.5 outline-none focus:ring-2 focus:ring-green-500"
      />
    </div>

    <div v-if="isSearching" class="text-neutral-400">Searching...</div>

    <div v-else-if="query && songs.length === 0" class="text-neutral-400">
      No results for "{{ query }}"
    </div>

    <div v-else-if="songs.length > 0">
      <SongListHeader :show-index="false" />
      <div class="divide-y divide-neutral-800/60">
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
      @updated="() => {}"
    />
  </div>
</template>