document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('myForm');

    form.addEventListener('submit', handleSubmit);
});

function handleSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    fetch('/your-server-endpoint', {
        method: 'POST',
        body: formData
    })
        .then(handleResponse)
        .catch(handleError);
}

function handleResponse(response) {
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
}

function handleError(error) {
    console.error('Error:', error);
    // Handle errors, such as displaying an error message to the user
}