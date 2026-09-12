import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useModalStore = defineStore('modal', () => {
  const isOpen = ref(false)
  const type = ref('confirm') // 'confirm' | 'prompt'
  const title = ref('')
  const message = ref('')
  const placeholder = ref('')
  const defaultValue = ref('')
  const danger = ref(false)
  let resolvePromise = null

  function confirm({ title: t, message: m, danger: d = false }) {
    type.value = 'confirm'
    title.value = t
    message.value = m
    danger.value = d
    isOpen.value = true
    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }

  function prompt({ title: t, message: m = '', placeholder: p = '', defaultValue: dv = '' }) {
    type.value = 'prompt'
    title.value = t
    message.value = m
    placeholder.value = p
    defaultValue.value = dv
    isOpen.value = true
    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }

  function resolve(value) {
    isOpen.value = false
    if (resolvePromise) {
      resolvePromise(value)
      resolvePromise = null
    }
  }

  return { isOpen, type, title, message, placeholder, defaultValue, danger, confirm, prompt, resolve }
})