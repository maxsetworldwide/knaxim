import modul from '@/store/error.module'
import { GET_ERROR, ERROR_LOOP } from '@/store/actions.type'
import { PUSH_ERROR, POP_ERROR, ADD_ERROR_LOOP, RESET_ERROR } from '@/store/mutations.type'
import { testAction } from './util'

describe('Error Store', function () {
  describe('Mutations', function () {
    let m = modul.mutations
    beforeEach(function () {
      this.state = {
        errors: [],
        errorLoop: Promise.resolve(true)
      }
    })
    it('queues errors', function () {
      m[PUSH_ERROR](this.state, 'hello')
      expect(this.state.errors).toEqual(['hello'])
    })
    it('dequeues errors', function () {
      this.state.errors.push('world')
      m[POP_ERROR](this.state, 'world')
      expect(this.state.errors.length).toBe(0)
    })
    it('adds error callback', async function () {
      m[ADD_ERROR_LOOP](this.state, (i) => {
        expect(i).toBeTrue()
        return 'theresa'
      })
      expect(await this.state.errorLoop).toBe('theresa')
    })
    it('restores to empty', function (done) {
      this.state.errors = false
      this.state.errorLoop = Promise.reject(new Error('needs reset'))
      m[RESET_ERROR](this.state)
      expect(this.state.errors).toEqual([])
      this.state.errorLoop.then((i) => {
        expect(i).toBeTrue()
        done()
      })
    })
  })
  describe('Actions', function () {
    let a = modul.actions
    beforeEach(function () {
      this.state = {
        errors: [
          new Error('1'),
          new Error('2')
        ],
        errorLoop: Promise.resolve(true)
      }
    })
    it('returns next error', async function () {
      await testAction(
        a[GET_ERROR],
        {},
        { mutations: [{ type: POP_ERROR, payload: new Error('1') }] },
        this.state
      )
    })
    it('loops through errors', async function () {
      let s = this.state
      let count = 0
      a[ERROR_LOOP]({
        commit: (type, payload) => {
          if (type === ADD_ERROR_LOOP) {
            s.errorLoop = s.errorLoop.then(payload)
          } else if (type !== RESET_ERROR) {
            fail('unexpected commit')
          }
        },
        dispatch: (t) => {
          expect(t).toBe(GET_ERROR)
          return Promise.resolve(s.errors.shift())
        },
        getters: {
          get availableErrors () {
            return s.errors.length > 0
          }
        }
      }, (e) => {
        count++
        expect(e).toBeDefined()
      })
      return s.errorLoop.then((i) => {
        expect(i).toBeTrue()
        expect(count).toBe(2)
        return true
      })
    })
  })
})
