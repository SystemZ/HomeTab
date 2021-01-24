import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Files from '../views/Files.vue'
import Tags from '../views/Tags.vue'
import Scan from '../views/Scan.vue'
import Login from '../views/Login.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/files',
    name: 'files',
    component: Files
  },
  {
    path: '/tags',
    name: 'tags',
    component: Tags
  },
  {
    path: '/scan',
    name: 'scan',
    component: Scan
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
]

const router = new VueRouter({
  routes
})

export default router
