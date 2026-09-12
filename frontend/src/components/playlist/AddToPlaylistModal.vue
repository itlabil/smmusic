<script setup>
import { onMounted, ref } from 'vue'
import { Check, ListMusic, Plus } from 'lucide-vue-next'
import { listPlaylists, createPlaylist, addSongToPlaylist } from '@/services/playlists'

const props = defineProps({
  songId: { type: Number, required: true },
})
const emit = defineEmits(['close', 'added'])

const playlists = ref([])
const isLoading = ref(true)
const addedToId = ref(null)
const isCreating = ref(false)
const newPlaylistName = ref('')

onMounted(async () => {
  const res = await listPlaylists()
  playlists.value = res.data.playlists || []
  isLoading.value = false
})

async function handleAdd(playlistId) {
  await addSongToPlaylist(playlistId, props.songId)
  addedToId.value = playlistId
  emit('added')
  setTimeout(() => emit('close'), 500)
}

async function handleCreateAndAdd() {
  if (!newPlaylistName.value.trim()) return
  isCreating.value = true
  try {
    const res = await createPlaylist(newPlaylistName.value.trim())
    await handleAdd(res.data.playlist.id)
  } finally {
    isCreating.value = false
  }
}
</script>

<template>
  <div
    class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4"
    @click.self="emit('close')"
  >
    <div class="bg-neutral-900 rounded-lg w-full max-w-sm max-h-[28rem] flex flex-col overflow-hidden">
      <div class="p-5 pb-3 flex-shrink-0">
        <h3 class="text-white font-semibold">Add to Playlist</h3>
      </div>

      <div class="flex-1 overflow-y-auto px-3">
        <div v-if="isLoading" class="px-2 py-6 text-neutral-400 text-sm text-center">
          Loading playlists...
        </div>

        <div v-else-if="playlists.length === 0" class="px-2 py-6 text-neutral-400 text-sm text-center">
          You don't have any playlists yet. Create one below.
        </div>

        <ul v-else class="space-y-0.5">
          <li v-for="playlist in playlists" :key="playlist.id">
            <button
              @click="handleAdd(playlist.id)"
              :disabled="addedToId !== null"
              class="w-full flex items-center gap-3 text-left px-3 py-2.5 rounded-lg hover:bg-neutral-800 transition disabled:cursor-default"
            >
              <div class="w-9 h-9 bg-neutral-800 rounded flex items-center justify-center flex-shrink-0">
                <ListMusic :size="16" class="text-neutral-500" />
              </div>
              <span class="flex-1 min-w-0 text-white text-sm truncate">{{ playlist.name }}</span>
              <Check v-if="addedToId === playlist.id" :size="18" class="text-green-500 flex-shrink-0" />
            </button>
          </li>
        </ul>
      </div>

      <div class="p-3 border-t border-neutral-800 flex-shrink-0">
        <div class="flex gap-2">
          <input
            v-model="newPlaylistName"
            type="text"
            placeholder="New playlist name"
            @keyup.enter="handleCreateAndAdd"
            class="flex-1 bg-neutral-800 text-white text-sm rounded-full px-4 py-2 outline-none focus:ring-2 focus:ring-green-500"
          />
          <button
            @click="handleCreateAndAdd"
            :disabled="!newPlaylistName.trim() || isCreating"
            class="bg-green-500 hover:bg-green-400 active:scale-[0.98] disabled:opacity-40 disabled:cursor-not-allowed text-black rounded-full p-2 transition flex-shrink-0"
            title="Create and add"
          >
            <Plus :size="18" />
          </button>
        </div>
      </div>

      <div class="p-3 pt-0">
        <button
          @click="emit('close')"
          class="w-full bg-neutral-800 hover:bg-neutral-700 active:scale-[0.98] text-white font-medium rounded-full py-2 text-sm transition"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
</template>