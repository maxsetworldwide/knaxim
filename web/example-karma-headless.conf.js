// Karma configuration
// Generated on Tue Oct 15 2019 21:08:53 GMT+0000 (Coordinated Universal Time)
var webpackConfig = require('../node_modules/@vue/cli-service/webpack.config.js')

// Setup headless chrome executabale.
process.env.CHROME_BIN = require('puppeteer').executablePath()

module.exports = function(config) {
  config.set({

    // Client configuration.
    client: {
      jasmine: {
        random: false
      }
    },

    // base path that will be used to resolve all patterns (eg. files, exclude)
    basePath: '../',

    // frameworks to use
    // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
    frameworks: ['jasmine'],

    // list of files / patterns to load in the browser
    files: [
      'test/**/*.spec.js'
    ],

    // list of files / patterns to exclude
    exclude: [
      '**/*.swp'
    ],

    // preprocess matching files before serving them to the browser
    // available preprocessors: https://npmjs.org/browse/keyword/karma-preprocessor
    preprocessors: {
      'test/setup.js': ['webpack', 'sourcemap'],
      '**/*.spec.js': ['webpack', 'sourcemap']
    },

    webpack: webpackConfig,

    // test results reporter to use
    // possible values: 'dots', 'progress'
    // available reporters: https://npmjs.org/browse/keyword/karma-reporter
    reporters: ['verbose', 'progress'],

    // web server port
    port: 9876,

    proxies: {
      '/api': {
        'target': 'http://localhost:8005/api',
        'changeOrigin': true
      }
    },

    // enable / disable colors in the output (reporters and logs)
    colors: true,

    // level of logging
    // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
    logLevel: config.LOG_INFO,

    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: true,

    // start these browsers
    // available browser launchers: https://npmjs.org/browse/keyword/karma-launcher
    browsers: ['ChromeHeadless'],
    flags: [
      '--disable-web-security',
      '--disable-gpu',
      '--no-sandbox'
    ]

    // Continuous Integration mode
    // if true, Karma captures browsers, runs the tests and exits
    singleRun: true,

    // Concurrency level
    // how many browser should be started simultaneous
    concurrency: Infinity
  })
}
