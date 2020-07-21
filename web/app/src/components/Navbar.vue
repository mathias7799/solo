<template>
  <nav>
    <div class="nav-container">
      <div class="navigation">
        <div class="logo">
          <img src="../assets/logo.svg" alt="Flexpool Logo" />
          <span>Flexpool SOLO</span>
        </div>
      </div>
      <div class="navbar-end">
        <div class="navbar-items">
          <NavbarItem label="Home"></NavbarItem>
        </div>
        <NavbarNodeStatus isActive isSyncing blockNumber="1234" />
      </div>
    </div>
  </nav>
</template>

<script>
import NavbarItem from "./NavbarItem.vue";
import NavbarNodeStatus from "./NavbarNodeStatus.vue";

export default {
  name: "Navbar",
  components: { NavbarItem, NavbarNodeStatus }
};
</script>

<style lang="scss" scoped>
nav {
  height: 60px;
  display: flex;
  align-content: center;
  justify-content: center;
  flex-direction: row;

  .nav-container {
    width: 100%;
    max-width: 1000px;
    display: flex;
    justify-content: space-between;

    .navigation {
      display: flex;
      justify-content: center;
      align-content: center;
      flex-direction: row;

      .logo {
        display: flex;
        justify-content: space-between;
        align-items: center;
        img {
          width: 50px;
          width: 50px;
        }

        span {
          font-weight: 700;
          font-size: 20px;
          margin-left: 20px;
        }
      }

      ul {
        display: flex;
        li {
          margin-left: 5px;
          list-style: none;
        }
      }
    }
  }

  .navbar-end {
    display: flex;
    justify-content: center;

    .navbar-items {
      margin-left: 25px;
      display: flex;
      align-items: center;

      a {
        font-weight: 500;
        font-size: 16px;
      }
    }

    .node-status {
      margin-left: 45px;
      display: flex;
      align-items: center;

      .block-number {
        color: #777;
      }
    }
  }
}
</style>