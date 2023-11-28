# Orderbook matching engine server

### Security

The orderbook matching engine microservice can be called only from the internal servers. Token validation outsourced to the first server in the internal group and than user's data is optimistically trusted (after the first server validated it by calling Auth
microservice). This is done to optimize load on the Auth server, reduce unnecessary requests and simplify authentication validation.

The enpoints should be protected by firewall upon deployment and called only from trusted sources (other microservices).

