<template>
  <div class="flex items-end gap-0.5 h-4 w-4">
    <span class="w-1 bg-green-500 rounded-sm animate-eq-bar1"></span>
    <span class="w-1 bg-green-500 rounded-sm animate-eq-bar2"></span>
    <span class="w-1 bg-green-500 rounded-sm animate-eq-bar3"></span>
  </div>
</template>

<style scoped>
@keyframes eq-bar1 {
  0%, 100% { height: 30%; }
  50% { height: 100%; }
}
@keyframes eq-bar2 {
  0%, 100% { height: 100%; }
  50% { height: 40%; }
}
@keyframes eq-bar3 {
  0%, 100% { height: 60%; }
  50% { height: 90%; }
}
.animate-eq-bar1 {
  animation: eq-bar1 0.8s ease-in-out infinite;
}
.animate-eq-bar2 {
  animation: eq-bar2 0.7s ease-in-out infinite;
}
.animate-eq-bar3 {
  animation: eq-bar3 0.9s ease-in-out infinite;
}
</style>