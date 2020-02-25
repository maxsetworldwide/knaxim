import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const FolderService = {
  create ({ name, content, group }) {
    return ApiService.put(`dir`, { name, content, group }).catch(error => {
      throw buildError('FolderService.create', error)
    })
  },

  list ({ group }) {
    return ApiService.query(`dir`, { group }).catch(error => {
      throw buildError('FolderService.list', error)
    })
  },

  info ({ name, group }) {
    return ApiService.query(`dir/${name}`, { group }).catch(error => {
      throw buildError('FolderService.info', error)
    })
  },

  add ({ name, fid, group }) {
    return ApiService.post(`dir/${name}/content`, { 'id': fid, 'group': group }).catch(error => {
      throw buildError('FolderService.add', error)
    })
  },

  remove ({ name, fid, group }) {
    return ApiService.delete(`dir/${name}/content`, { 'id': fid, 'group': group }).catch(error => {
      throw buildError('FolderService.remove', error)
    })
  },

  erase ({ name, group }) {
    return ApiService.delete(`dir/${name}`, { group }).catch(error => {
      throw buildError('FolderService.erase', error)
    })
  }
}

export default FolderService
