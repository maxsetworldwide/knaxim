import { shallowMount } from '@vue/test-utils'
import FileIcon from '@/components/file-icon'

// API options for test-utils - mount, shallowMount, etc.:
//   https://vue-test-utils.vuejs.org/api

// API options for mount/shallowMount - propsData, data, stubs, etc.:
//   https://vue-test-utils.vuejs.org/api/options.html#context

// Jasmine matchers - toBeTruthy, toBeDefined, etc.
//   https://jasmine.github.io/api/3.5/matchers.html

const shallowMountFa = (options = { props: {}, methods: {}, computed: {} }) => {
  return shallowMount(FileIcon, {
    propsData: {
      extention: '',
      webpage: false,
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

describe('FileIcon', () => {
  it('imports correctly', () => {
    const wrapper = shallowMountFa()
    expect(wrapper.is(FileIcon)).toBe(true)
  })
  const testExtensionWithIcon = function (ext, expectedIcon) {
    it(`uses ${expectedIcon} icon with ${ext} extension`, () => {
      const wrapper = shallowMountFa({ props: { extention: ext } })
      expect(wrapper.findAll('svg').length).toEqual(1)
      const svg = wrapper.find('use')
      expect(svg.attributes('href')).toContain(expectedIcon)
    })
  }
  const iconTests = [
    { ext: 'pdf', expectedIcon: '#pdf2' },
    { ext: 'doc', expectedIcon: '#doc' },
    { ext: 'docx', expectedIcon: '#doc' },
    { ext: 'csv', expectedIcon: '#csv' },
    { ext: 'txt', expectedIcon: '#txt' },
    { ext: 'ppt', expectedIcon: '#ppt' },
    { ext: 'pptx', expectedIcon: '#ppt' },
    { ext: 'xls', expectedIcon: '#xls' },
    { ext: 'xlsx', expectedIcon: '#xls' }
  ]
  for (let test of iconTests) {
    testExtensionWithIcon(test.ext, test.expectedIcon)
  }

  it('uses webpage icon when webpage prop is set', () => {
    const wrapper = shallowMountFa({ props: { extention: '', webpage: true } })
    expect(wrapper.findAll('svg').length).toEqual(1)
    expect(wrapper.find('use').attributes('href')).toContain('#webpage')
  })
  it('utilizes a text fallback with an extension with no existing svg', () => {
    const testExt = 'notexist'
    const wrapper = shallowMountFa({ props: { extention: testExt } })
    expect(wrapper.findAll('svg').length).toEqual(0)
    expect(wrapper.text().toLowerCase()).toContain(testExt.toLowerCase())
  })
})
