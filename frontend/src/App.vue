<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN">
    <n-modal-provider>
      <n-loading-bar-provider>
        <n-message-provider>
          <router-view />
        </n-message-provider>
      </n-loading-bar-provider>
    </n-modal-provider>
  </n-config-provider>
</template>

<script setup>
import { zhCN, dateZhCN } from 'naive-ui'

</script>

<style scoped></style>
