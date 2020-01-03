import Vue from 'vue'
import VueRouter from 'vue-router'
import Tasks from '../views/Tasks.vue'
import Notes from '../views/Notes.vue'
import Note from '../views/Note.vue'
import Login from '../views/Login.vue'

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
        path: '/login',
        name: 'login',
        component: Login
    },
]

const router = new VueRouter({
    routes
})

export default router
