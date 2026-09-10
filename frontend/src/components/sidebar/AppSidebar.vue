<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import api from '@/services/api'

const playlists = ref([])

onMounted(async () => {
  const res = await api.get('/playlists')
  playlists.value = res.data.playlists || []
})
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
    </nav>

    <div class="mt-6 px-6 flex-1 overflow-y-auto">
      <h2 class="text-xs font-semibold uppercase text-neutral-500 mb-2">Playlists</h2>
      <ul class="space-y-1">
        <li v-for="playlist in playlists" :key="playlist.id">
          <RouterLink
            :to="`/playlists/${playlist.id}`"
            class="block py-1.5 text-sm hover:text-white transition truncate"
          >
            {{ playlist.name }}
          </RouterLink>
        </li>
      </ul>
    </div>
  </aside>
</template>