import 'hamburgers/dist/hamburgers.min.css';

import Vue from 'vue'
import HighchartsVue from 'highcharts-vue';

Vue.use(HighchartsVue);

import App from './App.vue'

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
