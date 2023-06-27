let auth0 = null

async function configureClient() {
    auth0 = await createAuth0Client({
        domain: "on3zcc.eu.auth0.com",
        client_id: "BWgOvVNOQglGJDK10XFNjjp7aZGmNvo4",
        audience: "https://on3zcc.eu.auth0.com/api/v2/",
        cacheLocation: "localstorage",
        leeway: 300
    })
}

async function processLoginState() {
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
  $.ajax({
    url: "/data/shorten",
    type: "GET",
    data: {
        "length": 1000,
        "offset": 0,
    },
    beforeSend: function (request) {
      request.setRequestHeader("Authorization", token);
    },
    success: function (response) {
      // Sort the response data in descending order
      response.reverse();

      // Initialize the DataTable with the sorted data
      dataTable = $('#dataTable').DataTable({
        "lengthMenu": [10, 25, 50, 100],
        "pageLength": 10,
        "data": response,
        "order": [], // Remove default sorting
        "rowCallback": function (row, data) {
          $('td:eq(0)', row).html('<a href="' + window.location.origin + '/' + data.ShortURL + '">' + data.ShortURL + '</a>');
          $('td:eq(1)', row).html('<a href="' + data.LongURL + '">' + data.LongURL + '</a>');
          $('td:eq(2)', row).html(`
            <div class="field"><button type="submit" onclick='editURL(this, "${token}", "${data.ShortURL}", "${data.LongURL}")'><i class="fa-solid fa-pencil"></i></button></div>
          `);
        },
        "columns": [
          { data: "ShortURL" },
          { data: "LongURL" },
          { "defaultContent": "<div class=\"field\"><button type=\"submit\"><i class=\"fa-solid fa-pencil\"></i></button></div>" },
        ],
      });
    },
    error: function (error) {
      console.log("Error loading table:", error);
    }
  });
}

function toggleElementVisibility(elementId, isVisible) {
  const element = document.getElementById(elementId);
  if (isVisible) {
    element.classList.add('transition');
    element.classList.remove('hidden');
    element.style.display = ''; // Reset the display property
  } else {
    element.classList.remove('transition');
    element.classList.add('hidden');
    element.addEventListener('transitionend', function() {
      element.style.display = 'none'; // Set display: none after the transition is complete
    }, { once: true });
  }
}

async function updateUI() {
    if (localStorage.getItem('notification') === 'loggedout') {
        localStorage.removeItem('notification')
        Swal.fire({
            title: 'Success!',
            html: 'You\'re successfully logged out',
            icon: 'success',
            confirmButtonText: 'OK'
        })
    }
    const user = await auth0.getUser();
    if (localStorage.getItem('notification') === 'loggedin' && user) {
        localStorage.removeItem('notification')
        Swal.fire({
            title: 'Success!',
            html: 'You\'re successfully logged in',
            icon: 'success',
            confirmButtonText: 'OK'
        })
    }
    const isAuthenticated = await auth0.isAuthenticated()
    toggleElementVisibility('btn-login', !isAuthenticated);
    toggleElementVisibility('btn-logout', isAuthenticated);
    if (isAuthenticated) {
        let user = JSON.stringify(
            await auth0.getUser()
        )
        let user_obj = JSON.parse(user)
        document.getElementById("user").innerHTML = user_obj.name;
        let token = await auth0.getTokenSilently();
        $(document).ready(function () {
            loadTable(token);
        });
    }
    toggleElementVisibility('anonymous-form', !isAuthenticated);
    toggleElementVisibility('authenticated-form', isAuthenticated);
}

async function login() {
  localStorage.setItem('notification', 'loggedin')
  await auth0.loginWithRedirect({
      redirect_uri: window.location.href,
  })
}

async function logout() {
    localStorage.setItem('notification', 'loggedout')
    auth0.logout({
        returnTo: window.location.href,
    })
}

async function init() {
    await configureClient()
    try {
      await processLoginState()
    } catch (error) {
      console.log(error)
    }
    await updateUI()
}

window.addEventListener('load', init);
