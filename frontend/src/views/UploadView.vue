<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { uploadSong } from '@/services/songs'

const router = useRouter()
const selectedFile = ref(null)
const isUploading = ref(false)
const uploadProgress = ref(0)
const errorMessage = ref('')
const successMessage = ref('')

function handleFileChange(e) {
  const file = e.target.files[0]
  errorMessage.value = ''
  successMessage.value = ''

  if (!file) return

  const ext = file.name.toLowerCase().split('.').pop()
  if (ext !== 'flac' && ext !== 'mp3') {
    errorMessage.value = 'Only .flac and .mp3 files are allowed'
    selectedFile.value = null
    return
  }

  selectedFile.value = file
}

async function handleUpload() {
  if (!selectedFile.value) return

  isUploading.value = true
  uploadProgress.value = 0
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const res = await uploadSong(selectedFile.value, (progressEvent) => {
      uploadProgress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total)
    })

    successMessage.value = `"${res.data.song.title}" uploaded successfully${
      res.data.song.source_format === 'flac' ? ' — transcoding to MP3 in background' : ''
    }`
    selectedFile.value = null
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Upload failed. Please try again.'
  } finally {
    isUploading.value = false
    uploadProgress.value = 0
  }
}
</script>

<template>
  <div class="p-8 max-w-2xl">
    <h1 class="text-white text-2xl font-bold mb-6">Upload Song</h1>

    <div class="bg-neutral-900 rounded-lg p-8">
      <label
        class="flex flex-col items-center justify-center border-2 border-dashed border-neutral-700 rounded-lg py-12 cursor-pointer hover:border-green-500 transition"
      >
        <input type="file" accept=".flac,.mp3" class="hidden" @change="handleFileChange" />
        <p v-if="!selectedFile" class="text-neutral-400 text-sm">
          Click to select a FLAC or MP3 file
        </p>
        <p v-else class="text-white text-sm font-medium">{{ selectedFile.name }}</p>
      </label>

      <p v-if="errorMessage" class="text-red-500 text-sm mt-4">{{ errorMessage }}</p>
      <p v-if="successMessage" class="text-green-500 text-sm mt-4">{{ successMessage }}</p>

      <div v-if="isUploading" class="mt-4">
        <div class="w-full bg-neutral-800 rounded-full h-2">
          <div
            class="bg-green-500 h-2 rounded-full transition-all"
            :style="{ width: uploadProgress + '%' }"
          ></div>
        </div>
        <p class="text-neutral-400 text-xs mt-1">{{ uploadProgress }}%</p>
      </div>

      <button
        @click="handleUpload"
        :disabled="!selectedFile || isUploading"
        class="w-full mt-6 bg-green-500 hover:bg-green-400 disabled:opacity-50 disabled:cursor-not-allowed text-black font-bold rounded-full py-3 transition"
      >
        {{ isUploading ? 'Uploading...' : 'Upload' }}
      </button>
    </div>

    <button @click="router.push('/')" class="text-neutral-400 hover:text-white text-sm mt-4 underline">
      Back to Library
    </button>
  </div>
</template>