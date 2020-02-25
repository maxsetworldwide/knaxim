import Vue from 'vue'
import ApiService from '@/service/api'
import { buildError } from '@/service/util'

const FileService = {
  slice ({ fid, start, end }) {
    return ApiService.get(`file/${fid}/slice/${start}/${end}`).catch(error => {
      throw buildError('FileService.slice', error)
    })
  },

  search ({ fid, start, end, find }) {
    return ApiService.query(`file/${fid}/search/${start}/${end}`, {
      find
    }).catch(error => {
      throw buildError('FileService.search', error)
    })
  },

  create ({ file, folder, group }) {
    return ApiService.put(`file`, {
      group: group,
      file: file,
      dir: folder
    }).catch(error => {
      throw buildError('FileService.create', error)
    })
  },

  webpage ({ url, group, folder }) {
    return ApiService.put(`file/webpage`, {
      url: url,
      group: group,
      dir: folder
    }).catch(error => {
      throw buildError('FileService.webpage', error)
    })
  },

  rename ({ fid, name }) {
    return ApiService.post(`record/${fid}/name`, { name }).catch(error => {
      throw buildError('FileService.rename', error)
    })
  },

  info ({ fid }) {
    return ApiService.get(`file/${fid}`).catch(error => {
      throw buildError('FileService.info', error)
    })
  },

  downloadURL ({ fid }) {
    return Vue.axios.defaults.baseURL + `/file/${fid}/download`
  },

  viewURL ({ fid }) {
    return Vue.axios.defaults.baseURL + `/file/${fid}/view`
  },

  erase ({ fid }) {
    return ApiService.delete(`file/${fid}`).catch(error => {
      throw buildError('FileService.erase', error)
    })
  },

  list ({ shared, gid }) {
    return ApiService.query(`record${shared ? '/view' : ''}`, {
      group: gid
    }).catch(error => {
      throw buildError('FileService.list', error)
    })
  }
}

export default FileService
