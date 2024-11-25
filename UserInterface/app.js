document.addEventListener('DOMContentLoaded', function () {
  const app = document.getElementById('app');

  // Mock user data (to be replaced with API response later)
  let userData = null;

  // Static transaction data as fallback
  const staticTransactions = [
    { date: '2024-01-01', description: 'Deposit', amount: '$500' },
    { date: '2024-01-02', description: 'Withdrawal', amount: '$100' },
  ];

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

      // If password is "admin", directly render the dashboard
      // if (password === 'admin') {
      //   userData = {
      //     name: 'Admin User',
      //     email,
      //     userId,
      //     accountNumber: '1234567890',
      //   };
      //   renderDashboard();
      //   return;
      // }

      try {
        // Send POST request to /login
        const response = await fetch('http://localhost:8080/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body : JSON.stringify({
            user_id : userId,
            email,
            password
          })
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
  async function renderDashboard() {
    app.innerHTML = `
      <h1>Welcome, ${userData.name}</h1>
      <div class="dashboard">
        <p><strong>Email:</strong> ${userData.email}</p>
        <p><strong>User ID:</strong> ${userData.userId}</p>
        <p><strong>Account Number:</strong> ${userData.accountNumber}</p>
      </div>
      <div class="table-container">
        <h2>Transaction History</h2>
        <p>Loading transactions...</p>
      </div>
      <button id="logout">Logout</button>
    `;

    const logoutButton = document.getElementById('logout');
    logoutButton.addEventListener('click', function () {
      userData = null; // Clear user data
      renderLoginForm();
    });

    // Fetch transactions
    try {
      const transactionsResponse = await fetch(
        `http://localhost:8080/transactions?email=${encodeURIComponent(userData.email)}&userId=${encodeURIComponent(userData.userId)}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      let transactions = staticTransactions; // Default to static data

      if (transactionsResponse.ok) {
        // Parse the transactions if the response is OK
        transactions = await transactionsResponse.json();
      }

      // Render transactions
      renderTransactions(transactions);
    } catch (error) {
      console.error('Error fetching transactions:', error);
      renderTransactions(staticTransactions);
    }
  }

  // Render transaction history
  function renderTransactions(transactions) {
    const tableContainer = document.querySelector('.table-container');
    tableContainer.innerHTML = `
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
          ${transactions
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
    `;
  }

  // Initial render
  renderLoginForm();
});
