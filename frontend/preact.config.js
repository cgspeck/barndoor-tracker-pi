export default (config, env, helpers) => {
  if (env.isProd) {
    config.output.filename = (chunkData) => {
      return `${chunkData.chunk.name.slice(0, 6)}.js`
    };
    config.output.chunkFilename = `[id].c.js`;

    let cssPlugin = helpers.getPluginsByName(config, "MiniCssExtractPlugin")[0].plugin;
    cssPlugin.options.filename = `[contenthash:5].css`;
    cssPlugin.options.chunkFilename = `[contenthash:5].c.css`;
  }
};

