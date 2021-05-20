<link rel="stylesheet" href="styles.css">

# Food Delivery Client
Full Client that communicates with the backend Server system through REST API for Food Delivery Applications. Written in go/ golang with search, auto complete and suggestive search as well. The server also utilizes different flags that can be switced on and off, depending on environment settings and testing requirements.

<h2>Introduction</h2>

<p>The client app was an extension of the server application to allow different merchants on board to make their own changes to their menu items with us. As the clients' interfaces uses API keys that can be retrieved from the main server after loggin in, merchants are able to make changes to their menu items easily as long as they have their API key, and have a browser to access https://kvclient.sozoalvin.com</p>

<h2>Tech Stacks Employed</h2>
<p>Amazon Web Services EC2 - Ubuntu</p>
<p>Amazon Web Services RD2 - mySQL - Server only (client does not have a direct access to the database)</p>
<p>Amazon Web Servies Https LB - load balanced</p>

<h2>How to view the entire project in action</h2>
<p>It is best to open up both the server as well as the client along with the instructions to see the project in action.</p>
<p>Click <a href="https://kvserver.sozoalvin.com">here</a> for the backend server</p>
<p>Click <a href="https://kvclient.sozoalvin.com">here</a> for the client server</p>

<h2>Instructions on How to Use</h2>
<h4>Click on image for a direct link to the instructions in PDF format.</h4>
<a href="#"><img src="https://i.imgur.com/f8OZ7Mk.png" width="700"></a>

<h2>Development Envirnonment</h2> 
<p>For development purposes in your own environment, use the following command once you've in the working directory</p>

```go run .```

<p>You can access the server on <i>localhost</i> or <i>localhost:80</i> if the former doesn't work</p>
<p>If you are running the server the same computer/host, please consider changing the ListenAndServe argument from</p>

```http.ListenAndServe(":80", router)```

to

```http.ListenAndServe(":8080", router)``` 

<p>If you do that change, you'll want to access the server on <i>localhost:8080</i> instead.</p>

<h2>Testing for Production</h2>
<p>To check if your program is ready for production, use the following command.</p>

```go run . -productionFlg```

<p>You can access the server on <i>localhost</i> if the former doesn't work</p>
<p>When production flag is activated, your API calls would be routed via <i>https:</i> instead of <i>http</i></p>

<h2>Launching on an Actual Web Server</h2>

```go run . -productionFlag -domain www.yourdomainname.com```

<p>On an actual web server, your domain name to your server has to be specified. If you server has a valid certificate, your should be able to send your requests through https.</p>

Sample RestAPI query URL if -productionFlag and -domain are not specified

```http://localhost/api/v1/apivalidation?api=dummy_api_key_here_pleae_paste_your_own```

Sample RestAPI query URL if -productionFlag and -domain are specified

```https://www.yourdomainname.com/api/v1/apivalidation?api=dummy_api_key_here_pleae_paste_your_own```


<h2>Search Away to Discover Great Tasting Food!</h2>
<img src = "https://i.imgur.com/XFUVxvV.png" width="700">

<h2>Access Your Own Account's API</h2>
<p>Lost API? Revoke with a simple click and lock out all clients using your old API key</p>
<img src = "https://i.imgur.com/ojn0LQ4.png" width="700">

<h2>Easily Add Any Items to Cart</h2>
<img src = "https://i.imgur.com/FKjFCQ8.png" width="700">

<h2>Admin Users Get Special Settings for Service Recovery</h2>
<img src = "https://i.imgur.com/16XvZTt.png" width="700">

<h2>Regular Account Types at Cart Option</h2>
<img src = "https://i.imgur.com/pHjQuLW.png" width="700">

<h2>Customer's Checkout Page</h2>
<img src = "https://i.imgur.com/PVTs0E7.png" width="700">

<h2>Priority Queues for Service Recvery</h2>
<img src = "https://i.imgur.com/c63vP9z.png" width="700">

<h2>Driver Assignment</h2>
<img src = "https://i.imgur.com/1Nye0iX.png" width="700">
<img src = "https://i.imgur.com/bhKLJWl.png" width="700">

<h2>Read, Get Set, Dispatch and it's a GO!</h2>
<img src = "https://i.imgur.com/xJjm2yX.png" width="700">


<h2>Easily View All System Information</h2>
<img src = "https://i.imgur.com/DCrvoAL.png" width="700">
<img src = "https://i.imgur.com/qhb7rVt.png" width="700">

