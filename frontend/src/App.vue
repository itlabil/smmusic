<script setup>
import { useRoute } from 'vue-router'
import { computed } from 'vue'
import AppSidebar from '@/components/sidebar/AppSidebar.vue'
import MobileNav from '@/components/sidebar/MobileNav.vue'
import NowPlayingBar from '@/components/player/NowPlayingBar.vue'
import QueuePanel from '@/components/player/QueuePanel.vue'
import { usePlayerStore } from '@/stores/player'

const route = useRoute()
const isAuthPage = computed(() => route.name === 'login')
const player = usePlayerStore()
</script>

<template>
  <div v-if="isAuthPage">
    <router-view />
  </div>

  <div v-else class="h-screen flex flex-col bg-neutral-950">
    <div class="flex flex-1 overflow-hidden">
      <AppSidebar />
      <main class="flex-1 overflow-y-auto">
        <router-view />
      </main>
      <QueuePanel v-if="player.isQueueOpen" />
    </div>
    <NowPlayingBar />
    <MobileNav />
  </div>
</template>