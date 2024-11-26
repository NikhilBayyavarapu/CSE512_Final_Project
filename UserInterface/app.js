document.addEventListener('DOMContentLoaded', function () {
  const app = document.getElementById('app');

  let userData = null;

  const staticTransactions = [
    { date: '2024-01-01', description: 'Deposit', amount: '$500' },
    { date: '2024-01-02', description: 'Withdrawal', amount: '$100' },
  ];

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

      const userId = document.getElementById('userId').value;
      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      try {
        const response = await fetch('http://localhost:8080/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ user_id: userId, email, password }),
        });

        if (response.ok) {
          userData = (await response.json()).data;
          userData.user_id = userId ;
          renderDashboard();
        } else {
          errorMessage.style.display = 'block';
        }
      } catch (error) {
        console.error('Error during login:', error);
        errorMessage.style.display = 'block';
      }
    });
  }

  async function renderDashboard() {
    app.innerHTML = `
      <h1>Welcome, ${userData.name}</h1>
      <div class="dashboard">
        <table class="user-table">
          <tr>
            <th>Email:</th>
            <td>${userData.email}</td>
          </tr>
          <tr>
            <th>User ID:</th>
            <td>${userData.user_id}</td>
          </tr>
          <tr>
            <th>Account Number:</th>
            <td>${userData.accountNumber}</td>
          </tr>
          <tr>
            <th>Current Balance:</th>
            <td>$${userData.balance}</td>
          </tr>
        </table>
      </div>
      <button id="sendMoneyButton">Send Money</button>
      <div id="sendMoneyFormContainer"></div>
      <div class="table-container">
        <h2>Transaction History</h2>
        <p>Loading transactions...</p>
      </div>
      <button id="logout">Logout</button>
    `;

    document.getElementById('logout').addEventListener('click', function () {
      userData = null;
      renderLoginForm();
    });

    document.getElementById('sendMoneyButton').addEventListener('click', openSendMoneyForm);

    try {
      const transactionsResponse = await fetch(
        `http://localhost:8080/transactions?user_id=${userData.user_id}&email=${userData.email}`,
        { method: 'GET', headers: { 'Content-Type': 'application/json' } }
      );

      let transactions = staticTransactions;
      if (transactionsResponse.ok) {
        transactions = await transactionsResponse.json();
      }

      renderTransactions(transactions);
    } catch (error) {
      console.error('Error fetching transactions:', error);
      renderTransactions(staticTransactions);
    }
  }

  function openSendMoneyForm() {
    const sendMoneyButton = document.getElementById('sendMoneyButton');
    sendMoneyButton.disabled = true;
  
    const formContainer = document.getElementById('sendMoneyFormContainer');
    formContainer.innerHTML = `
      <table class="send-money-table" style="width: 100%; border-collapse: collapse; margin-top: 20px;">
        <tbody>
          <tr>
            <td><label for="receiverId">Receiver ID:</label></td>
            <td><input type="text" id="receiverId" placeholder="Receiver ID" required style="width: 100%; padding: 5px;" /></td>
          </tr>
          <tr>
            <td><label for="receiverEmail">Receiver Email:</label></td>
            <td><input type="email" id="receiverEmail" placeholder="Receiver Email" required style="width: 100%; padding: 5px;" /></td>
          </tr>
          <tr>
            <td><label for="receiverAccount">Receiver Account Number:</label></td>
            <td><input type="text" id="receiverAccount" placeholder="Receiver Account Number" required style="width: 100%; padding: 5px;" /></td>
          </tr>
          <tr>
            <td><label for="amount">Amount to Send:</label></td>
            <td><input type="number" id="amount" placeholder="Amount to Send" required min="1" style="width: 100%; padding: 5px;" /></td>
          </tr>
          <tr>
            <td><label for="confirm">Confirm Transaction:</label></td>
            <td>
              <input type="checkbox" id="confirm" name="confirm" style="margin-right: 10px;" />
              <label for="confirm">I confirm the transaction</label>
            </td>
          </tr>
          <tr>
            <td colspan="2" style="text-align: right; padding-top: 10px;">
              <button type="submit" id="sendNowButton" disabled style="margin-right: 10px; padding: 10px 15px;">Send Now</button>
              <button type="button" id="cancelButton" style="padding: 10px 15px;">Cancel</button>
            </td>
          </tr>
        </tbody>
      </table>
      <p id="errorMessage" style="color: red; display: none; margin-top: 10px;">Insufficient balance!</p>
    `;
  
    const confirmCheckbox = document.getElementById('confirm');
    const sendNowButton = document.getElementById('sendNowButton');
    const cancelButton = document.getElementById('cancelButton');
  
    confirmCheckbox.addEventListener('change', function () {
      sendNowButton.disabled = !this.checked;
    });
  
    const sendMoneyForm = formContainer.querySelector('.send-money-table');
    sendNowButton.addEventListener('click', function (e) {
      e.preventDefault();
  
      const amount = parseFloat(document.getElementById('amount').value);
      if (amount > userData.balance) {
        document.getElementById('errorMessage').style.display = 'block';
      } else {
        const payload = {
          receiverId: document.getElementById('receiverId').value,
          receiverEmail: document.getElementById('receiverEmail').value,
          receiverAccount: document.getElementById('receiverAccount').value,
          amount,
        };
  
        console.log('Simulated API Call Payload:', payload);
        alert(`$${amount} sent successfully!`);
        closeSendMoneyForm();
      }
    });
  
    cancelButton.addEventListener('click', closeSendMoneyForm);
  
    function closeSendMoneyForm() {
      formContainer.innerHTML = '';
      sendMoneyButton.disabled = false;
    }
  }
  
  
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

  renderLoginForm();
});
