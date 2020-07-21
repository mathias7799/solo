<template>
  <nav class="navbar navbar-expand-lg navbar-light bg-light">
    <a class="navbar-brand" href="/">
      <span class="logo-wrapper">
        <img
          src="./../assets/logo.svg"
          width="30"
          height="30"
          class="d-inline-block align-top"
          alt="Flexpool Logo"
          loading="lazy"
        />
      </span>
      <span class="ml-2 font-weight-bold">Flexpool SOLO</span>
    </a>
    <button
      class="navbar-toggler"
      type="button"
      data-toggle="collapse"
      data-target="#navbarNavDropdown"
      aria-controls="navbarNavDropdown"
      aria-expanded="false"
      aria-label="Toggle navigation"
    >
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNavDropdown">
      <ul class="navbar-nav">
        <NavbarItem name="Home" isSelected></NavbarItem>
        <NavbarItem name="Blocks"></NavbarItem>
        <NavbarItem name="Node"></NavbarItem>
      </ul>
    </div>
    <span class="navbar-text">{{ block_number }}</span>
  </nav>
</template>

<script>
module.exports = {
  components: {
    NavbarItem: require("./NavbarItem.vue").default
  },
  props: {
    block_number: Number
  }
};
</script>

<style scoped>
.navbar-brand {
  display: flex;
  justify-content: center;
  align-items: center;
}

.logo-wrapper {
  display: inline-flex;
  justify-content: center;
  flex-direction: center;
  background-color: white;
  border-radius: 10px;
  padding: 7px;
}
</style>