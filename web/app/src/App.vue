<template>
  <div id="app">
    <Navbar />
    <main>
      <router-view></router-view>
    </main>
  </div>
</template>

<script>
import Navbar from "./components/Navbar.vue";
//import Stats from "./components/stats-components/Stats.vue";

export default {
  name: "App",
  components: {
    Navbar,
    //Stats,
  },
};
</script>

<style lang="scss">
@import "style/_utils.scss";
@import "style/_fonts.scss";
* {
  padding: 0px;
  margin: 0px;
  font-family: "Montserrat", sans-serif;
}

main {
  width: 100%;
  max-width: 980px;
  margin: 0px auto;
  margin-top: 140px;
  min-height: calc(100vh - 500px);
  margin-bottom: 100px;
}

h1.center {
  text-align: center;
  width: 100%;
}

mark {
  background-color: #00000000;
  color: #0069ff;
}

h1.center {
  text-align: center;
}

div.center {
  text-align: center;
}

h1.mt {
  margin-top: 30px;
}

@media (max-width: 980px) {
  main {
    width: 90%;
    margin-top: 130px;
  }
}

input:focus,
select:focus,
textarea:focus,
button:focus {
  outline: none;
}

a {
  color: #0069ff;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

body {
  max-width: 100%;
}

@media (max-width: 980px) {
  main {
    margin-top: 130px;
  }
}

/* Highcharts Title */
tspan {
  font-weight: 600;
}
</style>
