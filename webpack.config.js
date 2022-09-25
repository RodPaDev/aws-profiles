"use strict";

const path = require("path");
const nodeExternals = require("webpack-node-externals");
const webpack = require("webpack");

module.exports = {
  target: "node",
  entry: "./src/index.js",
  mode: process.env.NODE_ENV || "development",
  externals: [nodeExternals()],
  module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: "babel-loader",
          options: {
            presets: ["@babel/preset-env"],
          },
        },
      },
    ],
  },
  output: {
    libraryTarget: "commonjs2",
    path: path.join(__dirname, "./dist"),
    filename: "[name].js",
  },
  plugins: [new webpack.LoaderOptionsPlugin({ minimize: true, debug: false })],
};
