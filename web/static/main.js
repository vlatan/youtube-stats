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
                content.innerHTML = `<ul>
                    <li>Id: ${data.id}</li>
                    <li>Title: ${data.title}</li>
                    <li>Privacy Status: ${data.privacyStatus}</li>
                    <li>Age Restriced: ${data.ageRestriced}</li>
                    <li>Embeddable: ${data.embeddable}</li>
                    <li>Region Restricted: ${data.regionRestricted}</li>
                    <li>Default Language: ${data.defaultLanguage}</li>
                    <li>Live Broadcast: ${data.liveBroadcast}</li>
                    <li>Duration: ${data.duration}</li>
                </ul>`;
            });
        });
    });
});