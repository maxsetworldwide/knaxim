import { shallowMount } from '@vue/test-utils'
import FileListBatch from '@/components/file-list-batch'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(FileListBatch, {
    stubs: ['batch-delete', 'b-dropdown-item', 'b-dropdown-divider', 'b-dropdown'],
    propsData: {
      checkedFiles: [{
        id: 'id-abc-123',
        name: 'fakeFile'
      }],
      ...options.props
    },
    methods: {
      ...options.methods
    },
    computed: {
      ...options.computed
    }
  })
}

describe('FileListBatch', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(FileListBatch)).toBe(true)
  })

  /** TODO: Evaluate the use of upstream logic; The following functions could easily
   * be pulled out of upstream code and into their own components, or some components.
   */
  it('emits an add-folder', (done) => {
    const wrapper = shallowMountFa()

    wrapper.vm.newFolder()
    wrapper.vm.$nextTick().then(() => {
      expect(wrapper.emitted()['add-folder']).toBeTruthy()
      done()
    })
  })

  it('emits a favorite', (done) => {
    const wrapper = shallowMountFa()

    wrapper.vm.addFavorite()
    wrapper.vm.$nextTick().then(() => {
      expect(wrapper.emitted().favorite).toBeTruthy()
      done()
    })
  })

  it('emits a share-file', (done) => {
    const wrapper = shallowMountFa()

    wrapper.vm.$emit('share-file')
    wrapper.vm.share()
    wrapper.vm.$nextTick().then(() => {
      expect(wrapper.emitted()['share-file']).toBeTruthy()
      done()
    })
  })
})
