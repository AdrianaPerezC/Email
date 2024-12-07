import { apiCall } from '@/services/utils'
import { ref } from 'vue'
import type { IResponseInformation,ISearchRequest,IEmailInformation } from '@/types/types';

export function useInformation() {
  const information = ref<IEmailInformation[]>([])
  const pages=ref(0)

  async function loadInformationOrder(request: ISearchRequest) {
    try {
      const response = await apiCall('/', { method: 'POST', data: request })
      if (response.error) {
        throw new Error('Error al cargar los emails ordenados.');
      }
      response.value = response as IResponseInformation
      pages.value=response.hits.total.value
      information.value = response.hits.hits as IEmailInformation[]

      return { success: true, response };
    } catch (error) {
      console.error("Error al cargar los emails: ", error);
      return { success: false, error: (error as Error).message };
    }
  }

  return {  information, pages, loadInformationOrder }
}
