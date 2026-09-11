<script setup>
import { Heart, ListPlus, X, Music, Pencil, GripVertical } from 'lucide-vue-next'
import { getCoverUrl } from '@/services/songs'
import { formatDate, formatDuration } from '@/utils/format'

defineProps({
  song: { type: Object, required: true },
  index: { type: Number, default: null },
  showRemove: { type: Boolean, default: false },
  showAddToPlaylist: { type: Boolean, default: true },
  draggable: { type: Boolean, default: false },
})
const emit = defineEmits(['play', 'toggle-like', 'add-to-playlist', 'remove', 'edit'])
</script>

<template>
  <div class="flex items-center gap-3 px-3 py-2.5 hover:bg-neutral-800/60 cursor-pointer group">
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <div v-if="draggable" class="text-neutral-600 group-hover:text-neutral-400 cursor-grab active:cursor-grabbing flex-shrink-0">
        <GripVertical :size="16" />
      </div>
      <span v-if="index !== null" class="text-neutral-400 text-sm w-4 text-right flex-shrink-0 hidden sm:block">
        {{ index + 1 }}
      </span>
      <div @click="emit('play')" class="w-10 h-10 bg-neutral-700 rounded flex-shrink-0 overflow-hidden flex items-center justify-center">
        <img v-if="getCoverUrl(song)" :src="getCoverUrl(song)" class="w-full h-full object-cover" alt="" />
        <Music v-else :size="16" class="text-neutral-500" />
      </div>
      <div @click="emit('play')" class="min-w-0">
        <p class="text-white text-sm font-medium truncate">{{ song.title }}</p>
        <p class="text-neutral-400 text-xs truncate">{{ song.artist }}</p>
      </div>
    </div>

    <div class="hidden sm:block w-40 flex-shrink-0 min-w-0">
      <p class="text-neutral-400 text-sm truncate">{{ song.album || '-' }}</p>
    </div>

    <div class="hidden sm:block w-28 flex-shrink-0">
      <p class="text-neutral-400 text-sm whitespace-nowrap">{{ formatDate(song.added_at) }}</p>
    </div>

    <div class="w-16 flex-shrink-0 text-right">
      <span class="text-neutral-400 text-sm hidden sm:block">
        {{ formatDuration(song.duration_seconds) }}
      </span>
    </div>

    <div class="w-12 flex-shrink-0 flex items-center justify-end">
      <span v-if="song.source_format === 'flac'" class="text-xs text-green-500 border border-green-500 rounded px-1.5 whitespace-nowrap">
        FLAC
      </span>
    </div>

    <div class="w-40 flex-shrink-0 flex items-center justify-end gap-1">
      <button
        v-if="showAddToPlaylist"
        @click="emit('add-to-playlist')"
        :class="[
          'p-1.5 rounded-full transition hover:bg-neutral-700',
          song.is_in_playlist ? 'text-green-500' : 'text-neutral-400 hover:text-white',
        ]"
        title="Add to playlist"
      >
        <ListPlus :size="18" />
      </button>

      <button
        @click="emit('toggle-like')"
        :class="[
          'p-1.5 rounded-full transition hover:bg-neutral-700',
          song.is_liked ? 'text-green-500' : 'text-neutral-400 hover:text-white',
        ]"
        :title="song.is_liked ? 'Unlike' : 'Like'"
      >
        <Heart :size="18" :fill="song.is_liked ? 'currentColor' : 'none'" />
      </button>

      <button
        @click="emit('edit')"
        class="p-1.5 rounded-full transition hover:bg-neutral-700 text-neutral-400 hover:text-white"
        title="Edit song"
      >
        <Pencil :size="16" />
      </button>

      <button
        v-if="showRemove"
        @click="emit('remove')"
        class="p-1.5 rounded-full transition hover:bg-neutral-700 text-neutral-400 hover:text-red-400"
        title="Remove from playlist"
      >
        <X :size="18" />
      </button>
    </div>
  </div>
</template>