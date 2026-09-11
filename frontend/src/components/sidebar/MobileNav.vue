<script setup>
import { RouterLink } from 'vue-router'
import { Home, Search, Heart, Library } from 'lucide-vue-next'
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()
</script>

<template>
  <nav class="md:hidden relative z-30 flex items-center justify-around bg-black border-t border-neutral-800 py-2 flex-shrink-0">
    <RouterLink
      to="/"
      class="flex flex-col items-center gap-1 px-3 py-1 text-neutral-400"
      active-class="text-white"
    >
      <Home :size="20" />
      <span class="text-xs">Home</span>
    </RouterLink>
    <RouterLink
      to="/search"
      class="flex flex-col items-center gap-1 px-3 py-1 text-neutral-400"
      active-class="text-white"
    >
      <Search :size="20" />
      <span class="text-xs">Search</span>
    </RouterLink>
    <RouterLink
      to="/liked"
      class="flex flex-col items-center gap-1 px-3 py-1 text-neutral-400"
      active-class="text-white"
    >
      <Heart :size="20" />
      <span class="text-xs">Liked</span>
    </RouterLink>
    <button @click="ui.toggleMobileMenu" class="flex flex-col items-center gap-1 px-3 py-1 text-neutral-400">
      <Library :size="20" />
      <span class="text-xs">Playlists</span>
    </button>
  </nav>
</template>