<template>
  <n-layout has-sider style="height: 100vh">
    <Sidebar v-model:collapsed="collapsed" />
    <n-layout>
      <Navbar @toggle-collapse="collapsed = !collapsed" />
      <AppMain />
    </n-layout>
  </n-layout>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'
import AppMain from './components/AppMain.vue'

const collapsed = ref(false)
</script>
