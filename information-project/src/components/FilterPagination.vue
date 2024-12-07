<script lang="ts" setup>
import { computed } from 'vue'

const props = defineProps({
  pages: {
    type: Number,
    required: true,
  },
  currentPage: {
    type: Number,
    required: true,
  },
})
const emit = defineEmits(['pageChanged']);
const maxVisible = 10; // Máximo de páginas visibles


const visiblePages = computed(() => {
  const half = Math.floor(maxVisible / 2);
  let start = Math.max(1, props.currentPage - half);
  const end = Math.min(props.pages, start + maxVisible - 1);

  start = Math.max(1, end - maxVisible + 1);

  return Array.from({ length: end - start + 1 }, (_, i) => start + i);
});

const goToPreviousPage = () => {
  if (props.currentPage > 1) {
    emit('pageChanged', props.currentPage - 1);
  }
};
const goToNextPage = () => {
  if (props.currentPage < props.pages) {
    emit('pageChanged', props.currentPage + 1);
  }
};
</script>
<template>
  <div class="flex items-center justify-between space-x-1">
    <button @click="goToPreviousPage" :disabled="currentPage === 1"
      class="rounded-md border border-slate-300 py-2 px-3 text-center text-sm transition-all shadow-sm hover:shadow-lg text-slate-600 hover:text-white hover:bg-slate-800 hover:border-slate-800 focus:text-white focus:bg-slate-800 focus:border-slate-800 active:border-slate-800 active:text-white active:bg-slate-800 disabled:pointer-events-none disabled:opacity-50 disabled:shadow-none ml-2">
      &lt;
    </button>
    <button v-for="page in visiblePages" :key="page" @click="$emit('pageChanged', page)" :class="[
      'rounded-md border border-slate-300 py-2 px-3 text-center text-sm transition-all shadow-sm hover:shadow-lg text-slate-600 hover:text-white hover:bg-slate-800 hover:border-slate-800 focus:text-white focus:bg-slate-800 focus:border-slate-800 active:border-slate-800 active:text-white active:bg-slate-800 disabled:pointer-events-none disabled:opacity-50 disabled:shadow-none ml-2',
      { 'bg-slate-600 text-white': page === currentPage }
    ]">
      {{ page }}
    </button>
    <button @click="goToNextPage" :disabled="currentPage === pages"
      class="min-w-9 rounded-md border border-slate-300 py-2 px-3 text-center text-sm transition-all shadow-sm hover:shadow-lg text-slate-600 hover:text-white hover:bg-slate-800 hover:border-slate-800 focus:text-white focus:bg-slate-800 focus:border-slate-800 active:border-slate-800 active:text-white active:bg-slate-800 disabled:pointer-events-none disabled:opacity-50 disabled:shadow-none ml-2">
      &gt;
    </button>
  </div>

</template>
