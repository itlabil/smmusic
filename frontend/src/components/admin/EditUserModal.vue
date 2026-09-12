<script setup>
import { ref } from 'vue'
import { updateUser } from '@/services/admin'

const props = defineProps({
  user: { type: Object, required: true },
})
const emit = defineEmits(['close', 'updated'])

const role = ref(props.user.role)
const newPassword = ref('')
const isSaving = ref(false)
const errorMessage = ref('')

async function handleSave() {
  errorMessage.value = ''
  isSaving.value = true

  try {
    await updateUser(props.user.id, {
      role: role.value !== props.user.role ? role.value : undefined,
      newPassword: newPassword.value || undefined,
    })
    emit('updated')
    emit('close')
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Failed to update user'
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
    <div class="bg-neutral-900 rounded-lg p-6 w-full max-w-sm">
      <h3 class="text-white font-semibold mb-1">Edit User</h3>
      <p class="text-neutral-400 text-sm mb-4">{{ user.username }}</p>

      <div class="space-y-3">
        <div>
          <label class="block text-neutral-300 text-sm mb-1">Role</label>
          <select
            v-model="role"
            class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
          >
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div>
          <label class="block text-neutral-300 text-sm mb-1">New Password (optional)</label>
          <input
            v-model="newPassword"
            type="text"
            minlength="6"
            placeholder="Leave blank to keep current password"
            class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
      </div>

      <p v-if="errorMessage" class="text-red-500 text-sm mt-3">{{ errorMessage }}</p>

      <div class="flex gap-3 mt-6">
        <button
          @click="handleSave"
          :disabled="isSaving"
          class="flex-1 bg-green-500 hover:bg-green-400 active:scale-[0.98] disabled:opacity-50 disabled:hover:bg-green-500 text-black font-bold rounded-full py-2 shadow-lg shadow-green-500/10 transition"
        >
          {{ isSaving ? 'Saving...' : 'Save' }}
        </button>
        <button
          @click="emit('close')"
          class="flex-1 bg-neutral-800 hover:bg-neutral-700 active:scale-[0.98] text-white font-medium rounded-full py-2 transition"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
</template>