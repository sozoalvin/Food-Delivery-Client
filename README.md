<link rel="stylesheet" href="styles.css">

# Food Delivery Client
Full Client that communicates with the backend Server system through REST API for Food Delivery Applications. Written in go/ golang with search, auto complete and suggestive search as well. The server also utilizes different flags that can be switced on and off, depending on environment settings and testing requirements.

<img src = "https://i.imgur.com/2rdSode.png" width="700">

<h2>Introduction</h2>

<p>The client app was an extension of the server application to allow different merchants on board to make their own changes to their menu items with us. As the clients' interfaces uses API keys that can be retrieved from the main server after loggin in, merchants are able to make changes to their menu items easily as long as they have their API key, and have a browser to access https://kvclient.sozoalvin.com</p>

<h2>Tech Stacks Employed</h2>
<p>Amazon Web Services EC2 - Ubuntu</p>
<p>Amazon Web Services RD2 - mySQL - Server only (client does not have a direct access to the database)</p>
<p>Amazon Web Servies Https LB - load balanced</p>

<h2>View the Entire Project in Action</h2>
<p>It is best to open up both the server as well as the client along with the instructions to see the project in action.</p>
<p>Click <a href="https://kvserver.sozoalvin.com">here</a> for the backend server</p>
<p>Click <a href="https://kvclient.sozoalvin.com">here</a> for the client server</p>

<h2>Understand How RESTAPI works in the Project</h2>
<p>Click<a href="https://github.com/sozoalvin/Food-Delivery-Client/blob/master/Instructions/Understanding_RESTAPI.pdf"> Here </a> to learn more about how RESTAPI was deployed on the backend server as well as the client's</p>

<h2>Instructions on How to Use</h2>
<h4>Click on image for a direct link to the instructions in PDF format.</h4>
<a href="https://github.com/sozoalvin/Food-Delivery-Client/blob/master/Instructions/Instructions%20on%20Navigating%20Client%20Program%20KV%20Food%20Delivery%20.pdf"><img src="https://i.imgur.com/sgKILTQ.png" width="700"></a>

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

<h2>API Validating is Easy</h2>
<img src = "https://i.imgur.com/izEYc3a.png" width="700">

<h2>Merchant Account Setup</h2>
<p>Set up your merchant account with your brand's name, detailed location, postal code and operating hours. Now it even includes a cut off timing!</p>
<img src = "https://i.imgur.com/d0onXT9.png" width="700">

<h2>Adding a Food Menu Item</h2>
<img src = "https://i.imgur.com/PsVqYUA.png" width="700">
<p>Sample food item being added</p>
<img src ="https://i.imgur.com/eWoAvl1.png" width="700"> 
<p>Did you know that items added by merchants are immediately searchable and indexed by the main server?</p>
<img src = "https://i.imgur.com/usoiJz7.png" width="700"> 
<img src = "https://i.imgur.com/X74JVYF.png" width="700"> 

<h2>Editing a Food Menu Item</h2>
<img src = "https://i.imgur.com/AswqL4D.png" width="700">
<p>Editing a food menu item doens't have to be complicated when it's only 2 clicks + 2 forms</p>
<img src ="https://i.imgur.com/QUHaoya.png" width="700"> 

<h2>Deleting a Food Menu Item</h2>
<p>Deleting a food menu item is also very straightforward</p>
<img src = "https://i.imgur.com/EVAwxEz.png" width="700">

<h2>View All Menu Records</h2>
<img src = "https://i.imgur.com/6HBNKx3.png" width="700"> 


