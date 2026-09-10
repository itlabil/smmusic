import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import LoginView from '@/views/LoginView.vue'
import HomeView from '@/views/HomeView.vue'
import UploadView from '@/views/UploadView.vue'
import LikedSongsView from '@/views/LikedSongsView.vue'
import PlaylistDetailView from '@/views/PlaylistDetailView.vue'
import AdminUsersView from '@/views/AdminUsersView.vue'
import SearchView from '@/views/SearchView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/upload',
      name: 'upload',
      component: UploadView,
      meta: { requiresAuth: true },
    },
    {
      path: '/liked',
      name: 'liked',
      component: LikedSongsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/playlists/:id',
      name: 'playlist-detail',
      component: PlaylistDetailView,
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/users',
      name: 'admin-users',
      component: AdminUsersView,
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/search',
      name: 'search',
      component: SearchView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login' }
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return { name: 'home' }
  }

  if (to.name === 'login' && authStore.isAuthenticated) {
    return { name: 'home' }
  }
})

export default router