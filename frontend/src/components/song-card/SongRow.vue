<script setup>
defineProps({
  song: { type: Object, required: true },
})
const emit = defineEmits(['play'])
</script>

<template>
  <div
    @click="emit('play')"
    class="flex items-center gap-3 px-3 py-2 rounded hover:bg-neutral-800 cursor-pointer group"
  >
    <div class="w-10 h-10 bg-neutral-700 rounded flex-shrink-0"></div>
    <div class="min-w-0 flex-1">
      <p class="text-white text-sm font-medium truncate">{{ song.title }}</p>
      <p class="text-neutral-400 text-xs truncate">{{ song.artist }}</p>
    </div>
    <span v-if="song.source_format === 'flac'" class="text-xs text-green-500 border border-green-500 rounded px-1.5">
      FLAC
    </span>
  </div>
</template>