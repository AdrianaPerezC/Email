<script lang="ts" setup>
import { computed, ref, onMounted} from 'vue';
import type { IEmail, ISearchRequest, InputSearch } from '../types/types';
import { NUMBER_ROWS } from '../types/types';
import FilterType from '@/components/FilterType.vue';
import FilterPagination from './FilterPagination.vue';
import { useInformation } from '@/services/informationEmail';
import MessageInformation from './MessageInformation.vue';
import MailCard from './MailCard.vue';

const searchFilter = ref('');
const selectedEmail = ref({} as IEmail);
const orderFilter = ref({} as ISearchRequest)
orderFilter.value.size = NUMBER_ROWS
const currentPage = ref(1);
const sortBy = ref('')
const sortOrder = ref('asc')
const informationMessage = ref('');
const errorMessage = ref('');
const { information, pages, loadInformationOrder } = useInformation()

const filterItems = computed(() => {
    let items = information.value
    if (searchFilter.value !== '') {
        items = items.filter(item => item._source.Subject.includes(searchFilter.value) || item._source.To.includes(searchFilter.value) || item._source.From.includes(searchFilter.value))
    }
    return items
})

const sortData = async (column: string) => {
    if (sortBy.value === column) {
        sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
    } else {
        sortBy.value = column;
        sortOrder.value = 'asc';
    }

    if (!orderFilter.value.from) {
        orderFilter.value.from = 0
        orderFilter.value.query = { term: "a" }
    }
    orderFilter.value.sort_fields = Order(sortOrder.value) + sortBy.value
    await loadInformationOrder(orderFilter.value);
}
const Order = (orderBy: string) => {
    if (orderBy == 'asc') {
        return '+'
    } else {
        return '-'
    }
}
const filterSearch = async (search: InputSearch) => {
    const informationSearch: ISearchRequest= {
        query: {
            term: search.word
        },
        from: 0,
        size: NUMBER_ROWS,
        sort_fields: '-subject'
    }
    await loadInformationOrder(informationSearch);
    currentPage.value = 1
    orderFilter.value = informationSearch
    sortBy.value = 'subject'
    sortOrder.value = 'desc'
    informationMessage.value = ''
    selectedEmail.value = {} as IEmail
}
const handlePageChanged = async (page: number) => {
    currentPage.value = page
    loadInformationManagement(page)
}

const selectEmail = (email: IEmail) => {
    selectedEmail.value = email
}

onMounted(() => {
    sortData('subject')
}
);
async function loadInformationManagement(page: number) {
    orderFilter.value.from = page * NUMBER_ROWS
    const result = await loadInformationOrder(orderFilter.value);
    if (!result.success) {
        errorMessage.value = result.error || "";
        informationMessage.value = '';
    }
}
</script>

<template>
    <MessageInformation :information="informationMessage" :error="errorMessage"></MessageInformation>
    <div class="max-w-3xl mx-auto flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
        <div class="flex-1">
            <div class="relative overflow-x-auto shadow-md sm:rounded-lg">
                <div class="p-4">
                    <label for="table-search" class="sr-only">Search</label>
                    <div class="relative mt-1">
                        <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                            <svg class="w-5 h-5 text-gray-500 dark:text-gray-400" fill="currentColor"
                                viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                                <path fill-rule="evenodd"
                                    d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                                    clip-rule="evenodd"></path>
                            </svg>
                        </div>
                        <FilterType @filter="filterSearch"></FilterType>
                    </div>
                </div>
                <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                    <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                        <tr>
                            <th @click="sortData('subject')" scope="col" class="px-6 py-3 cursor-pointer">
                                Subject
                                <span v-if="sortBy === 'subject'">
                                    <i v-if="sortOrder === 'asc'" class="ml-2">↑</i>
                                    <i v-else class="ml-2">↓</i>
                                </span>
                            </th>
                            <th @click="sortData('from')" scope="col" class="px-6 py-3 cursor-pointer">
                                From
                                <span v-if="sortBy === 'from'">
                                    <i v-if="sortOrder === 'asc'" scope="col" class="ml-2">↑</i>
                                    <i v-else class="ml-2">↓</i>
                                </span>
                            </th>
                            <th @click="sortData('to')" scope="col" class="px-6 py-3 cursor-pointer">
                                To
                                <span v-if="sortBy === 'to'">
                                    <i v-if="sortOrder === 'asc'" class="ml-2">↑</i>
                                    <i v-else class="ml-2">↓</i>
                                </span>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in filterItems" :key="item._source.MessageID" @click="selectEmail(item._source)"
                            class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 cursor-pointer transition">
                            <td class="px-6 py-4">
                                {{ item._source.Subject }}
                            </td>
                            <td class="px-6 py-4">
                                {{ item._source.From }}
                            </td>
                            <td class="px-6 py-4">
                                {{ item._source.To }}
                            </td>
                        </tr>
                    </tbody>
                </table>
                <FilterPagination :pages="pages" :currentPage="currentPage" @pageChanged="handlePageChanged">
                </FilterPagination>
            </div>
        </div>

        <div class="flex-none w-full sm:w-96 max-h-screen overflow-y-auto border rounded shadow-sm">
            <MailCard v-if="selectedEmail" :email="selectedEmail"></MailCard>
        </div>
    </div>
</template>