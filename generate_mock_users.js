const fs = require('fs');
const { faker } = require('@faker-js/faker');

// Function to generate user data
function generateUserData(numUsers) {
    const users = [];
    const generatedAccountNumbers = new Set();

    for (let i = 0; i < numUsers; i++) {
        // Generate a unique account number
        let accountNumber;
        do {
            accountNumber = faker.number.int({ min: 100000000, max: 999999999 });
        } while (generatedAccountNumbers.has(accountNumber));
        generatedAccountNumbers.add(accountNumber);

        const user = {
            user_id: 100 + i, // Auto-incremented number starting at 100
            first_name: faker.person.firstName(),
            last_name: faker.person.lastName(),
            email: faker.internet.email(),
            current_balance: faker.number.int({ min: 10000, max: 999999 }),
            password: faker.internet.password(8, false, /[A-Za-z0-9]/), // Random alphanumeric password
            account_number: accountNumber,
        };

        users.push(user);
    }

    return users;
}

// Function to write user data to a JSON file
function writeDataToJson(data, filename) {
    fs.writeFileSync(filename, JSON.stringify(data, null, 2), "utf-8");
    console.log(`Data successfully written to ${filename}`);
}

// Main
const numUsers = 3000;
const filename = "mock_data_userInfo.json";

console.log("Generating user data...");
const userData = generateUserData(numUsers);

console.log("Writing data to JSON file...");
writeDataToJson(userData, filename);
