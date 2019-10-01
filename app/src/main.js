import Vue from 'vue'
import App from './App.vue'
import VueDragRotateResize from '@chenxingyu_o/vue-drag-rotate-resize'
import '@chenxingyu_o/vue-drag-rotate-resize/lib/vue-drag-rotate-resize.css'
import './registerServiceWorker'

Vue.use(VueDragRotateResize)

Vue.config.productionTip = false

new Vue({
  render: h => h(App)
}).$mount('#app')
