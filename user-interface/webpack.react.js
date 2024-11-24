const path = require("path");
const webpack = require('webpack');
const UI = require(__dirname + "/../settings/ui-settings.js");

module.exports = function(_env, argv) {
  const isProduction = argv.mode === "production";
  const isDevelopment = !isProduction;

  return {
    devtool: isDevelopment && "cheap-module-source-map",
    entry:[ "./react/index.js" , "./sass/app.scss" , "./sass/dark.scss"],
    output: {
      path: path.resolve(__dirname, "../public"),
      filename: "js/app-unlocked.js",
      publicPath: "../"
    },
    plugins: [
      new webpack.EnvironmentPlugin({
        MIX_IMAGE_DIMENSIONS_W: UI.dimensions_w,
        MIX_IMAGE_DIMENSIONS_H: UI.dimensions_h,
        MIX_IMAGE_DIMENSIONS_SMALL_W: UI.dimensions_small_w,
        MIX_IMAGE_DIMENSIONS_SMALL_H: UI.dimensions_small_h,
        MIX_UPLOAD_RULES: UI.rules,
        MIX_UPLOAD_SMALL_RULES: UI.rules_small,
        MIX_APP_URL: UI.host_addr,
        MIX_APP_HOSTNAME:UI.host_name,
        MIX_VERSION_NO: UI.version_no,
        MIX_EXTRA_INFO: UI.extra_info,
        MIX_FREE_MODE: UI.free_mode,
        MIX_BOARDS: UI.boards,
      })
    ],
    module: {
      rules: [
        {
          test: /\.jsx?$/,
          exclude: /node_modules/,
          use: {
            loader: "babel-loader",
            options: {
              cacheDirectory: true,
              cacheCompression: false,
            }
          }
        }, {
          test: /\.scss$/,
          exclude: /node_modules/,
          use: [
              {
                  loader: 'file-loader',
                  options: { outputPath: '../public/css/', name: '[name]-unlocked.css'}
              },
              'sass-loader'
          ]
        }
      ]
    },
    resolve: {
      extensions: [".js", ".jsx"]
    }
  };
};
