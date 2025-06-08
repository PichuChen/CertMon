import { createRouter, createWebHistory } from 'vue-router'
import CertMonHome from './components/CertMonHome.vue'
import CertMonList from './components/CertMonList.vue'
import CertMonDetail from './components/CertMonDetail.vue'

const routes = [
  { path: '/', name: 'home', component: CertMonHome },
  { path: '/domains', name: 'list', component: CertMonList },
  { path: '/domains/:domain', name: 'detail', component: CertMonDetail, props: true },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
