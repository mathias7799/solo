module.exports = {
    plugins =[
        new webpack.ProvidePlugin({
            $: 'jquery',
            'window.$': 'jquery',
            jquery: 'jquery',
            'window.jQuery': 'jquery',
            jQuery: 'jquery'
        })
    ]
}