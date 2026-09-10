<script setup>
import { onMounted, ref } from 'vue'
import { listUsers, createUser } from '@/services/admin'

const users = ref([])
const newUsername = ref('')
const newPassword = ref('')
const newRole = ref('user')
const errorMessage = ref('')
const successMessage = ref('')
const isCreating = ref(false)

async function loadUsers() {
  const res = await listUsers()
  users.value = res.data.users || []
}

onMounted(loadUsers)

async function handleCreate() {
  errorMessage.value = ''
  successMessage.value = ''
  isCreating.value = true

  try {
    await createUser(newUsername.value, newPassword.value, newRole.value)
    successMessage.value = `User "${newUsername.value}" created successfully`
    newUsername.value = ''
    newPassword.value = ''
    newRole.value = 'user'
    await loadUsers()
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Failed to create user'
  } finally {
    isCreating.value = false
  }
}
</script>

<template>
  <div class="p-8 max-w-3xl">
    <h1 class="text-white text-2xl font-bold mb-6">Admin — User Management</h1>

    <div class="bg-neutral-900 rounded-lg p-6 mb-8">
      <h2 class="text-white font-semibold mb-4">Create New User</h2>

      <form @submit.prevent="handleCreate" class="space-y-4">
        <div class="flex gap-4">
          <div class="flex-1">
            <label class="block text-neutral-300 text-sm mb-1">Username</label>
            <input
              v-model="newUsername"
              type="text"
              required
              class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
          <div class="flex-1">
            <label class="block text-neutral-300 text-sm mb-1">Password</label>
            <input
              v-model="newPassword"
              type="text"
              required
              minlength="6"
              class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
          <div class="w-32">
            <label class="block text-neutral-300 text-sm mb-1">Role</label>
            <select
              v-model="newRole"
              class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
            >
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>
        </div>

        <p v-if="errorMessage" class="text-red-500 text-sm">{{ errorMessage }}</p>
        <p v-if="successMessage" class="text-green-500 text-sm">{{ successMessage }}</p>

        <button
          type="submit"
          :disabled="isCreating"
          class="bg-green-500 hover:bg-green-400 active:scale-[0.98] disabled:opacity-50 disabled:hover:bg-green-500 text-black font-bold rounded-full px-6 py-2 shadow-lg shadow-green-500/10 transition"
        >
          {{ isCreating ? 'Creating...' : 'Create User' }}
        </button>
      </form>
    </div>

    <div class="bg-neutral-900 rounded-lg p-6">
      <h2 class="text-white font-semibold mb-4">All Users</h2>

      <table class="w-full text-left">
        <thead>
          <tr class="text-neutral-500 text-xs uppercase border-b border-neutral-800">
            <th class="pb-2">Username</th>
            <th class="pb-2">Role</th>
            <th class="pb-2">Created</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="border-b border-neutral-800/50">
            <td class="py-2 text-white text-sm">{{ user.username }}</td>
            <td class="py-2 text-sm">
              <span
                :class="[
                  'px-2 py-0.5 rounded text-xs',
                  user.role === 'admin' ? 'bg-green-900 text-green-400' : 'bg-neutral-800 text-neutral-400',
                ]"
              >
                {{ user.role }}
              </span>
            </td>
            <td class="py-2 text-neutral-400 text-sm">
              {{ new Date(user.created_at).toLocaleDateString() }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>