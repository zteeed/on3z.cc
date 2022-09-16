const form = document.getElementById('generate-short-url');
form.onsubmit = function(e) {
    const xhttp = new XMLHttpRequest();
    xhttp.open("POST", "/data/shorten", true);
    xhttp.setRequestHeader("Content-Type", "application/json");
    const params = {"longURL": document.getElementById("url").value};
    const data = JSON.stringify(params);
    xhttp.send(data);
    xhttp.onload = function() {
        if (this.status === 201) {
            const result = JSON.parse(this.responseText);
            Swal.fire({
                title: 'Success!',
                html: 'Here is your short link: <u><a href="' + result.shortURL + '">'  +result.shortURL + '</a></u>',
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
    e.preventDefault();
};
