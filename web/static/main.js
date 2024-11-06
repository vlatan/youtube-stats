document.addEventListener('DOMContentLoaded', () => {
    const content = document.getElementById('content')
    const form = document.getElementById('myForm');

    form.addEventListener('submit', (event) => {
        event.preventDefault();
        const formData = new FormData(event.target);
        fetch('/', {
            method: 'POST',
            body: formData
        }).then(response => {
            if (response.ok) {
                response.json().then(data => {
                    content.innerHTML = data.title;
                    return
                });
            }
            content.innerHTML = "Sorry no content";
        });
    });
});