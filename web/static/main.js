document.addEventListener('DOMContentLoaded', () => {
    const content = document.getElementById('content')
    const form = document.getElementById('myForm');

    form.addEventListener('submit', event => {
        event.preventDefault();
        const formData = new FormData(event.target);
        fetch('/', {
            method: 'POST',
            body: formData
        }).then(response => {
            if (!response.ok) {
                content.innerHTML = "Not been able to fetch the info for this video.";
                return;
            }
            response.json().then(data => {
                const regionRestricted = data.regionRestricted.length ? data.regionRestricted : false;
                const defaultLanguage = data.defaultLanguage || "none";
                console.log(data);

                for (const [key, value] of Object.entries(data)) {
                    document.querySelector(`[data-id=${key}]`).innerText = value;
                }

                var countryElements = document.getElementById('countries').childNodes;
                var countryCount = countryElements.length;
                for (var i = 0; i < countryCount; i++) {
                    countryElements[i].onclick = () => {
                        alert('You clicked on ' + this.getAttribute('data-name'));
                    }
                }
            });
        });
    });
});