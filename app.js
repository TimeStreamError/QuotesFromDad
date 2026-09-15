document.getElementById('login-form').addEventListener('submit', async (event) => {
  event.preventDefault();
  await login();
});

// prepare for the edit quotes section
console.log("Setting up quote event listener");
document.getElementById("allquotes-section").addEventListener("click", function(event) {
    const button = event.target.closest(".edit-quote");
    console.log("Quote button clicked");
    if (!button) {
        return;
    }
    console.log(button.dataset.id)
    showEditQuote(button.dataset.id);
    });

// the JST token is 'stored' here for the session
var authToken;

async function submitQuote(event) {

    event.preventDefault();

    const response = await fetch("/admin/quotes", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            quotation: document.getElementById("quote").value,
            author: document.getElementById("author").value
        })
    });

    const html = await response.text();
    document.documentElement.innerHTML = html;
}

async function login() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  try {
    const res = await fetch('/admin/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to login: ${data.error}`);
    }

    if (data.token) {
      authToken = data.token
      document.getElementById('login-section').innerHTML = `Success!`;
    } else {
      alert('Login failed. Please check your credentials.');
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

async function signup() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  const authorization = document.getElementById('authorization').value

  try {
    const res = await fetch('/admin/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, authorization }),
    });
    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to create user: ${data.error}`);
    }
    console.log('User created!');
    await login();
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

async function allQuotes() {

  // Show the DIV that displays all the quotes
  showSection("allquotes-section");
  console.log("Starting to get all quotes")
  try { // in case something goes very wrong, beyond stuff handled via status

    //fetch all the quotes via the REST API
    const res = await fetch('/api/quotes', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`,
      },
    });

    // the JST token isn't valid. We need to check the refresh token which is stored in a cookie
    if (res.status == 401) {
        // ask for a new JWT token via refresh token
        console.log("Asking for new JWT token");
        const refreshRes = await fetch('/api/refresh');

        if (!refreshRes.ok) {
          showLogin();
          return;
        }

        const data = await refreshRes.json();
        authToken = data.token
        console.log("Trying fetch again")
        // now we can try the original request again
        const res = await fetch('/api/quotes', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${authToken}`,
          },
        });
    }

    // so initial fetch did not return a 401 Unauthorized status...
    if (!res.ok) {
      // ... but something else went wrong
      const data = await res.json();
      throw new Error(`Failed to retrieve all quotes: ${data.error}`); 
    }
    // ... and we got a 200 OK status.
    const data = await res.text();
    document.getElementById("allquotes-section").innerHTML = data;

  } catch (error) { // something in the http code went bad (not handled status stuff)
    alert:(`Error: ${error.message}`)
  }

}

function showHome() {
  showSection("home-section");
}

function showLogin() {
  showSection("login-section");
}

function showSection(id) {
    document.querySelectorAll("main > div").forEach(section => {
        section.classList.add("hidden");
    });

    document.getElementById(id).classList.remove("hidden");
}

async function showEditQuote(id) {

    const response = await fetch(`/api/quotes/${id}`);
    const quote = await response.json();

    // Put quote information into your edit form
    document.getElementById("quote-text").value = quote.Quote;
    let author = ""
    if (quote.Author.Valid) {
      author = quote.Author.String;
    }
    document.getElementById("quote-author").value = author;
    document.getElementById("quote-id").value = id;

    // Switch sections
    showSection("edit-quote-section");
}

async function putEditedQuote(event) {
  event.preventDefault();

    quoteID = document.getElementById("quote-id").value;
    const response = await fetch(`/api/quotes/${quoteID}`, {
      method: "PUT",
      headers: {
          "Content-Type": "application/json"
      },
      body: JSON.stringify({
          quotation: document.getElementById("quote-text").value,
          author: document.getElementById("quote-author").value,
          id: quoteID,
      })
    });
    const responseJson = await response.json();
    document.getElementById('edit-quote-section').innerHTML = `Success!`;
    console.log(responseJson)
}