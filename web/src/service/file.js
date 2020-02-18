import Vue from 'vue'
import ApiService from '@/service/api'

const FileService = {
  slice ({ fid, start, end }) {
    return ApiService.get(`file/${fid}/slice/${start}/${end}`).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  search ({ fid, start, end, find }) {
    return ApiService.query(`file/${fid}/search/${start}/${end}`, {
      find
    }).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  create ({ file, folder, group }) {
    return ApiService.put(`file`, {
      group: group,
      file: file,
      dir: folder
    }).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  webpage ({ url, group, folder }) {
    return ApiService.put(`file/webpage`, {
      url: url,
      group: group,
      dir: folder
    }).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  rename ({ fid, name }) {
    return ApiService.post(`record/${fid}/name`, { name }).catch(error => {
      throw new Error(`FileService ${error}`)
    })
  },

  info ({ fid }) {
    return ApiService.get(`file/${fid}`).catch(error => {
      throw new Error(`FileService ${error}`)
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
      throw new Error(`FileService ${error}`)
    })
  },

  list ({ shared, gid }) {
    return ApiService.query(`record${shared ? '/view' : ''}`, {
      group: gid
    }).catch(error => {
      throw new Error(`FilesService->list: ${error}`)
    })
  }
}

export default FileService
