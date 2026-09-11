<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Home, Search, Heart, Upload, Users, Plus, Library, X } from 'lucide-vue-next'
import { listPlaylists, createPlaylist } from '@/services/playlists'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const playlists = ref([])
const router = useRouter()
const authStore = useAuthStore()
const ui = useUiStore()

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
  ui.closeMobileMenu()
  router.push(`/playlists/${res.data.playlist.id}`)
}

function handleNavigate() {
  ui.closeMobileMenu()
}

defineExpose({ loadPlaylists })
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <div
    v-if="ui.isMobileMenuOpen"
    class="fixed inset-0 bg-black/60 z-40 md:hidden"
    @click="ui.closeMobileMenu"
  ></div>

  <aside
    :class="[
      'w-64 bg-black h-full flex-col text-neutral-300 flex-shrink-0 z-50 transition-transform duration-200',
      'fixed inset-y-0 left-0 md:static md:flex md:translate-x-0',
      ui.isMobileMenuOpen ? 'flex translate-x-0' : 'hidden -translate-x-full',
    ]"
  >
    <div class="p-6 flex items-center justify-between">
      <h1 class="text-white text-xl font-bold tracking-tight">SMMusic</h1>
      <button @click="ui.closeMobileMenu" class="text-neutral-400 hover:text-white md:hidden">
        <X :size="20" />
      </button>
    </div>

    <nav class="px-3 space-y-1">
      <RouterLink
        to="/"
        @click="handleNavigate"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-green-500 bg-neutral-800"
      >
        <Home :size="18" />
        Home
      </RouterLink>
      <RouterLink
        to="/search"
        @click="handleNavigate"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-green-500 bg-neutral-800"
      >
        <Search :size="18" />
        Search
      </RouterLink>
      <RouterLink
        to="/liked"
        @click="handleNavigate"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-green-500 bg-neutral-800"
      >
        <Heart :size="18" />
        Liked Songs
      </RouterLink>
      <RouterLink
        to="/upload"
        @click="handleNavigate"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-green-500 bg-neutral-800"
      >
        <Upload :size="18" />
        Upload
      </RouterLink>
      <RouterLink
        v-if="authStore.isAdmin"
        to="/admin/users"
        @click="handleNavigate"
        class="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 hover:text-white transition text-sm font-medium"
        active-class="text-green-500 bg-neutral-800"
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
            @click="handleNavigate"
            class="block py-1.5 px-2 -mx-2 rounded text-sm text-neutral-400 hover:text-white hover:bg-neutral-800 transition truncate"
            active-class="text-green-500 font-medium bg-neutral-800"
          >
            {{ playlist.name }}
          </RouterLink>
        </li>
      </ul>
    </div>
  </aside>
</template>