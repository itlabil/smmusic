<script setup>
import { Menu } from 'lucide-vue-next'
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()
</script>

<template>
  <div class="md:hidden flex items-center justify-between px-4 py-3 bg-black border-b border-neutral-800 flex-shrink-0">
    <h1 class="text-white text-lg font-bold">SMMusic</h1>
    <button @click="ui.toggleMobileMenu" class="text-white p-1">
      <Menu :size="22" />
    </button>
  </div>
</template>