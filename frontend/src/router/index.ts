import { createRouter, createWebHistory } from 'vue-router'
import PortalView from '@/views/PortalView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'portal',
      component: PortalView
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('super_token')

  if (to.meta.requiresAuth && !token) {
    next('/')
  } else {
    next()
  }
})

export default router
