<script setup>
import { ref, watch } from 'vue'
import { useModalStore } from '@/stores/modal'

const modal = useModalStore()
const inputValue = ref('')

watch(
  () => modal.isOpen,
  (open) => {
    if (open && modal.type === 'prompt') {
      inputValue.value = modal.defaultValue
    }
  },
)

function handleConfirm() {
  modal.resolve(modal.type === 'prompt' ? inputValue.value.trim() : true)
}

function handleCancel() {
  modal.resolve(modal.type === 'prompt' ? null : false)
}
</script>

<template>
  <div
    v-if="modal.isOpen"
    class="fixed inset-0 bg-black/70 flex items-center justify-center z-[60] p-4"
    @click.self="handleCancel"
  >
    <div class="bg-neutral-900 rounded-lg p-6 w-full max-w-sm">
      <h3 class="text-white font-semibold mb-2">{{ modal.title }}</h3>
      <p v-if="modal.message" class="text-neutral-400 text-sm mb-4">{{ modal.message }}</p>

      <input
        v-if="modal.type === 'prompt'"
        v-model="inputValue"
        type="text"
        :placeholder="modal.placeholder"
        autofocus
        @keyup.enter="handleConfirm"
        class="w-full bg-neutral-800 text-white rounded px-3 py-2 outline-none focus:ring-2 focus:ring-green-500 mb-4"
      />

      <div class="flex gap-3">
        <button
          @click="handleConfirm"
          :class="[
            'flex-1 font-bold rounded-full py-2 active:scale-[0.98] transition',
            modal.danger
              ? 'bg-red-500 hover:bg-red-400 text-white shadow-lg shadow-red-500/10'
              : 'bg-green-500 hover:bg-green-400 text-black shadow-lg shadow-green-500/10',
          ]"
        >
          {{ modal.type === 'prompt' ? 'OK' : (modal.danger ? 'Delete' : 'OK') }}
        </button>
        <button
          @click="handleCancel"
          class="flex-1 bg-neutral-800 hover:bg-neutral-700 active:scale-[0.98] text-white font-medium rounded-full py-2 transition"
        >
          Cancel
        </button>
      </div>
    </div>
  </div>
</template>