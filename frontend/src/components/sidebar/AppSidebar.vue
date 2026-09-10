<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { listPlaylists, createPlaylist } from '@/services/playlists'
import { useAuthStore } from '@/stores/auth'

const playlists = ref([])
const router = useRouter()
const authStore = useAuthStore()

async function loadPlaylists() {
  const res = await listPlaylists()
  playlists.value = res.data.playlists || []
}

onMounted(loadPlaylists)

async function handleCreatePlaylist() {
  const name = prompt('Playlist name:')
  if (!name || !name.trim()) return

  const res = await createPlaylist(name.trim())
  playlists.value.unshift(res.data.playlist)
  router.push(`/playlists/${res.data.playlist.id}`)
}

defineExpose({ loadPlaylists })
</script>

<template>
  <aside class="w-60 bg-black h-full flex flex-col text-neutral-300 flex-shrink-0">
    <div class="p-6">
      <h1 class="text-white text-xl font-bold">SMMusic</h1>
    </div>

    <nav class="px-3 space-y-1">
      <RouterLink
        to="/"
        class="flex items-center gap-3 px-3 py-2 rounded hover:bg-neutral-800 hover:text-white transition"
        active-class="text-white bg-neutral-800"
      >
        Home
      </RouterLink>
      <RouterLink
        to="/liked"
        class="flex items-center gap-3 px-3 py-2 rounded hover:bg-neutral-800 hover:text-white transition"
        active-class="text-white bg-neutral-800"
      >
        Liked Songs
      </RouterLink>
      <RouterLink
        to="/upload"
        class="flex items-center gap-3 px-3 py-2 rounded hover:bg-neutral-800 hover:text-white transition"
        active-class="text-white bg-neutral-800"
      >
        Upload
      </RouterLink>
      <RouterLink
        v-if="authStore.isAdmin"
        to="/admin/users"
        class="flex items-center gap-3 px-3 py-2 rounded hover:bg-neutral-800 hover:text-white transition"
        active-class="text-white bg-neutral-800"
      >
        Users
      </RouterLink>
    </nav>

    <div class="mt-6 px-6 flex items-center justify-between">
      <h2 class="text-xs font-semibold uppercase text-neutral-500">Playlists</h2>
      <button @click="handleCreatePlaylist" class="text-neutral-400 hover:text-white text-lg leading-none">
        +
      </button>
    </div>

    <div class="mt-2 px-6 flex-1 overflow-y-auto">
      <ul class="space-y-1">
        <li v-for="playlist in playlists" :key="playlist.id">
          <RouterLink
            :to="`/playlists/${playlist.id}`"
            class="block py-1.5 text-sm hover:text-white transition truncate"
            active-class="text-white"
          >
            {{ playlist.name }}
          </RouterLink>
        </li>
      </ul>
    </div>
  </aside>
</template>