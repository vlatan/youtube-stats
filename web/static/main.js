document.addEventListener('DOMContentLoaded', () => {
    // const content = document.getElementById('content')
    const form = document.getElementById('myForm');
    const badVideo = document.getElementById('badVideo');
    const infoCells = document.querySelectorAll('td[data-id]');

    var countryElements = document.getElementById('countries').children;
    for (let countryElem of countryElements) {
        let countryName = countryElem.getAttribute('data-name');
        countryElem.innerHTML = `<title>${countryName}</title>`;
    }

    form.addEventListener('submit', event => {

        // reset the previous data, styles and events after submit
        badVideo.innerText = "";
        for (let infoCell of infoCells) {
            infoCell.innerText = "";
        }
        for (let countryElem of countryElements) {
            countryElem.removeAttribute("style");
            countryElem.onmouseenter = () => { };
            countryElem.onmouseleave = () => { };
        }

        event.preventDefault();
        const formData = new FormData(event.target);
        fetch('/', {
            method: 'POST',
            body: formData
        }).then(response => {
            if (!response.ok) {
                badVideo.innerText = "Not been able to fetch the info for this video.";
                return;
            }

            response.json().then(data => {
                // const regionRestricted = data.regionRestricted.length ? data.regionRestricted : false;
                // const defaultLanguage = data.defaultLanguage || "none";

                for (const [key, value] of Object.entries(data)) {
                    document.querySelector(`td[data-id=${key}]`).innerText = value;
                }

                for (let countryCode of data.regionRestricted) {
                    let country = document.querySelector(`path[data-id=${countryCode}]`);
                    country.style.fill = "red";
                    country.onmouseenter = event => {
                        event.target.style.fill = "crimson";
                    }
                    country.onmouseleave = event => {
                        event.target.style.fill = "red";
                    }
                }
            });
        });
    });
});