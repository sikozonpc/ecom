## Here is a list of improvements that could be made to the project:

- Refactor the product quantity/stock to be an [atomic and concurrent safe](https://www.freecodecamp.org/news/acid-databases-explained/#what-does-atomicity-mean). The current implementation might lead to invalid stock values if multiple requests are made at the same time. Not only that, every time a product quantity is updated we need to query the products table *(violates database normalization principle)*. 

- Implement the user addresses feature. This way we can store the user's address and use it in the checkout instead of an hardcoded value.

- Implement the user's order history. This way we can store the user's orders and show them in the user's profile.

- Implement the cancel order endpoint. This way we can allow the user to cancel an order if it's not yet shipped.