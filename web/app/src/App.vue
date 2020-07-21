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

input:focus,
select:focus,
textarea:focus,
button:focus {
  outline: none;
}
</style>
