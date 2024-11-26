const fs = require('fs');

// Input JSON file containing user data
const userFile = "mock_data_userInfo.json"; // Replace with your file name
const outputFile = "mock_transactions.json";

// Load user data
const users = JSON.parse(fs.readFileSync(userFile, 'utf-8'));

// Extract user IDs
const userIds = users.map(user => user.user_id);
const userMap = Object.fromEntries(users.map(user => [user.user_id, `${user.first_name} ${user.last_name}`]));

// Function to generate a random timestamp between two dates
function randomDate(start, end) {
    const startTime = start.getTime();
    const endTime = end.getTime();
    return new Date(startTime + Math.random() * (endTime - startTime));
}

// Transaction types
const transactionTypes = ["transfer", "deposit", "withdrawal"];

// Generate and write transactions in batches
function generateTransactions(numTransactions, batchSize) {
    const startDate = new Date(2022, 0, 1); // January 1, 2022
    const endDate = new Date();

    // Create a write stream
    const stream = fs.createWriteStream(outputFile, { flags: 'w', encoding: 'utf-8' });

    stream.write('[\n'); // Start JSON array

    let isFirstTransaction = true;

    for (let i = 0; i < numTransactions; i += batchSize) {
        const transactions = [];
        const upperLimit = Math.min(i + batchSize, numTransactions);

        for (let txnId = i + 1; txnId <= upperLimit; txnId++) {
            const txnType = transactionTypes[Math.floor(Math.random() * transactionTypes.length)];
            const senderId = userIds[Math.floor(Math.random() * userIds.length)];
            let receiverId = txnType === "transfer" ? userIds[Math.floor(Math.random() * userIds.length)] : senderId;
            let amount = Math.floor(Math.random() * (15000 - 500 + 1)) + 500;
            const timestamp = Math.floor(randomDate(startDate, endDate).getTime() / 1000); // Unix timestamp in seconds
            const status = Math.random() < 0.02 ? "failed" : "completed";
            const senderName = userMap[senderId];
            const receiverName = userMap[receiverId];

            let remarks;
            if (txnType === "deposit") {
                remarks = `Deposit of ${amount} by ${senderName}`;
                receiverId = senderId; // Same user ID for deposit
            } else if (txnType === "withdrawal") {
                remarks = `Withdrawal of ${amount} by ${senderName}`;
                receiverId = senderId; // Same user ID for withdrawal
                amount = -amount; // Withdrawal amount as negative
            } else { // Transfer
                remarks = `Transfer of ${amount} from ${senderName} to ${receiverName}`;
            }

            const transaction = {
                transaction_id: txnId,
                sender_id: senderId,
                amount: amount,
                receiver_id: receiverId,
                remarks: remarks,
                dateTimeStamp: timestamp,
                status: status
            };

            transactions.push(transaction);
        }

        // Write batch to file
        const jsonLines = transactions.map(txn => JSON.stringify(txn)).join(',\n');
        if (!isFirstTransaction) {
            stream.write(',\n');
        }
        stream.write(jsonLines);

        isFirstTransaction = false;

        console.log(`Written ${Math.min(i + batchSize, numTransactions)} of ${numTransactions} transactions...`);
    }

    stream.write('\n]\n'); // Close JSON array
    stream.end();

    console.log(`Mock transaction data generated successfully in '${outputFile}'.`);
}

// Main
const numTransactions = 450000; // Set this to the desired number of transactions
const batchSize = 10000; // Adjust batch size to balance memory usage and write speed

console.log("Generating and writing transactions...");
generateTransactions(numTransactions, batchSize);
