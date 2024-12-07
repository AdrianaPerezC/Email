<script lang="ts" setup>
const props = defineProps({
  information: String, error:
    String
})
</script>
<template>
  <div class="text-xs italic mt-4">
    <div v-if="information" class="text-green-700">
      <pre>{{ props.information }}</pre>
    </div>

    <div v-if="error">
      <p style="color: red;">{{ props.error }}</p>
    </div>
  </div>
</template>