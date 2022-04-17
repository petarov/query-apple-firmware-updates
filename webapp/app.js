(async () => {
    const DEVICE_TEMPLATE = document.getElementById('device-item-template').innerHTML;
    const RESULTS = document.getElementById('results');
    const TERM = document.getElementById('term');

    const doSearch = () => {
        RESULTS.innerHTML = 'Searching ...Please wait';

        const term = document.getElementById('term').value;
        const display = (results) => {
            if (results.length > 0) {
                RESULTS.innerHTML = '';
                for (const device of results) {
                    RESULTS.insertAdjacentHTML("beforeend", getDeviceHtml(DEVICE_TEMPLATE, device));
                }
            } else {
                RESULTS.innerHTML = 'No devices found';
            }
        };

        const search = (term, cb) => {
            fetch(`/api/devices/search?key=${term}`).then(response => {
                if (response.status / 100 === 2) {
                    response.json().then(data => cb(null, data));
                } else {
                    cb(`HTTP Err: ${response.status} ${response.statusText}`, null);
                }
            }).catch(err => cb(err, null));
        };

        search(term, (err, json) => {
            if (err) {
                RESULTS.innerHTML = err;
                if (json != null) {
                    RESULTS.innerHTML += "<br>";
                    RESULTS.innerHTML += JSON.parse(json).errorMessage;
                }
            } else {
                display(json);
            }
        });
    };

    document.getElementById('search').onclick = () => doSearch();
    TERM.onkeyup = async (event) => event.key === 'Enter' && doSearch();
    document.getElementById('clear').onclick = () => {
        RESULTS.innerHTML = 'Search results will appear here';
        TERM.value = '';
        TERM.focus();
    };
})();

function getDeviceHtml(template, device) {
    var tpl = template.slice(0);
    tpl = tpl.replace(/__PRODUCT__/g, device.product);
    tpl = tpl.replace(/__NAME__/g, device.name);
    tpl = tpl.replace(/__LAST_CHECK__/g, new Date(device.last_checked_on).toLocaleString());
    var iconName = 'fab fa-apple';
    // if (device.product.indexOf('iPhone') != -1) {
    //     iconName = 'fas fa-mobile-alt';
    if (device.product.indexOf('iPad') != -1) {
        iconName = 'fas fa-tablet';
    } else if (device.product.indexOf('iMac') != -1) {
        iconName = 'fas fa-desktop';
    } else if (device.product.indexOf('MacBook') != -1) {
        iconName = 'fas fa-laptop';
    // } else if (device.product.indexOf('iPod') != -1) {
    //     iconName = 'fas fa-mp3-player';
    }
    tpl = tpl.replace(/__ICON_NAME__/g, iconName);
    return tpl;
}
