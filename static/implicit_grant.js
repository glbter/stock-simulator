// Check if the user is already authenticated
if (localStorage.getItem('accessToken')) {
    // Call the fetchStockData function to populate the stock list if the user is authenticated
    console.log(localStorage.getItem('accessToken'))
} else {

    const token = window.location.hash.substring(1);

    if (token) {
        localStorage.setItem('accessToken', token)

    } else {
        window.location.replace("https://dev-r1fw6swpk4dcvcyi.us.auth0.com/authorize?client_id=4bFggkRTQofaujBo0quaHkgpx2UfoDyr&response_type=token&audience=https://stocks-simulatior/&redirect_uri=https://stocks-spa.s3.eu-central-1.amazonaws.com/implicit_grant.html");

    }

}