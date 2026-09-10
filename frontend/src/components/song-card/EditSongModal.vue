<script setup>
import { ref } from 'vue'
import { updateSong, uploadCover, getCoverUrl } from '@/services/songs'

const props = defineProps({
  song: { type: Object, required: true },
})
const emit = defineEmits(['close', 'updated'])

const title = ref(props.song.title)
const artist = ref(props.song.artist)
const album = ref(props.song.album || '')
const genre = ref(props.song.genre || '')
const selectedCoverFile = ref(null)
const coverPreview = ref(getCoverUrl(props.song))
const isSaving = ref(false)
const errorMessage = ref('')

function handleCoverChange(e) {
  const file = e.target.files[0]
  if (!file) return
  selectedCoverFile.value = file
  coverPreview.value = URL.createObjectURL(file)
}

async function handleSave() {
  errorMessage.value = ''
  isSaving.value = true

  try {
    await updateSong(props.song.id, {
      title: title.value,
      artist: artist.value,
      album: album.value || null,
      genre: genre.value || null,
    })

    if (selectedCoverFile.value) {
      await uploadCover(props.song.id, selectedCoverFile.value)
    }

    emit('updated')
    emit('close')
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Failed to update song'
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div
    class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4"
    @click.self="emit('close')"
  >
    <div class="bg-neutral-900 rounded-lg p-6 w-full max-w-md">
      <h3 class="text-white font-semibold mb-4">Edit Song</h3>

      <div class="flex justify-center mb-4">
        <label class="cursor-pointer group relative">
          <div class="w-24 h-24 bg-neutral-800 rounded overflow-hidden flex items-center justify-center">
            <img v-if="coverPreview" :src="coverPreview" class="w-full h-full object-cover" alt="" />
            <span v-else class="text-neutral-500 text-xs text-center px-2">No cover</span>
          </div>
          <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition flex items-center justify-center rounded">
            <span class="text-white text-xs">Change</span>
          </div>
          <input type="file" accept=".jpg,.jpeg,.png" class="hidden" @change="handleCoverChange" />
        </label>
      </div>

      <div class="space-y-3">
        <div>
          <label class="block text-neutral-300 text-sm mb-1">Title</label>
          <input
            v-model="title"
            type="text"
            class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
        <div>
          <label class="block text-neutral-300 text-sm mb-1">Artist</label>
          <input
            v-model="artist"
            type="text"
            class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
        <div>
          <label class="block text-neutral-300 text-sm mb-1">Album</label>
          <input
            v-model="album"
            type="text"
            class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
        <div>
          <label class="block text-neutral-300 text-sm mb-1">Genre</label>
          <input
            v-model="genre"
            type="text"
            class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
      </div>

      <p v-if="errorMessage" class="text-red-500 text-sm mt-3">{{ errorMessage }}</p>

      <div class="flex gap-3 mt-6">
        <button
          @click="handleSave"
          :disabled="isSaving"
          class="flex-1 bg-green-500 hover:bg-green-400 disabled:opacity-50 text-black font-bold rounded-full py-2 transition"
        >
          {{ isSaving ? 'Saving...' : 'Save' }}
        </button>
        <button
          @click="emit('close')"
          class="flex-1 bg-neutral-800 hover:bg-neutral-700 text-white rounded-full py-2 transition"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
</template>