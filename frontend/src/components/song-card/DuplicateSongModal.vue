<script setup>
const props = defineProps({
  existingSong: { type: Object, required: true },
  canMerge: { type: Boolean, required: true },
})
const emit = defineEmits(['merge', 'new', 'close'])
</script>

<template>
  <div class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4">
    <div class="bg-neutral-900 rounded-lg p-6 w-full max-w-sm">
      <h3 class="text-white font-semibold mb-2">Possible Duplicate</h3>
      <p class="text-neutral-400 text-sm mb-4">
        A song with the same title and artist already exists:
      </p>

      <div class="bg-neutral-800 rounded-lg p-3 mb-4">
        <p class="text-white text-sm font-medium">{{ existingSong.title }}</p>
        <p class="text-neutral-400 text-xs">{{ existingSong.artist }}</p>
        <span class="inline-block mt-1 text-xs text-neutral-500 uppercase">
          Currently has: {{ existingSong.flac_path ? 'FLAC' : '' }}{{ existingSong.flac_path && existingSong.mp3_path ? ' + ' : '' }}{{ existingSong.mp3_path ? 'MP3' : '' }}
        </span>
      </div>

      <p v-if="canMerge" class="text-neutral-400 text-sm mb-4">
        You can add this file as an additional quality option for the existing song, or upload it as a separate new song.
      </p>
      <p v-else class="text-neutral-400 text-sm mb-4">
        The existing song already has this format, so merging isn't possible. You can still upload this as a separate new song.
      </p>

      <div class="flex flex-col gap-2">
        <button
          v-if="canMerge"
          @click="emit('merge')"
          class="w-full bg-green-500 hover:bg-green-400 active:scale-[0.98] text-black font-bold rounded-full py-2.5 shadow-lg shadow-green-500/10 transition"
        >
          Add to Existing Song
        </button>
        <button
          @click="emit('new')"
          class="w-full bg-neutral-800 hover:bg-neutral-700 active:scale-[0.98] text-white font-medium rounded-full py-2.5 transition"
        >
          Upload as New Song
        </button>
        <button
          @click="emit('close')"
          class="w-full text-neutral-400 hover:text-white text-sm py-2 transition"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
</template>