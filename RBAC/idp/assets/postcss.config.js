module.exports = {
  plugins: [
    require('tailwindcss'),
    require('autoprefixer'),
    require('cssnano')({
      preset: 'default'
    }),
    require('@fullhuman/postcss-purgecss')({
      content: [
        './templates/*.html',
      ],
      whitelistPatterns: [/^opened/,/^ck-editor/,],
      defaultExtractor: content => content.match(/[A-Za-z0-9-_:/]+/g) || []
    })
  ]
}
