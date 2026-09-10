<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

async function handleLogin() {
  errorMessage.value = ''
  isLoading.value = true

  try {
    await authStore.login(username.value, password.value)
    router.push('/')
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Login failed. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-black flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <h1 class="text-white text-3xl font-bold text-center mb-8">SMMusic</h1>

      <form @submit.prevent="handleLogin" class="bg-neutral-900 rounded-lg p-8 space-y-5">
        <div>
          <label class="block text-neutral-300 text-sm font-medium mb-1">Username</label>
          <input
            v-model="username"
            type="text"
            required
            class="w-full bg-neutral-800 text-white rounded px-4 py-2.5 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>

        <div>
          <label class="block text-neutral-300 text-sm font-medium mb-1">Password</label>
          <input
            v-model="password"
            type="password"
            required
            class="w-full bg-neutral-800 text-white rounded px-4 py-2.5 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>

        <p v-if="errorMessage" class="text-red-500 text-sm">{{ errorMessage }}</p>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full bg-green-500 hover:bg-green-400 disabled:opacity-50 text-black font-bold rounded-full py-3 transition"
        >
          {{ isLoading ? 'Logging in...' : 'Log In' }}
        </button>
      </form>
    </div>
  </div>
</template>