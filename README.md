## KafkaFilter

Search through kafka topic for a word/phrase, live stream search as well.

#### Setup

Requires an existing Kafka Cluster + topics to search.

Run main.go, access website/index.html via browser

##### Disclaimers

Currently only setup to take api server on `localhost:5000`

Web interface doesn't allow changing topics, API server takes topic as part of the URL.

Web interface is also quite ugly... sorry

Only got a 2-3 hours in on this, and spent some of that on trying a new go kafka client.


### ScreenShots

![Initial Screen](/screenshots/Initial.png)

![Search for Devops](/screenshots/Search-Devops.png)

