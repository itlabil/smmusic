<script setup>
import { onMounted, ref } from 'vue'
import { listPlaylists, addSongToPlaylist } from '@/services/playlists'

const props = defineProps({
  songId: { type: Number, required: true },
})
const emit = defineEmits(['close', 'added'])

const playlists = ref([])

onMounted(async () => {
  const res = await listPlaylists()
  playlists.value = res.data.playlists || []
})

async function handleAdd(playlistId) {
  await addSongToPlaylist(playlistId, props.songId)
  emit('added')
  emit('close')
}
</script>

<template>
  <div
    class="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
    @click.self="emit('close')"
  >
    <div class="bg-neutral-900 rounded-lg p-6 w-80 max-h-96 overflow-y-auto">
      <h3 class="text-white font-semibold mb-4">Add to Playlist</h3>

      <div v-if="playlists.length === 0" class="text-neutral-400 text-sm">
        No playlists yet. Create one from the sidebar first.
      </div>

      <ul v-else class="space-y-1">
        <li v-for="playlist in playlists" :key="playlist.id">
          <button
            @click="handleAdd(playlist.id)"
            class="w-full text-left px-3 py-2 rounded hover:bg-neutral-800 text-white text-sm"
          >
            {{ playlist.name }}
          </button>
        </li>
      </ul>

      <button @click="emit('close')" class="mt-4 text-neutral-400 hover:text-white text-sm underline">
        Cancel
      </button>
    </div>
  </div>
</template>