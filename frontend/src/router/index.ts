import Vue from 'vue'
import VueRouter from 'vue-router'
import Tasks from '../views/Tasks.vue'
import Notes from '../views/Notes.vue'
import Note from '../views/Note.vue'
import Counters from '../views/Counters.vue'
import Counter from '../views/Counter.vue'
import Log from '../views/Log.vue'
import Events from '../views/Events.vue'
import Devices from '../views/Devices.vue'
import Pantry from '../views/Pantry.vue'
import Login from '../views/Login.vue'
import Settings from '../views/Settings.vue'
import Scan from '../views/Scan.vue'
import Files from '../views/Files.vue'
import Tags from '../views/Tags.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'tasks-home',
    component: Tasks
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: Tasks
  },
  {
    path: '/notes',
    name: 'notes',
    component: Notes
  },
  {
    path: '/note/:id',
    name: 'note',
    component: Note
  },
  {
    path: '/counters',
    name: 'counters',
    component: Counters
  },
  {
    path: '/counter/:id',
    name: 'counter',
    component: Counter
  },
  {
    path: '/log',
    name: 'log',
    component: Log
  },
  {
    path: '/events',
    name: 'events',
    component: Events
  },
  {
    path: '/devices',
    name: 'devices',
    component: Devices
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
    path: '/pantry',
    name: 'pantry',
    component: Pantry
  },
  {
    path: '/settings',
    name: 'settings',
    component: Settings
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
