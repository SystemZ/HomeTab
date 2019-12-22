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
  return withZero(t.getHours()) + ':' + withZero(t.getMinutes()) +
      ' ' + withZero(t.getDate()) + '/' + withZero(t.getMonth() + 1) + '/' + t.getFullYear();
});

new Vue({
  router,
  store,
  //@ts-ignore
  vuetify,
  render: h => h(App)
}).$mount('#app')
