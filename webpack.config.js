const Glob = require("glob");
const path = require('path');

const Webpack = require("webpack");

const CopyWebpackPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const ManifestPlugin = require("webpack-manifest-plugin");

const CleanObsoleteChunks = require('webpack-clean-obsolete-chunks');
// const { CleanWebpackPlugin } = require('clean-webpack-plugin');

const UglifyJsPlugin = require("uglifyjs-webpack-plugin");
const LiveReloadPlugin = require('webpack-livereload-plugin');


// 这是修改后的配置文件，
// 1. 所有入口都在 js/css 下，所以入口文件都是 js 结尾，引用的文件可以是 jsx, ts, tsx; 
// 2. 不会 copy 任何内容到 并且 ./public/assets 目录下，需要手工 copy；
// 3. 支持 react，在 .babelrc 中配置 presets： { "presets": ["@babel/preset-env", "@babel/preset-react"] }


// let msg = path.resolve(__dirname, 'assets/js/src');
// console.log("xlog: ", msg);

let configurator;
configurator = {
  entries: function () {
    const entries = {
      application: [
        './node_modules/jquery-ujs/src/rails.js',
        './assets/css/application.scss',
      ],
    };

    // only files in ./assets/js or css folder, refer to CopyWebpackPlugin ignore option
    Glob.sync("./assets/{js,css}/**/*.*").forEach((entry) => {
      if (entry === './assets/css/application.scss') {
        return
      }

      // if not these files, ignore
      if ((/(ts|js|tsx|jsx|s[ac]ss|go)$/i).test(entry) == false) {
        return
      }

      // ignore files name start with underscore
      let key = entry.replace(/(\.\/assets\/(src|js|jsx|css|go)\/)|\.(ts|js|tsx|jsx|s[ac]ss|go)/g, '')
      if (key.startsWith("_")) {
        return
      }

      if (entries[key] == null) {
        entries[key] = [entry];
        return
      }

      entries[key].push(entry)
    })
    return entries
  },

  plugins() {
    var plugins = [
      new Webpack.ProgressPlugin(),
      new CleanObsoleteChunks(),
      // new CleanWebpackPlugin(),
      new Webpack.ProvidePlugin({$: "jquery", jQuery: "jquery"}),
      new MiniCssExtractPlugin({filename: "[name].[contenthash].css"}),
      // new CopyWebpackPlugin([{from: "./assets", to: ""}], {
      //   copyUnmodified: true,
      //   ignore: ["css/**", "js/**", "src/**", "vendors/**"]
      // }),
      new Webpack.LoaderOptionsPlugin({minimize: true, debug: false}),
      new ManifestPlugin({fileName: "manifest.json"})
    ];

    return plugins
  },

  moduleOptions: function () {
    return {
      rules: [
        {
          test: /\.s[ac]ss$/,
          use: [
            MiniCssExtractPlugin.loader,
            {loader: "css-loader", options: {sourceMap: true}},
            {
              loader: "postcss-loader",
              options: {
                ident: 'postcss',
                plugins: () => [require('tailwindcss'), require("autoprefixer")],
                sourceMap: true
              }
            },
            {loader: "sass-loader", options: {sourceMap: true}}
          ]
        },
        {test: /\.tsx?$/, use: "ts-loader", exclude: /node_modules/},
        {test: /\.jsx?$/, loader: "babel-loader", exclude: /node_modules/},
        {test: /\.(woff|woff2|ttf|svg)(\?v=\d+\.\d+\.\d+)?$/, use: "url-loader"},
        {test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, use: "file-loader"},
        {test: require.resolve("jquery"), use: "expose-loader?jQuery!expose-loader?$"},
        {test: /\.go$/, use: "gopherjs-loader"}
      ]
    }
  },

  buildConfig: function () {
    // NOTE: If you are having issues with this not being set "properly", make
    // sure your GO_ENV is set properly as `buffalo build` overrides NODE_ENV
    // with whatever GO_ENV is set to or "development".
    const env = process.env.NODE_ENV || "development";

    var config = {
      mode: env,
      entry: configurator.entries(),
      output: {filename: "[name].[hash].js", path: `${__dirname}/public/assets`},
      plugins: configurator.plugins(),
      module: configurator.moduleOptions(),
      resolve: {
        extensions: ['.ts', '.js', '.jsx', ".tsx", '.json'],

        // import 相对路径的查找
        alias: {
          src: path.resolve(__dirname, 'assets/js/src'),
        },
      },
    };

    if (env === "development") {
      config.plugins.push(new LiveReloadPlugin({appendScriptTag: true}))
      return config
    }

    const uglifier = new UglifyJsPlugin({
      uglifyOptions: {
        beautify: false,
        mangle: {keep_fnames: true},
        output: {comments: false},
        compress: {}
      }
    })

    config.optimization = {
      minimizer: [uglifier]
    }

    return config
  }
};

module.exports = configurator.buildConfig()
