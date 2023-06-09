const host = 'https://5nd4ve1r26.execute-api.eu-central-1.amazonaws.com'

const clientHost = 'http://stocks-spa.s3-website.eu-central-1.amazonaws.com/'

let isString = value => typeof value === 'string' || value instanceof String;

async function requestStockData(tickerId, ticker) {
    var url = host+`/api/ticker`
    if (isString(tickerId) && tickerId !== '') {
        url += '?ticker_ids='+tickerId
    }

    if (isString(ticker) && ticker !== '') {
        url += '?ticker_names='+ticker
    }

    const res = await axios.get(url);
    return res.data
}

async function fetchStockData(ticker, tickerName) {
    try {
        const stockList = document.getElementById('stockList');

        stockList.innerHTML = '';

        const listItem = document.createElement('li');
        const symbolSpan = document.createElement('span');
        const priceHighSpan = document.createElement('span');
        const changeLowSpan = document.createElement('span');

        symbolSpan.textContent = "Код активу";
        priceHighSpan.textContent = "Вища ціна";
        changeLowSpan.textContent = "Нижча ціна";

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

        stockList.appendChild(listItem);

        const data = await requestStockData(ticker, tickerName);
        data.forEach(stock => {
            const { Ticker, High, Low, Open, Close, Volume } = stock;
            const { ID, Name, Description } = Ticker;

            const divBtn = document.createElement('div')
            divBtn.onclick = function (tickerID) {
                return function () {
                    showTickerDetailsMenu(tickerID)
                }
            }(ID)
            const listItem = document.createElement('li');
            const symbolSpan = document.createElement('span');
            const priceHighSpan = document.createElement('span');
            const changeLowSpan = document.createElement('span');
            const button = document.createElement('button')

            button.textContent = "Торгувати";
            button.onclick = function (tickerID, tickerName) {
                return function (){
                    showTradeMenu(tickerID, tickerName)
                }
            }(ID, Name)

            symbolSpan.textContent = Name;
            priceHighSpan.textContent = appendUAHCurrency(High);
            changeLowSpan.textContent = appendUAHCurrency(Low);

            listItem.appendChild(symbolSpan);
            listItem.appendChild(priceHighSpan);
            listItem.appendChild(changeLowSpan);

            listItem.appendChild(button)

            stockList.appendChild(listItem);
        });
    } catch (error) {
        console.log('Error fetching stock data:', error);
    }
}

async function showMainMenu() {
    sessionStorage.setItem('currentPageId', 'main_menu')

    document.getElementById('trade_menu').style.display = 'none'
    document.getElementById('portfolio_menu').style.display = 'none'
    document.getElementById('main-menu-return-btn').style.display = 'none'
    document.getElementById('ticker_details_menu').style.display = 'none'

    document.getElementById('main_menu').style.display = 'inherit'


    fetchStockData()
}

async function showInvestmentPortfolioMenu() {
    sessionStorage.setItem('currentPageId', 'portfolio_menu')

    handleAuthorization()

    document.getElementById('main_menu').style.display = 'none'
    document.getElementById('trade_menu').style.display = 'none'
    document.getElementById('ticker_details_menu').style.display = 'none'

    document.getElementById('portfolio_menu').style.display = 'inherit'
    document.getElementById('main-menu-return-btn').style.display = 'inherit'

    try {
        const stockList = document.getElementById('stockListInvestment');

        stockList.innerHTML = '';

        const res = await axios
            .get(host+'/api/portfolio')
            .catch((err) => {
                if (err.response && err.response.status === 401) {
                    reAuthorize()
                }
            });

        const data = await res.data

        //List headers
        addToPortfolioTickerList(stockList, "Код активу", "Вища ціна", "Нижча ціна", "Кількість")

        //Total info
        addElementToPortfolioTickerList(stockList, '', data.Total.High, data.Total.Low, data.Total.Amount)

        data.All.forEach(stock => {
            const { TickerID, Name, High, Low, Open, Close, Amount } = stock;
            addElementToPortfolioTickerList(stockList, Name, High, Low, Amount)
        });
    } catch (error) {
        console.log('Error fetching stock data:', error);
    }
}

function addToPortfolioTickerList(stockList, Name, High, Low, Amount) {
    const listItem = document.createElement('li');
    const symbolSpan = document.createElement('span');
    const priceHighSpan = document.createElement('span');
    const changeLowSpan = document.createElement('span')
    const amountSpan = document.createElement('span');

    symbolSpan.textContent = Name;
    priceHighSpan.textContent = High;
    changeLowSpan.textContent = Low;
    amountSpan.textContent = Amount;

    listItem.appendChild(symbolSpan);
    listItem.appendChild(priceHighSpan);
    listItem.appendChild(changeLowSpan);
    listItem.appendChild(amountSpan);

    stockList.appendChild(listItem);
}

function addElementToPortfolioTickerList(stockList, Name, High, Low, Amount) {
    addToPortfolioTickerList(stockList, Name, appendUAHCurrency(High), appendUAHCurrency(Low), Amount)
}

async function searchTicker() {
    const input = document.getElementById('search-ticker-input').value

    if (!isString(input) || input === "" || !/^[a-zA-Z]+$/.test(input)) {
        return
    }

    fetchStockData("", input)
}

async function showTradeMenu(tickerID, tickerName) {
    sessionStorage.setItem('currentPageId', 'trade_menu')

    handleAuthorization()

    document.getElementById('main_menu').style.display = 'none'
    document.getElementById('portfolio_menu').style.display = 'none'
    document.getElementById('ticker_details_menu').style.display = 'none'

    document.getElementById('trade_menu').style.display = 'inherit'
    document.getElementById('main-menu-return-btn').style.display = 'inherit'

    if (!isString(tickerName)) {
        tickerName = ''
    }
    document.getElementById("trade-ticker-code").value = tickerName
}

async function tradeTicker() {
    const action = document.getElementById("action").value
    const ticker = document.getElementById("trade-ticker-code").value
    const amount = document.getElementById("trade-ticker-amount").value

    if (!isString(ticker) || ticker === "" || !/^[a-zA-Z]+$/.test(ticker)) {
        showTradeError("Неправильний формат коду активу")
        return
    }

    const tickerAmount = parseInt(amount, 10)
    if (!(Number.isInteger(tickerAmount) && tickerAmount > 0)) {
        showTradeError("Неправильний формат кількості")
        return
    }

    try {
        const data = await requestStockData('', ticker)
        if (data.length === 0) {
            showTradeError("Актив не знайдено")
            return
        }

        await axios.
            post(host+'/api/ticker/trade',{
                amount: tickerAmount,
                action: action,
                ticker_id: data[0].Ticker.ID
            })
            .catch((err) => {
                if (err.response && err.response.status === 401) {
                    reAuthorize()
                }
            })

    } catch (e) {
        showTradeError("Виникла помилка")
    }

}

function showTradeError(text) {
    document.getElementById("trade-result-success").style.display = 'none'
    const errMsg = document.getElementById("trade-result-fail")
    errMsg.textContent = text
    errMsg.style.display = 'inherit'
}

function showTradeSuccess(text) {
    document.getElementById("trade-result-fail").style.display = 'none'
    const successMsg = document.getElementById("trade-result-success")
    successMsg.textContent = text
    successMsg.style.display = 'inherit'
}

async function showTickerDetailsMenu(tickerID) {
    sessionStorage.setItem('currentPageId', 'ticker_details_menu')

    document.getElementById('main_menu').style.display = 'none'
    document.getElementById('portfolio_menu').style.display = 'none'
    document.getElementById('trade_menu').style.display = 'none'

    document.getElementById('ticker_details_menu').style.display = 'inherit'
    document.getElementById('main-menu-return-btn').style.display = 'inherit'

    const data = await requestStockData(tickerID, '')
    if (data.length === 0) {
        return
    }

    const tk = data[0]

    document.getElementById("ticker-details-ticker-name").value = tk.Ticker.Name
    document.getElementById("ticker-details-ticker-high").value = appendUAHCurrency(tk.high)
    document.getElementById("ticker-details-ticker-low").value = appendUAHCurrency(tk.low)
    document.getElementById("ticker-details-ticker-open").value = appendUAHCurrency(tk.open)
    document.getElementById("ticker-details-ticker-close").value = appendUAHCurrency(tk.close)
    document.getElementById("ticker-details-ticker-volume").value = appendUAHCurrency(tk.volume)
}

async function showWebsite() {
    const pageID = sessionStorage.getItem('currentPageId')

    if (pageID === 'main_menu' || pageID === null || pageID === '') {
        await showMainMenu()
    }

    if (pageID === 'trade_menu') {
        await showTradeMenu()
    }

    if (pageID === 'portfolio_menu') {
        await showInvestmentPortfolioMenu()
    }
}

function handleAuthorization() {
    if (localStorage.getItem('accessToken')) {
        const token = localStorage.getItem('accessToken')
        axios.defaults.headers.common['Authorization'] = 'Bearer '+token
    } else {
        // get token from url
        const token = window.location.hash.substring(1).split("&")[0].split("=")[1];
        if (token) {
            localStorage.setItem('accessToken', token)
            axios.defaults.headers.common['Authorization'] = 'Bearer '+token
        } else {
            window.location.replace("https://dev-r1fw6swpk4dcvcyi.us.auth0.com/authorize?client_id=4bFggkRTQofaujBo0quaHkgpx2UfoDyr&response_type=token&audience=https://stocks-simulatior/&redirect_uri="+clientHost);
        }
    }
}

function reAuthorize() {
    localStorage.removeItem('accessToken')
    window.location.replace("https://dev-r1fw6swpk4dcvcyi.us.auth0.com/authorize?client_id=4bFggkRTQofaujBo0quaHkgpx2UfoDyr&response_type=token&audience=https://stocks-simulatior/&redirect_uri="+clientHost);
}

function appendUAHCurrency(val) {
    return new Intl.NumberFormat('uk-UA', { style: 'currency', currency: 'UAH' }).format(val)
}

showWebsite();
