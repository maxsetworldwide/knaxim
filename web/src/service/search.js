import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const SearchService = {

  folderFiles ({ name, find, group }) {
    return ApiService.query(`dir/${name}/search`, { find, group }).catch(error => {
      throw buildError('SearchService.folderFiles', error)
    })
  },

  groupFiles ({ gid, find }) {
    return ApiService.query(`group/${gid}/search`, { find }).catch(error => {
      throw buildError('SearchService.groupFiles', error)
    })
  },

  publicFiles ({ find }) {
    return ApiService.query(`public/search`, { find }).catch(error => {
      throw buildError('SearchService.publicFiles', error)
    })
  },
  userFiles ({ find }) {
    return ApiService.query(`user/search`, { find }).catch(error => {
      throw buildError('SearchService.userFiles', error)
    })
  }
}

export default SearchService
