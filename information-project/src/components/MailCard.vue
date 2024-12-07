<script setup lang="ts">
import type { IEmail } from '@/types/types';
const props = defineProps<{ email: IEmail }>();

</script>
<template>
    <section class="max-w-2xl px-6 py-8 mx-auto bg-white dark:bg-gray-900">
        <main class="mt-8">
            <h2 class="text-gray-700 dark:text-gray-200 font-bold">{{ props.email.Subject }}</h2>
            <p class="mt-2 leading-loose text-gray-600 dark:text-gray-300">
                {{ props.email.Body }}
            </p>
        </main>
    </section>
</template>