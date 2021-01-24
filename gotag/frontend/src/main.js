import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import '@mdi/font/css/materialdesignicons.css'

Vue.config.productionTip = false
Vue.filter('prettyTimeDate', (str) => {
  function withZero (num) {
    if (num < 10) {
      return '0' + num
    }
    return String(num)
  }

  const t = new Date(str)
  // this date looks wrong, show dash
  if (isNaN(t.getTime())) {
    return '-'
  }
  // event not finished yet / unknown, show dash
  if (t.getFullYear() == 1) {
    return '-'
  }
  // show proper date
  return withZero(t.getHours()) + ':' + withZero(t.getMinutes()) +
    ' ' + withZero(t.getDate()) + '/' + withZero(t.getMonth() + 1) + '/' + t.getFullYear()
})
const mixin = {
  computed: {
    apiUrl () {
      if (process.env.NODE_ENV === 'production') {
        return ''
      } else {
        return 'http://127.0.0.1:4000'
      }
    },
    lsToken () {
      return 'authToken'
    },
  },
  methods: {
    urlToThumb (file, width) {
      if (file.mime === 'image/gif') {
        return this.apiUrl + '/img/full/' + file.sha256
      }
      return this.apiUrl + '/img/thumbs/' + width + '/' + width + '/' + file.sha256
    },
    isVideo (mime) {
      if (mime === 'video/webm' || mime === 'video/mp4') {
        return true
      }
      return false
    },
  }
}
Vue.mixin(mixin)

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
