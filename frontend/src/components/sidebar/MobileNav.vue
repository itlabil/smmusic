<script setup>
import { RouterLink } from 'vue-router'
import { Home, Heart, Upload } from 'lucide-vue-next'
</script>

<template>
  <nav class="md:hidden flex items-center justify-around bg-black border-t border-neutral-800 py-2 flex-shrink-0">
    <RouterLink
      to="/"
      class="flex flex-col items-center gap-1 px-4 py-1 text-neutral-400"
      active-class="text-white"
    >
      <Home :size="20" />
      <span class="text-xs">Home</span>
    </RouterLink>
    <RouterLink
      to="/liked"
      class="flex flex-col items-center gap-1 px-4 py-1 text-neutral-400"
      active-class="text-white"
    >
      <Heart :size="20" />
      <span class="text-xs">Liked</span>
    </RouterLink>
    <RouterLink
      to="/upload"
      class="flex flex-col items-center gap-1 px-4 py-1 text-neutral-400"
      active-class="text-white"
    >
      <Upload :size="20" />
      <span class="text-xs">Upload</span>
    </RouterLink>
  </nav>
</template>