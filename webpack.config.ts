import webpack from 'webpack'
import path from 'path'
import { writeFileSync } from 'fs'

import { name, version, description } from './package.json'

console.log('Copying CLI config')
writeFileSync(
  path.resolve(__dirname, 'src', 'config.json'),
  JSON.stringify({ name, version, description }, null, 2),
  { encoding: 'utf-8' }
)
console.log('Finished Copy CLI config')

const config: webpack.Configuration = {
  entry: './src/index.ts',
  target: 'node',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'build')
  },
  resolve: {
    extensions: ['.webpack.js', '.ts', '.js']
  },
  module: {
    rules: [{ test: /\.ts$/, loader: 'ts-loader', exclude: /node_modules/ }]
  }
}

export default config
