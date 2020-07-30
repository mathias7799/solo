import $ from 'jquery';

function getCurrencyDetails(chainId) {
    var coingeckoId;
    var ticker;
    if (chainId == 1) {
        coingeckoId = "ethereum"
        ticker = "ETH"
    } else if (chainId == 61) {
        coingeckoId = "ethereum-classic"
        ticker = "ETC"
    }
    return { coingeckoId, ticker }
}

function getCurrencyPrice(coingeckoId, callback, callbackError) {
    $.get("https://api.coingecko.com/api/v3/simple/price", { vs_currencies: "usd", ids: coingeckoId }, callback).fail(callbackError);
}

export {
    getCurrencyDetails,
    getCurrencyPrice
}
