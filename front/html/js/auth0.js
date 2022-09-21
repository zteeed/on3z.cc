let auth0 = null

window.onload = async () => {
    await configureClient()
    await processLoginState()
    updateUI()
}

const configureClient = async () => {
    auth0 = await createAuth0Client({
        domain: "llilgq.eu.auth0.com",
        client_id: "32DMIWDgPTwointJxN8exbcJtj1Wpqeo",
        audience: "https://llilgq.eu.auth0.com/api/v2/",
        cacheLocation: "localstorage",
        leeway: 300
    })
}

const processLoginState = async () => {
    // Check code and state parameters
    const query = window.location.search
    if (query.includes("code=") && query.includes("state=")) {
        // Process the login state
        await auth0.handleRedirectCallback()
        // Use replaceState to redirect the user away and remove the querystring parameters
        window.history.replaceState({}, document.title, window.location.pathname)
    }
}

function loadTable(token) {
    $('#dataTable').DataTable({
        "lengthMenu": [10, 25, 50, 100],
        "pageLength": 10,
        "ajax": {
            "url": "/data/shorten",
            "type": "GET",
            "dataSrc": "",
            "data": {
                "length": 1000,
                "offset": 0,
            },
            "beforeSend": function (request) {
                request.setRequestHeader("Authorization", token);
            }
        },
        "rowCallback": function (row, data) {
            $('td:eq(0)', row).html('<a href="' + window.location.origin + '/' + data.ShortURL + '">' + data.ShortURL + '</a>');
            $('td:eq(1)', row).html('<a href="' + data.LongURL + '">' + data.LongURL + '</a>');
            $('td:eq(2)', row).html(`
                <div class=\"field\"><button type=\"submit\" onclick='editURL(this, "` + token + `", "` + data.ShortURL + `", "` + data.LongURL + `")'><i class=\"fa-solid fa-pencil\"></i></button></div>
            `);
        },
        "columns": [
            {data: "ShortURL"},
            {data: "LongURL"},
            {"defaultContent": "<div class=\"field\"><button type=\"submit\"><i class=\"fa-solid fa-pencil\"></i></button></div>"},
        ],
    });
}

const updateUI = async () => {
    if (localStorage.getItem('notification') === 'loggedout') {
        localStorage.removeItem('notification')
        Swal.fire({
            title: 'Success!',
            html: 'You\'re successfully logged out',
            icon: 'success',
            confirmButtonText: 'OK'
        })
    }
    if (localStorage.getItem('notification') === 'loggedin') {
        localStorage.removeItem('notification')
        Swal.fire({
            title: 'Success!',
            html: 'You\'re successfully logged in',
            icon: 'success',
            confirmButtonText: 'OK'
        })
    }
    const isAuthenticated = await auth0.isAuthenticated()
    document.getElementById("btn-login").style.display = isAuthenticated ? "none" : ""
    document.getElementById("btn-logout").style.display = !isAuthenticated ? "none" : ""
    if (isAuthenticated) {
        let user = JSON.stringify(
            await auth0.getUser()
        )
        let user_obj = JSON.parse(user)
        document.getElementById("user").innerHTML = user_obj.name;
        let token = await auth0.getTokenSilently();
        $(document).ready(function () {
            loadTable(token)
        });
    }
    document.getElementById("anonymous-form").style.display = isAuthenticated ? "none" : ""
    document.getElementById("authenticated-form").style.display = !isAuthenticated ? "none" : ""
}

const login = async () => {
    localStorage.setItem('notification', 'loggedin')
    await auth0.loginWithRedirect({
        redirect_uri: window.location.href,
    })
}

const logout = () => {
    localStorage.setItem('notification', 'loggedout')
    auth0.logout({
        returnTo: window.location.href,
    })
}
