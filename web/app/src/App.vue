<template>
  <div id="app">
    <Navbar />
    <main>
      <!--<Charts />-->
    </main>
  </div>
</template>

<script>
import Navbar from "./components/Navbar.vue";
//import Charts from "./components/Charts.vue";

export default {
  name: "App",
  components: {
    Navbar
    //Charts
  }
};
</script>

<style lang="scss">
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
  main {
    margin-top: 130px;
  }

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
      }
    }
  }

  .navbody li:last-child a:hover {
    background-color: #0059ef;
  }

  .navbody ul li:first-child {
    border-top: 1px solid #efefef;
  }
}
</style>
