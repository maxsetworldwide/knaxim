<script>
export default {
  name: 'error-control',
  props: {
    // Listen for global error messages.
    global: {
      type: Boolean,
      defaul: false
    }
  },

  data () {
    return {
      // TODO: Track error history for debugging purposes...
      errors: [],
      lastError: ''
    }
  },
  mounted () {
    // Listen for and handle any global error messages.
    if (this.global) {
      this.$root.$on('app::set::error', (data) => {
        this.setError(data)
      })
    }
  },

  methods: {
    /* Handle all incoming errors; notify the parent.
     *
     * Use setError in the parent with $refs or anywhere in the slot using
     * the returned slot function.
    */
    setError (error) {
      this.lastError = error
      this.$emit('error', this.lastError)
    }
  },

  render () {
    return this.$scopedSlots.default({
      setError: this.setError,
      lastError: this.lastError
    })
  }
}
</script>
