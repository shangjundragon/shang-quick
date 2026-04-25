<template>
  <HelloWorld/>
</template>


<script setup>
import HelloWorld from './components/HelloWorld.vue'
import {onMounted} from "vue";
import request from "@/utils/request.jsx";

onMounted(() => {
  request({
    method: 'get',
    url: '/health'
  }).then(res => {
    console.log('res', res)
  })
})
</script>
