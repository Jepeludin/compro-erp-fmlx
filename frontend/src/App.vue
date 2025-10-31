<template>
  <div id="app-layout">
    <router-view />
  </div>
</template>

<script setup>
// Tidak perlu script untuk layout sederhana ini
</script>

<style>
/* Anda bisa pindahkan CSS ini ke src/style.css
  jika Anda mau. Ini untuk memastikan tampilan dasar
  terlihat rapi.
*/
body {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin: 0;
  background-color: #f4f7f6;
}
</style>