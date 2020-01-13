import ApiService from '@/service/api'

const SearchService = {

  folderFiles ({ name, find, group }) {
    return ApiService.query(`dir/${name}/search`, { find, group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  },

  groupFiles ({ gid, find }) {
    return ApiService.query(`group/${gid}/search`, { find }).catch(error => {
      throw new Error(`GroupService ${error}`)
    })
  },

  publicFiles ({ find }) {
    return ApiService.query(`public/search`, { find }).catch(error => {
      throw new Error(`SearchService ${error}`)
    })
  },
  userFiles ({ find }) {
    return ApiService.query(`user/search`, { find }).catch(error => {
      throw new Error(`SearchService ${error}`)
    })
  }
}

export default SearchService
