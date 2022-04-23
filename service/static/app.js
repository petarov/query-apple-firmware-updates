(async () => {
    const DEVICE_TEMPLATE = document.getElementById('device-item-template').innerHTML;
    const UPDATE_TEMPLATE = document.getElementById('update-item-template').innerHTML;
    const RESULTS = document.getElementById('results');
    const TERM = document.getElementById('term');

    const doSearch = () => {
        RESULTS.innerHTML = 'Searching ...Please wait';

        const term = document.getElementById('term').value;
        const display = (results) => {
            if (results.length > 0) {
                RESULTS.innerHTML = '';
                for (const device of results) {
                    RESULTS.insertAdjacentHTML('beforeend', getDeviceHtml(DEVICE_TEMPLATE, device));
                }

                const updateButtons = document.querySelectorAll('.update-btn');
                for (const btn of updateButtons) {
                    btn.addEventListener('click', (event) => {
                        const product = event.target.dataset.product;
                        getUpdates(event.target.dataset.product, (err, json) => {
                            if (err) {
                                document.getElementById("updates" + product).innerHTML = err;
                            } else {
                                var container = document.getElementById("updates-body-" + product);
                                if (json.updates) {
                                    for (let i = 0; i < json.updates.length; i++) {
                                        container.insertAdjacentHTML('beforeend', getUpdateHtml(UPDATE_TEMPLATE, i + 1, json.updates[i]));
                                    }
                                }
                                container.parentElement.classList.remove('is-hidden');
                            }
                        });
                        event.target.remove();
                    });
                }
            } else {
                RESULTS.innerHTML = 'No devices found';
            }
        };

        const getUpdates = (product, cb) => {
            fetch(`api/updates/${product}`).then(response => {
                if (response.status / 100 === 2) {
                    response.json().then(data => cb(null, data));
                } else {
                    cb(`HTTP Err: ${response.status} ${response.statusText}`, null);
                }
            }).catch(err => cb(err, null));
        };

        const search = (term, cb) => {
            fetch(`api/devices/search?key=${term}`).then(response => {
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

function getUpdateHtml(template, row, update) {
    var tpl = template.slice(0);
    tpl = tpl.replace(/__ROW__/g, row);
    tpl = tpl.replace(/__VERSION__/g, update.version);
    tpl = tpl.replace(/__BUILD__/g, update.build_id);
    tpl = tpl.replace(/__RELEASE_DATE__/g, new Date(update.released_on).toLocaleString());
    tpl = tpl.replace(/__FILESIZE__/g, Math.round((update.attributes.size/(1024*1024*1024)) * 100) / 100 + "GB");
    tpl = tpl.replace(/__FILENAME__/g, update.attributes.filename);
    tpl = tpl.replace(/__LINK__/g, update.attributes.url);
    return tpl;
}