import Vue from 'vue'
import VueRouter from 'vue-router'
import Tasks from '../views/Tasks.vue'
import Notes from '../views/Notes.vue'

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
]

const router = new VueRouter({
    routes
})

export default router
