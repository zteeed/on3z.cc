function add_short_url(url, dom_url_id, token = "") {
    const xhttp = new XMLHttpRequest();
    xhttp.open("POST", "/data/shorten", true);
    xhttp.setRequestHeader("Content-Type", "application/json");
    if (token !== "") {
        xhttp.setRequestHeader("Authorization", token);
    }
    const params = {"longURL": document.getElementById(dom_url_id).value};
    const data = JSON.stringify(params);
    xhttp.send(data);
    xhttp.onload = function () {
        if (this.status === 201) {
            const result = JSON.parse(this.responseText);
            Swal.fire({
                title: 'Success!',
                html: 'Here is your short link: <u><a href="' + result.shortURL + '">' + result.shortURL + '</a></u>',
                icon: 'success',
                confirmButtonText: 'OK'
            })
        } else {
            Swal.fire({
                title: 'Error!',
                text: 'Cannot generate short url for now. Try again later',
                icon: 'error',
                confirmButtonText: 'OK'
            })
        }
    }
}

const form = document.getElementById('generate-short-url');
form.onsubmit = function (e) {
    add_short_url("/data/shorten", "url");
    e.preventDefault();
};
const form_authenticated = document.getElementById('generate-short-url-authenticated');
form_authenticated.onsubmit = async function (e) {
    let token = await auth0.getTokenSilently();
    add_short_url("/data/shorten", "url-authenticated", token);
    e.preventDefault();
};

function editURL(e, token, shortURL, longURL) {
    Swal.fire({
        title: 'Edit the URL destination',
        html:
            'You are going to update the following short url: <b><u>' + window.location.origin + '/' + shortURL + '</u></b><br><br>' +
            'Previous longURL for this short link is: <b><u>' + longURL + '</u></b><br><br>' +
            'Enter the new longURL to associate:',
        input: 'text',
        inputAttributes: {
            autocapitalize: 'off'
        },
        showCancelButton: true,
        showLoaderOnConfirm: true,
        confirmButtonText: 'Yes, update it!',
        cancelButtonText: 'No, cancel!',
        confirmButtonColor: '#5468D4',
        cancelButtonColor: '#dc3741',
        preConfirm: (newURL) => {
            return fetch(`/data/shorten`, {
                method: 'PUT',
                headers: {
                    'Authorization': token
                },
                body: JSON.stringify({
                    longURL: newURL,
                    shortURL: shortURL,
                })
            }).then(response => {
                if (!response.ok) {
                    throw new Error(response.statusText)
                }
                return response.text()
            }).catch(error => {
                Swal.showValidationMessage(
                    `Request failed: ${error}`
                )
            })
        },
        allowOutsideClick: () => !Swal.isLoading()
        }).then((result) => {
        if (result.isConfirmed) {
            Swal.fire({
                title: 'Success!',
                html: 'Short link destination successfully updated for <b><u>' + window.location.origin + '/' + shortURL + '</u></b>',
                icon: 'success',
                confirmButtonText: 'OK', }).then((result) => {location.reload()
            })
            // TODO: Find a way to just reload table ajax after confirm
        }
    })
}
