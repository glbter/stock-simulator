// Simulated stock data
const data = [
    {
        "Ticker": {
            "ID": "3439d561-b4db-4455-aff9-da2119573574",
            "Name": "AAPL2",
            "Description": ""
        },
        "High": 3,
        "Low": 2,
        "Open": 4,
        "Close": 5,
        "Volume": 0,
        "DataDate": "2010-01-20T00:00:00Z"
    },
    {
        "Ticker": {
            "ID": "aad17418-6764-4ecd-90ed-bb1d7091edcc",
            "Name": "AAPL",
            "Description": ""
        },
        "High": 21,
        "Low": 19,
        "Open": 19.5,
        "Close": 20,
        "Volume": 50,
        "DataDate": "2023-05-26T00:00:00Z"
    }
    , {
        "Ticker": {
            "ID": "aad17418-6764-4ecd-90ed-bb1d7091edcc",
            "Name": "TSLA",
            "Description": ""
        },
        "High": 21,
        "Low": 19,
        "Open": 19.5,
        "Close": 20,
        "Volume": 50,
        "DataDate": "2023-05-26T00:00:00Z"
    }
];


// Function to fetch stock data and populate the list
async function fetchStockData() {
    try {
        // Get the stock list element
        const stockList = document.getElementById('stockList');

        // Clear any existing stock items
        stockList.innerHTML = '';

        // Create list item elements
        const listItem = document.createElement('li');
        const symbolSpan = document.createElement('span');
        const priceHighSpan = document.createElement('span');
        const changeLowSpan = document.createElement('span');

        // Set the stock data to the element text content
        symbolSpan.textContent = "Код активу";
        priceHighSpan.textContent = "Вища ціна";
        changeLowSpan.textContent = "Нижча ціна";

        // Append the spans to the list item
        listItem.appendChild(symbolSpan);
        listItem.appendChild(priceHighSpan);
        listItem.appendChild(changeLowSpan);

        const button = document.createElement('button')

        button.textContent = "Торгувати";
        button.onclick = function () {
            location.href='https://stocks-spa.s3.eu-central-1.amazonaws.com/implicit_grant.html'
        }

        button.style.backgroundColor = "white"
        listItem.appendChild(button)

        // Append the list item to the stock list
        stockList.appendChild(listItem);

        // Iterate through the stock data and create list items
        data.forEach(stock => {
            const { Ticker, High, Low, Open, Close, Volume } = stock;
            const { ID, Name, Description } = Ticker;

            // Create list item elements
            const listItem = document.createElement('li');
            const symbolSpan = document.createElement('span');
            const priceHighSpan = document.createElement('span');
            const changeLowSpan = document.createElement('span');
            const button = document.createElement('button')

            button.textContent = "Торгувати";
            button.onclick = function () {
                location.href='https://stocks-spa.s3.eu-central-1.amazonaws.com/implicit_grant.html'
            }

            // Set the stock data to the element text content
            symbolSpan.textContent = Name;
            priceHighSpan.textContent = High+'₴';
            changeLowSpan.textContent = Low+'₴';

            // Append the spans to the list item
            listItem.appendChild(symbolSpan);
            listItem.appendChild(priceHighSpan);
            listItem.appendChild(changeLowSpan);

            listItem.appendChild(button)


            // Append the list item to the stock list
            stockList.appendChild(listItem);
        });
    } catch (error) {
        console.log('Error fetching stock data:', error);
    }
}

fetchStockData();