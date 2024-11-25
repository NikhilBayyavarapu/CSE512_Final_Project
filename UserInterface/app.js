document.addEventListener('DOMContentLoaded', function () {
    const app = document.getElementById('app');
  
    // Mock user data (to be replaced with API response later)
    let userData = null;
  
    // Render the login form
    function renderLoginForm() {
      app.innerHTML = `
        <h1>Welcome to MyBank</h1>
        <form id="loginForm">
          <input type="text" id="userId" placeholder="User ID" required />
          <input type="email" id="email" placeholder="Email" required />
          <input type="password" id="password" placeholder="Password" required />
          <button type="submit">Login</button>
        </form>
        <p id="errorMessage" style="color: red; text-align: center; display: none;">Invalid login credentials. Please try again.</p>
      `;
  
      const loginForm = document.getElementById('loginForm');
      const errorMessage = document.getElementById('errorMessage');
  
      loginForm.addEventListener('submit', async function (e) {
        e.preventDefault();
  
        // Get user inputs
        const userId = document.getElementById('userId').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
  
        try {
          // Send POST request to /login
          const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ userId, email, password }),
          });
  
          if (response.ok) {
            // Parse the response JSON
            userData = await response.json();
            renderDashboard();
          } else {
            // Handle login failure
            errorMessage.style.display = 'block';
          }
        } catch (error) {
          console.error('Error during login:', error);
          errorMessage.style.display = 'block';
        }
      });
    }
  
    // Render the user dashboard
    function renderDashboard() {
      app.innerHTML = `
        <h1>Welcome, ${userData.name}</h1>
        <div class="dashboard">
          <p><strong>Email:</strong> ${userData.email}</p>
          <p><strong>User ID:</strong> ${userData.userId}</p>
          <p><strong>Account Number:</strong> ${userData.accountNumber}</p>
        </div>
        <div class="table-container">
          <h2>Transaction History</h2>
          <table>
            <thead>
              <tr>
                <th>Date</th>
                <th>Description</th>
                <th>Amount</th>
              </tr>
            </thead>
            <tbody>
              ${userData.transactions
                .map(
                  (txn) => `
                    <tr>
                      <td>${txn.date}</td>
                      <td>${txn.description}</td>
                      <td>${txn.amount}</td>
                    </tr>
                  `
                )
                .join('')}
            </tbody>
          </table>
        </div>
        <button id="logout">Logout</button>
      `;
  
      const logoutButton = document.getElementById('logout');
      logoutButton.addEventListener('click', function () {
        userData = null; // Clear user data
        renderLoginForm();
      });
    }
  
    // Initial render
    renderLoginForm();
  });
  