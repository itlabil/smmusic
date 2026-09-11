<script setup>
defineProps({
  showIndex: { type: Boolean, default: true },
})
</script>

<template>
  <div class="hidden lg:flex items-center gap-3 px-3 py-2 border-b border-neutral-800 text-neutral-400 text-xs uppercase tracking-wide">
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <span v-if="showIndex" class="w-4 text-right flex-shrink-0">#</span>
      <span class="w-10 flex-shrink-0"></span>
      <span>Title</span>
    </div>
    <div class="w-40 flex-shrink-0">Album</div>
    <div class="w-28 flex-shrink-0">Date Added</div>
    <div class="w-12 flex-shrink-0"></div>
    <div class="w-12 flex-shrink-0 text-right pr-2">Duration</div>
    <div class="w-40 flex-shrink-0"></div>
  </div>
</template>