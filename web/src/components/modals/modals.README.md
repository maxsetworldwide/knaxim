## Creating and using Modals

Modals are components just like everything else, but there are a few things to keep in mind:
- We are essentially creating a component that wraps a `b-modal` and puts our own content inside of it. [There are many options available with these, as posted in the documentation](https://bootstrap-vue.js.org/docs/components/modal/)
- In order to show a b-modal, you can either call `this.$bvModal.show(id)` or call the modal with `$refs` and call its `show()` method. Since using the former method will show multiple modals with the same ID, I would recommend using the latter and including functions called `show()` and `hide()` that wrap `b-modal`'s `show()` and `hide()` functions, then calling them from the parent. Otherwise, make sure to keep your IDs unique.
- Modals can be placed anywhere in your template. I personally have been placing them wherever their trigger button is.
- When styling the modal, you can use the `content-class` and other attributes of `b-modal`, however, if your css is scoped in your component, you will need to make the class `::v-deep`. If you know of or find a better way to do this, update this!
- I have also made a confirmation modal which utilizes the "old" method of dynamically creating a modal and mounting it via a global function call. This is not very extensible at the moment, so making more of these is not recommended. Note that these modals don't utilize props. They only utilize data.


Copyright August 2020 Maxset Worldwide Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
