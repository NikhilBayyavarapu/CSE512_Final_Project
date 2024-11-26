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
          userData.user_id = userId;
          userData.password = password ;
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
    const formatter = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    });
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
            <td>${formatter.format(userData.balance)}</td>
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
      console.log(userData) ;
      const transactionsResponse = await fetch(
        `http://localhost:8080/transactions?sender_id=${userData.user_id}`,
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
      <table>
        <tr>
          <td><label for="receiverName">Receiver Name:</label></td>
          <td><input type="text" id="receiverName" placeholder="Receiver's Name" required /></td>
          <td><span id="receiverNameError" style="color: red; display: none;">This field is required</span></td>
        </tr>
        <tr>
          <td><label for="receiverId">Receiver ID:</label></td>
          <td><input type="text" id="receiverId" placeholder="Receiver's User ID" required /></td>
          <td><span id="receiverIdError" style="color: red; display: none;">This field is required and must be a positive number</span></td>
        </tr>
        <tr>
          <td><label for="receiverEmail">Receiver Email:</label></td>
          <td><input type="email" id="receiverEmail" placeholder="Receiver's Email" required /></td>
          <td><span id="receiverEmailError" style="color: red; display: none;">Please enter a valid email</span></td>
        </tr>
        <tr>
          <td><label for="receiverAccount">Receiver Account Number:</label></td>
          <td><input type="text" id="receiverAccount" placeholder="Receiver's Account Number" required /></td>
          <td><span id="receiverAccountError" style="color: red; display: none;">This field is required and must be a positive number</span></td>
        </tr>
        <tr>
          <td><label for="amount">Amount to Send:</label></td>
          <td><input type="number" id="amount" placeholder="Amount" required min="0.01" /></td>
          <td><span id="amountError" style="color: red; display: none;">Amount must be greater than 0 and less than or equal to your balance</span></td>
        </tr>
        <tr>
          <td><input type="checkbox" id="confirm" /> I confirm this transaction</td>
          <td></td>
          <td><span id="confirmError" style="color: red; display: none;">You must confirm the transaction</span></td>
        </tr>
      </table>
      <div style="margin-top: 10px;">
        <button id="sendNowButton" disabled>Send Now</button>
        <button id="cancelButton">Cancel</button>
      </div>
    `;
  
    const confirmCheckbox = document.getElementById('confirm');
    const sendNowButton = document.getElementById('sendNowButton');
    const cancelButton = document.getElementById('cancelButton');
    const amountInput = document.getElementById('amount');
    const receiverIdInput = document.getElementById('receiverId');
    const receiverEmailInput = document.getElementById('receiverEmail');
    const receiverAccountInput = document.getElementById('receiverAccount');
    const receiverNameInput = document.getElementById('receiverName');
    const receiverNameError = document.getElementById('receiverNameError');
    const receiverIdError = document.getElementById('receiverIdError');
    const receiverEmailError = document.getElementById('receiverEmailError');
    const receiverAccountError = document.getElementById('receiverAccountError');
    const amountError = document.getElementById('amountError');
    const confirmError = document.getElementById('confirmError');
  
    confirmCheckbox.addEventListener('change', function () {
      sendNowButton.disabled = !this.checked;
      confirmError.style.display = this.checked ? 'none' : 'inline';
    });
  
    // Form validation on the "Send Now" button click
    sendNowButton.addEventListener('click', async function (e) {
      e.preventDefault();
  
      // Reset error messages
      receiverNameError.style.display = 'none';
      receiverIdError.style.display = 'none';
      receiverEmailError.style.display = 'none';
      receiverAccountError.style.display = 'none';
      amountError.style.display = 'none';
      confirmError.style.display = 'none';
  
      let isValid = true;
  
      // Validate all fields
      if (!receiverNameInput.value) {
        receiverNameError.style.display = 'inline';
        isValid = false;
      }
      const receiverId = parseFloat(receiverIdInput.value);
      if (!receiverId || receiverId <= 0) {
        receiverIdError.style.display = 'inline';
        isValid = false;
      }
      if (!receiverEmailInput.value || !validateEmail(receiverEmailInput.value)) {
        receiverEmailError.style.display = 'inline';
        isValid = false;
      }
      const receiverAccount = parseFloat(receiverAccountInput.value);
      if (!receiverAccount || receiverAccount <= 0) {
        receiverAccountError.style.display = 'inline';
        isValid = false;
      }
      const amount = parseFloat(amountInput.value);
      if (!amount || amount <= 0 || amount > userData.balance) {
        amountError.style.display = 'inline';
        isValid = false;
      }
      // If all validations pass
      if (isValid) {
        // Show confirmation popup
        const confirmation = window.confirm(
          `Are you sure you want to send $${amount.toFixed(2)} to ${receiverNameInput.value}?`
        );
  
        if (confirmation) {
          // Prepare the payload for the transaction
          const payload = {
            sender_id: parseInt(userData.user_id),
            receiver_id: receiverId,
            account_number: receiverAccount,
            amount: amount,
            remarks: `Transfer of $${amount.toFixed(2)} from ${userData.name} to ${receiverNameInput.value}`,
            dateTimeStamp: Math.floor(Date.now() / 1000), // Current Unix timestamp
          };
  
          try {
            // Hit the /handletransaction API
            const response = await fetch('http://localhost:8080/transaction', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(payload),
            });
  
            if (response.ok) {
              alert('Transaction completed successfully!');
              formContainer.innerHTML = '';
              // After transaction, refresh the transactions and dashboard
              await fetchTransactionsAndUpdateState();
            } else {
              alert('Transaction failed. Please try again.');
            }
          } catch (error) {
            console.error('Error during transaction:', error);
            alert('Transaction failed. Please try again.');
          }
        } else {
          alert('Transaction canceled.');
        }
      }
    });
  
    cancelButton.addEventListener('click', function () {
      formContainer.innerHTML = '';
      sendMoneyButton.disabled = false;
    });
  }

  async function fetchTransactionsAndUpdateState() {
    try {
      // Fetch updated user data and transactions
      const user_id = userData.user_id ;
      const userResponse = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userData.user_id, email: userData.email, password: userData.password }),
      });
      if (userResponse.ok) {
        let userDetail = await userResponse.json();
        userData.balance = userDetail.data.balance ; 
        renderDashboard(); // Refresh the dashboard
      }
    } catch (error) {
      console.error('Error fetching updated data:', error);
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
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          ${transactions
            .map((txn) => {
              let amount = txn.amount;

              // Check if the transaction is a transfer and the current user is the sender
              if (txn.remarks.includes("Transfer") && txn.remarks.includes(`from ${userData.name}`)) {
                amount = -Math.abs(txn.amount); // Make the amount negative if the user is the sender
              }

              // Format the amount for display
              const amountStyle = amount < 0 ? 'color: red;' : '';
              const statusIcon = txn.status === 'completed' ? '✔' : '✘';

              return `
                <tr>
                  <td>${new Date(txn.dateTimeStamp * 1000).toLocaleDateString()}</td>
                  <td>${txn.remarks}</td>
                  <td style="${amountStyle}">${amount < 0 ? '-' : ''}$${Math.abs(amount)}</td>
                  <td style="text-align: center;">${statusIcon}</td>
                </tr>
              `;
            })
            .join('')}
        </tbody>
      </table>
    `;
  }
  function validateEmail(email) {
    const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return re.test(String(email).toLowerCase());
  }

  renderLoginForm();
});
