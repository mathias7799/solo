// Preload CSS
import 'hamburgers/dist/hamburgers.min.css';

// Vue
import Vue from 'vue'

// Vue plugins
import HighchartsVue from 'highcharts-vue';
import VueRouter from 'vue-router';

// Components
import App from './App.vue'
import Stats from './components/stats-components/Stats.vue';

Vue.use(HighchartsVue);
Vue.use(VueRouter);



const router = new VueRouter({
  routes: [
    { "path": "/", component: Stats }
  ],
  mode: "history"
})


Vue.config.productionTip = false

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
