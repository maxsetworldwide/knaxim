import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const OwnerService = {
  id (id) {
    return ApiService.get(`owner/id/${id}`).catch(error => {
      throw buildError('OwnerService.id', error)
    })
  },
  name (name) {
    return ApiService.get(`owner/name/${name}`).catch(error => {
      throw buildError('OwnerService.name', error)
    })
  }
}

export default OwnerService
