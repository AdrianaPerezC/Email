<script setup lang="ts">
import DataTable from '@/components/DataTable.vue';
</script>

<template>
  <header class="flex p-4 bg-gray-100 justify-between items-center">
    <nav class="flex items-center w-full justify-between">
      <div class="flex items-center">
      </div>
    </nav>
  </header>
  <div class="flex flex-col items-center w-full p-3">
    <div class="flex justify-between">
      <h1 class="text-xl font-bold text-blue-500 mr-72">Emails:</h1>
    </div>
    <DataTable></DataTable>
  </div>
</template>