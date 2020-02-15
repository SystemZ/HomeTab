import Vue from 'vue'
import VueRouter from 'vue-router'
import Files from '../views/Files.vue'
import Tags from '../views/Tags.vue'
import Login from '../views/Login.vue'

Vue.use(VueRouter)

const routes = [
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
    path: '/login',
    name: 'login',
    component: Login
  },
]

const router = new VueRouter({
  routes
})

export default router
