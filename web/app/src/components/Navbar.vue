<template>
  <nav>
    <div class="navbody">
      <div class="header noselect">
        <a class="logo" href="/">
          <img src="../assets/logo.svg" alt="Logo" />
          flexpool SOLO
        </a>
        <button class="hamburger hamburger--spin" type="button">
          <span class="hamburger-box">
            <span class="hamburger-inner"></span>
          </span>
        </button>
      </div>
      <ul class="navul">
        <li>
          <a href="stats">Stats</a>
        </li>
        <li>
          <a href="blocks">Blocks</a>
        </li>
        <li>
          <a href="miners">Node</a>
        </li>
        <li>
          <a href="docs/help" class="help">Docs</a>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script>
import $ from "jquery";
export default {
  name: "Navbar",
  mounted() {
    $(".hamburger").click(function() {
      var nav = $("nav");
      var body = $("body");
      var navul = $(".navbody ul");

      var hamburger = $(".hamburger");
      if (nav.height() <= 100) {
        nav.animate({ height: "100%" }, 150);
        body.css("overflow", "hidden");
        hamburger.addClass("is-active");
        navul.css("visibility", "visible");
      } else {
        nav.css({ height: "100px" });
        body.css("overflow", "visible");
        hamburger.removeClass("is-active");
        navul.css("visibility", "hidden");
      }
    });
  }
};
</script>
<!--
<style lang="scss" scoped>
nav {
  position: fixed;
  width: 100%;
  top: 0;
  height: 100px;
  box-shadow: 0px 10px 30px 0px #0000000a;
  display: flex;
  justify-content: center;
  background-color: white;
  z-index: 99;
}

.navbody {
  width: 100%;
  max-width: 980px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  .logo {
    display: flex;
    align-items: center;
    font-family: "Montserrat", sans-serif;
    font-weight: bold;
    font-size: 25px;
    cursor: pointer;
    height: 100%;
    padding: 0px 20px;
    color: black;
    text-decoration: none;
    img {
      width: 45px;
      height: 45px;
      margin-left: 5px;
      padding-right: 10px;
    }
  }

  ul {
    list-style: none;
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    margin-right: 10px;
    height: 100%;
    align-items: center;
  }

  li {
    font-family: "Montserrat", sans-serif;
    font-weight: 400;
    display: flex;
    height: 100%;
    width: 100px;

    a {
      display: block;
      font-weight: bold;
      color: black;
      text-decoration: none;
      transition: 0.05s;
      font-size: 17px;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
    }
  }
}

.navbody a:hover,
.navbody .logo:hover {
  background-color: #fafafa;
  transition: 0.05s;
}

.navbody li a:active,
.navbody .logo:active {
  background-color: #f5f5f5;
}

.navbody li a.help {
  background-color: #0069ff;
  color: white;
}

.navbody li a.help:hover {
  background-color: #0065fa;
  color: white;
}

.navbody li a.help:active {
  background-color: #0050e5;
}

.navbody li a:hover {
  color: #0069ff;
  transition: 0.05s;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  height: 100px !important;
}

.hamburger {
  visibility: hidden;
  font-size: 40px;
}

@media (max-width: 980px) {
  nav {
    overflow: hidden;
  }

  .navbody .logo {
    padding: 0px 15px;

    img {
      margin-left: 0px;
    }
  }

  .header {
    width: 90%;
  }

  .hamburger {
    visibility: visible;
  }

  .navbody {
    flex-direction: column;
    justify-content: flex-start;

    ul {
      display: block;
      height: 0px;
      margin-right: 0px;
      flex-direction: column;
      width: 100%;
      visibility: hidden;

      li {
        display: flex;
        align-items: center;
        width: 100%;
        height: 80px;
        justify-content: center;
        border-bottom: 1px solid #efefef;

        a {
          font-size: 25px;
          width: 100%;
        }

        :last-child a:hover {
          background-color: #0059ef;
        }

        :first-child {
          border-top: 1px solid #efefef;
        }
      }
    }
  }
}
</style>-->