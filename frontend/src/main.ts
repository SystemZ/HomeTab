import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import retryTimes = jest.retryTimes;
//import 'roboto-fontface/css/roboto/roboto-fontface.css'
//import '@mdi/font/css/materialdesignicons.css'

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
        return '';
      } else {
        return 'http://127.0.0.1:3000';
      }
    },
    lsToken(): string {
      return 'authToken';
    },
    // color templates
    // button colors
    btnPrimary(): string {
      return '#5C9DA0';
    },
    btnSecondary(): string {
      return '#DE6461';
    },
    btnAccent(): string {
      return '#DE9B61';
    },
    //checkbox colors
    checkPrimary(): string {
      return '#DE9B61';
    },
    // chip colors
    chipPrimary(): string {
      return '#009688';
    },
    // is this element dark?
    componentDark(): boolean {
      return true
    },
    // input colors
      inputPrimary(): string {
      return '#F06292';
    },
    // navbar colors
      navbarPrimary(): string {
      return '#4A148C';
      },
    // TODO pagination colors - tables
    // https://stackoverflow.com/questions/58936262/how-to-change-default-styling-of-v-data-table-footer
    // pagination colors
    pagePrimary(): string {
      return '#827717';
    },
    // picker colors
      pickerPrimary(): string {
      return '#673AB7';
    },
    // progress bar colors
      progressPrimary(): string {
      return '#FFB74D';
    },
    // tags color
    tagPrimary(): string {
      return '#FF5722';
    },
    // toolbar colors
      toolbarPrimary(): string {
      return '#673AB7';
      },
    // snackbar colors
      snackbarPrimary(): string {
      return '#FF80AB';
    },
    // switch colors
      switchPrimary(): string {
      return '#827717';
    },
  },
  methods: {
    urlToThumb(file: string, width: string): string {
      // @ts-ignore
      if (file.mime === 'image/gif') {
        // @ts-ignore
        return this.apiUrl + '/img/full/' + file.sha256
      }
      // @ts-ignore
      return this.apiUrl + '/img/thumbs/' + width + '/' + width + '/' + file.sha256
    },
    isVideo(mime: string) {
      return mime === 'video/webm' || mime === 'video/mp4';
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
