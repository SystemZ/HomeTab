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
import Login from '../views/Login.vue'
import Settings from '../views/Settings.vue'

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
