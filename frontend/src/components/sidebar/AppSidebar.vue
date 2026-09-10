<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Home, Heart, Upload, Users, Plus, Library } from 'lucide-vue-next'
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
  <aside class="w-60 bg-black h-full flex-col text-neutral-300 flex-shrink-0 hidden md:flex">
    <div class="p-6">
      <h1 class="text-white text-xl font-bold tracking-tight">SMMusic</h1>
    </div>

    <nav class="px-3 space-y-1">
      <RouterLink
        to="/"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-white bg-neutral-800"
      >
        <Home :size="18" />
        Home
      </RouterLink>
      <RouterLink
        to="/liked"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-white bg-neutral-800"
      >
        <Heart :size="18" />
        Liked Songs
      </RouterLink>
      <RouterLink
        to="/upload"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-white bg-neutral-800"
      >
        <Upload :size="18" />
        Upload
      </RouterLink>
      <RouterLink
        v-if="authStore.isAdmin"
        to="/admin/users"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-white bg-neutral-800"
      >
        <Users :size="18" />
        Users
      </RouterLink>
    </nav>

    <div class="mt-6 px-6 flex items-center justify-between">
      <div class="flex items-center gap-2 text-neutral-400">
        <Library :size="16" />
        <h2 class="text-xs font-semibold uppercase tracking-wide">Playlists</h2>
      </div>
      <button
        @click="handleCreatePlaylist"
        class="text-neutral-400 hover:text-white transition p-0.5 hover:bg-neutral-800 rounded"
        title="Create playlist"
      >
        <Plus :size="16" />
      </button>
    </div>

    <div class="mt-2 px-6 flex-1 overflow-y-auto pb-4">
      <ul class="space-y-1">
        <li v-for="playlist in playlists" :key="playlist.id">
          <RouterLink
            :to="`/playlists/${playlist.id}`"
            class="block py-1.5 text-sm text-neutral-400 hover:text-white transition truncate"
            active-class="text-white font-medium"
          >
            {{ playlist.name }}
          </RouterLink>
        </li>
      </ul>
    </div>
  </aside>
</template>