const path = require("path");

module.exports = function(_env, argv) {
  const isProduction = argv.mode === "production";
  const isDevelopment = !isProduction;

  return {
    devtool: isDevelopment && "cheap-module-source-map",
    entry:[ "./react/index.js" , "./sass/app.scss"],
    output: {
      path: path.resolve(__dirname, "../public"),
      filename: "js/app.js",
      publicPath: "../"
    },
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
                  options: { outputPath: '../public/css/', name: 'app.css'}
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