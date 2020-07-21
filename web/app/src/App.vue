<template>
  <div id="app">
    <Navbar />
  </div>
</template>

<script>
import Navbar from "./components/Navbar.vue";

export default {
  name: "App",
  components: {
    Navbar
  }
};
</script>

<style lang="scss" scoped>
* {
  font-family: "Open Sans", sans-serif;
}

main {
  max-width: 1000px;
  margin: 0px auto;
}
</style>
