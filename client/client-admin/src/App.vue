<script setup>
import {RouterView } from 'vue-router'
</script>

<template>
  <header>
    <div class="header-wrapper">
      <div class="header">
        <div class="logo">Logo</div>
        <div class="nav-menu">
          <nav @click="$router.push('/admin/users')">Users</nav>
          <nav @click="$router.push('/admin/music')">Music</nav>
          <nav @click="$router.push('/admin/films')">Films</nav>
          <nav @click="$router.push('/admin/video')">Video</nav>

        </div>
      </div>
    </div>
  </header>
  <body>
    <div class="content-wrapper">
      <RouterView />
    </div>
  </body>
</template>

<style scoped>

*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.header-wrapper {
  width: 100%;
  height: 90px;
  background-color: #181818;
}

.header {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.logo {
  position: absolute;
  left: 20px;
  cursor: pointer;
}

.nav-menu {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 300px;
}

.nav-menu nav {
  cursor: pointer;
    color: rgb(148, 148, 148);
  transition: all .3s ease-in-out;
}
.nav-menu nav:hover {
  color: whitesmoke;
}

.content-wrapper {
  background-color: #242424;
  width: 100%;
  height: 100vh;
}

</style>
