import Vue from 'vue'
import App from './App'
import router from './router'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import 'font-awesome-webpack'

Vue.use(Vuetify)

Vue.config.productionTip = false

// config
Vue.mixin({
  created: function () {
    this.$config = {
      apiUrl: 'http://localhost:3000/api/v1/'
    }
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
