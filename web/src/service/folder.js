import ApiService from '@/service/api'

const FolderService = {
  create ({ name, content, group }) {
    return ApiService.put(`dir`, { name, content, group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  },

  list ({ group }) {
    return ApiService.query(`dir`, { group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  },

  info ({ name, group }) {
    return ApiService.query(`dir/${name}`, { group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  },

  add ({ name, fid, group }) {
    return ApiService.post(`dir/${name}/content`, { 'id': fid, 'group': group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  },

  remove ({ name, fid, group }) {
    return ApiService.delete(`dir/${name}/content`, { 'id': fid, 'group': group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  },

  erase ({ name, group }) {
    return ApiService.delete(`dir/${name}`, { group }).catch(error => {
      throw new Error(`FolderService ${error}`)
    })
  }
}

export default FolderService
