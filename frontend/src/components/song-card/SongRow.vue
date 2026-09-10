<script setup>
defineProps({
  song: { type: Object, required: true },
  showRemove: { type: Boolean, default: false },
  showAddToPlaylist: { type: Boolean, default: true },
})
const emit = defineEmits(['play', 'toggle-like', 'add-to-playlist', 'remove'])
</script>

<template>
  <div
    class="flex items-center gap-3 px-3 py-2 rounded hover:bg-neutral-800 cursor-pointer group"
  >
    <div @click="emit('play')" class="w-10 h-10 bg-neutral-700 rounded flex-shrink-0"></div>
    <div @click="emit('play')" class="min-w-0 flex-1">
      <p class="text-white text-sm font-medium truncate">{{ song.title }}</p>
      <p class="text-neutral-400 text-xs truncate">{{ song.artist }}</p>
    </div>
    <span v-if="song.source_format === 'flac'" class="text-xs text-green-500 border border-green-500 rounded px-1.5">
      FLAC
    </span>

    <button
      v-if="showAddToPlaylist"
      @click="emit('add-to-playlist')"
      class="text-neutral-500 hover:text-white opacity-0 group-hover:opacity-100 transition text-sm px-2"
      title="Add to playlist"
    >
      + Playlist
    </button>

    <button
      @click="emit('toggle-like')"
      :class="[
        'transition text-lg px-2',
        song.is_liked ? 'text-green-500' : 'text-neutral-500 hover:text-white opacity-0 group-hover:opacity-100',
      ]"
      :title="song.is_liked ? 'Unlike' : 'Like'"
    >
      {{ song.is_liked ? '♥' : '♡' }}
    </button>

    <button
      v-if="showRemove"
      @click="emit('remove')"
      class="text-neutral-500 hover:text-red-400 opacity-0 group-hover:opacity-100 transition text-sm px-2"
      title="Remove from playlist"
    >
      ✕
    </button>
  </div>
</template>