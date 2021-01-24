import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        loggedIn: false,
    },
    mutations: {
        userLogged(state, newVal) {
            state.loggedIn = newVal
        },
    },
    actions: {
        setLoggedIn({dispatch, state, commit}) {
            commit('userLogged', true)
        },
        setLoggedOut({dispatch, state, commit}) {
            commit('userLogged', false)
        },
    },
    modules: {}
})
