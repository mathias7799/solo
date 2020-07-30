function getSi(num) {
    if (isNaN(num)) {
        return [1, ""]
    }
    const DIVIDER = 1000;
    const SI_SYMBOLS = "kMGTPEZY";
    if (num <= 1000) {
        return [1, ""];
    }
    for (var i = 0; i < SI_SYMBOLS.length; i++) {
        var divider_new = 1;
        for (var a = 0; a <= i; a++) {
            divider_new = divider_new * DIVIDER;
        }

        var divided_sum = num / divider_new;

        if (divided_sum < 1000) {
            return [divider_new, SI_SYMBOLS[i]];
        }
    }
    return [Math.pow(10, 24), "Y"];
}

export {
    getSi
}