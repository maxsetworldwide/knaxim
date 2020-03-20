import { LOAD_FOLDERS, LOAD_FOLDER, PUT_FILE_FOLDER, REMOVE_FILE_FOLDER, HANDLE_SERVER_STATE, LOAD_SERVER } from '@/store/actions.type'
import { FOLDER_LOADING, SET_FOLDER, FOLDER_ADD, FOLDER_REMOVE, ACTIVATE_GROUP, ACTIVATE_FOLDER, DEACTIVATE_FOLDER, PUSH_ERROR } from '@/store/mutations.type'
import { testAction } from './util'
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import modul from '@/store/folder.module'

describe('Folder Store', function () {
  beforeAll(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
    // this.getters = {}
    // for (const getter in modul.getter)
  })
  beforeEach(function () {
    this.state = JSON.parse(JSON.stringify(modul.state))
  })
  describe('Mutations', function () {
    const m = modul.Mutations
    it('clears active folders on changing groups', function () {
      this.state.active.push('to be removed')
      m[ACTIVATE_GROUP](this.state)
      expect(this.state.active).toEqual([])
    })
    it('activates a folder', function () {
      this.state.active.push('first').push('second')
      m[ACTIVATE_FOLDER](this.state, 'second')
      expect(this.state.active).toEqual(['second', 'first'])
    })
    it('deactivates a folder', function () {
      this.state.active.push('to be removed')
      m[DEACTIVATE_FOLDER](this.state, 'to be removed')
      expect(this.state.active).toEqual([])
    })
    it('adjusts loading state', function () {
      m[FOLDER_LOADING](this.state, 5)
      expect(this.state.loading).toBe(5)
      m[FOLDER_LOADING](this.state, -5)
      expect(this.state.loading).toBe(0)
    })
    it('sets folder', function () {
      m[SET_FOLDER](this.state, {
        name: 'userfolder',
        files: ['1', '2']
      })
      m[SET_FOLDER](this.state, {
        name: 'groupFolder',
        files: ['3', '4'],
        group: 'group'
      })
      expect(this.state.user).toEqual({
        userfolder: ['1', '2']
      })
      expect(this.state.group).toEqual({
        group: {
          groupFolder: ['3', '4']
        }
      })
    })
    it('adds to a folder', function () {
      m[FOLDER_ADD](this.state, {
        name: 'folder',
        fid: 'fid'
      })
      m[FOLDER_ADD](this.state, {
        name: 'folder',
        group: 'group',
        fid: 'fid'
      })
      expect(this.state.user).toEqual({
        folder: ['fid']
      })
      expect(this.state.group).toEqual({
        group: {
          folder: ['fid']
        }
      })
    })
    it('removes from folder', function () {
      this.state.user.folder = ['fid', 'other']
      this.state.group.group = {
        folder: ['fid', 'other']
      }
      m[FOLDER_REMOVE](this.state, {
        name: 'folder',
        fid: 'fid'
      })
      m[FOLDER_REMOVE](this.state, {
        group: 'group',
        name: 'folder',
        fid: 'fid'
      })
      expect(this.state.user).toEqual({
        folder: ['other']
      })
      expect(this.state.group).toEqual({
        group: {
          folder: ['other']
        }
      })
    })
  })
  describe('Actions', function () {})
})
