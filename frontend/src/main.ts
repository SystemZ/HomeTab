import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false
Vue.filter('prettyTimeDate', (str: string) => {
  function withZero(num: number): string {
    if (num < 10) {
      return '0' + num;
    }
    return String(num);
  }

  const t = new Date(str);
  // this date looks wrong, show dash
  if (isNaN(t.getTime())) {
    return "-"
  }
  // event not finished yet / unknown, show dash
  if (t.getFullYear() == 1) {
    return "-"
  }
  // show proper date
  return withZero(t.getHours()) + ':' + withZero(t.getMinutes()) +
      ' ' + withZero(t.getDate()) + '/' + withZero(t.getMonth() + 1) + '/' + t.getFullYear();
});
const mixin = {
  computed: {
    apiUrl(): string {
      if (process.env.NODE_ENV === 'production') {
        return 'https://tasktab.lvlup.pro';
      } else {
        return 'http://127.0.0.1:3000';
      }
    },
    lsToken(): string {
      return 'authToken';
    },
  }
}
Vue.mixin(mixin);

new Vue({
  router,
  store,
  //@ts-ignore
  vuetify,
  render: h => h(App)
}).$mount('#app')
