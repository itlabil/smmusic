<script setup>
import { ref, computed } from 'vue'
import { Heart, ListPlus, X, Music, Pencil, GripVertical, MoreVertical, Trash2 } from 'lucide-vue-next'
import { getCoverUrl } from '@/services/songs'
import { formatDate, formatDuration } from '@/utils/format'
import { usePlayerStore } from '@/stores/player'
import PlayingIndicator from '@/components/song-card/PlayingIndicator.vue'

const props = defineProps({
  song: { type: Object, required: true },
  index: { type: Number, default: null },
  showRemove: { type: Boolean, default: false },
  showAddToPlaylist: { type: Boolean, default: true },
  draggable: { type: Boolean, default: false },
  canDelete: { type: Boolean, default: false },
})
const emit = defineEmits(['play', 'toggle-like', 'add-to-playlist', 'remove', 'edit', 'delete'])

const player = usePlayerStore()
const isCurrentSong = computed(() => player.currentSong?.id === props.song.id)

const showMobileMenu = ref(false)

function handleMobileAction(action) {
  emit(action)
  showMobileMenu.value = false
}
</script>

<template>
  <div class="flex items-center gap-3 px-3 py-2.5 hover:bg-neutral-800/60 cursor-pointer group">
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <div v-if="draggable" class="text-neutral-600 group-hover:text-neutral-400 cursor-grab active:cursor-grabbing flex-shrink-0">
        <GripVertical :size="16" />
      </div>
      <span v-if="index !== null && !isCurrentSong" class="hidden lg:flex text-neutral-400 text-sm w-4 justify-end flex-shrink-0">
        {{ index + 1 }}
      </span>
      <span v-else-if="isCurrentSong" class="hidden lg:flex w-4 justify-end flex-shrink-0">
        <PlayingIndicator />
      </span>
      <div @click="emit('play')" class="w-10 h-10 bg-neutral-700 rounded flex-shrink-0 overflow-hidden flex items-center justify-center relative">
        <img v-if="getCoverUrl(song)" :src="getCoverUrl(song)" class="w-full h-full object-cover" alt="" />
        <Music v-else :size="16" class="text-neutral-500" />
        <div v-if="isCurrentSong" class="lg:hidden absolute inset-0 bg-black/50 flex items-center justify-center">
          <PlayingIndicator />
        </div>
      </div>
      <div @click="emit('play')" class="min-w-0">
        <p :class="['text-sm font-medium truncate', isCurrentSong ? 'text-green-500' : 'text-white']">
          {{ song.title }}
        </p>
        <p class="text-neutral-400 text-xs truncate">{{ song.artist }}</p>
      </div>
    </div>

    <div class="hidden lg:block w-40 flex-shrink-0 min-w-0">
      <p class="text-neutral-400 text-sm truncate">{{ song.album || '-' }}</p>
    </div>

    <div class="hidden lg:block w-28 flex-shrink-0">
      <p class="text-neutral-400 text-sm whitespace-nowrap">{{ formatDate(song.added_at) }}</p>
    </div>

    <div class="lg:w-12 flex-shrink-0 flex items-center justify-end">
      <span v-if="song.source_format === 'flac'" class="text-xs text-green-500 border border-green-500 rounded px-1.5 whitespace-nowrap">
        FLAC
      </span>
    </div>

    <div class="lg:w-12 flex-shrink-0 text-right">
      <span class="text-neutral-400 text-sm whitespace-nowrap">
        {{ formatDuration(song.duration_seconds) }}
      </span>
    </div>

    <div class="w-auto lg:w-40 flex-shrink-0 flex items-center justify-end gap-1 relative">
      <div v-if="showAddToPlaylist" class="hidden lg:block">
        <button
          @click="emit('add-to-playlist')"
          :class="[
            'p-1.5 rounded-full transition hover:bg-neutral-700',
            song.is_in_playlist ? 'text-green-500' : 'text-neutral-400 hover:text-white',
          ]"
          title="Add to playlist"
        >
          <ListPlus :size="18" />
        </button>
      </div>

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

      <div class="hidden lg:block">
        <button
          @click="emit('edit')"
          class="p-1.5 rounded-full transition hover:bg-neutral-700 text-neutral-400 hover:text-white"
          title="Edit song"
        >
          <Pencil :size="16" />
        </button>
      </div>

      <div v-if="canDelete" class="hidden lg:block">
        <button
          @click="emit('delete')"
          class="p-1.5 rounded-full transition hover:bg-neutral-700 text-neutral-400 hover:text-red-400"
          title="Delete song"
        >
          <Trash2 :size="16" />
        </button>
      </div>

      <div v-if="showRemove" class="hidden lg:block">
        <button
          @click="emit('remove')"
          class="p-1.5 rounded-full transition hover:bg-neutral-700 text-neutral-400 hover:text-red-400"
          title="Remove from playlist"
        >
          <X :size="18" />
        </button>
      </div>

      <div class="lg:hidden">
        <button
          @click.stop="showMobileMenu = !showMobileMenu"
          class="p-1.5 rounded-full transition hover:bg-neutral-700 text-neutral-400 hover:text-white"
        >
          <MoreVertical :size="18" />
        </button>

        <div
          v-if="showMobileMenu"
          class="absolute right-0 top-full mt-1 bg-neutral-800 rounded-lg shadow-xl overflow-hidden z-10 w-44"
          @click.stop
        >
          <button
            v-if="showAddToPlaylist"
            @click="handleMobileAction('add-to-playlist')"
            :class="[
              'w-full flex items-center gap-2 px-4 py-2.5 text-sm text-left hover:bg-neutral-700 transition',
              song.is_in_playlist ? 'text-green-500' : 'text-white',
            ]"
          >
            <ListPlus :size="16" />
            Add to Playlist
          </button>
          <button
            @click="handleMobileAction('edit')"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-left text-white hover:bg-neutral-700 transition"
          >
            <Pencil :size="16" />
            Edit Song
          </button>
          <button
            v-if="canDelete"
            @click="handleMobileAction('delete')"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-left text-red-400 hover:bg-neutral-700 transition"
          >
            <Trash2 :size="16" />
            Delete Song
          </button>
          <button
            v-if="showRemove"
            @click="handleMobileAction('remove')"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-left text-red-400 hover:bg-neutral-700 transition"
          >
            <X :size="16" />
            Remove
          </button>
        </div>
      </div>
    </div>
  </div>
</template>