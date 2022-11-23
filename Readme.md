# Diving into Go:book:
This is repo with some test tasks and my experience with learning Go!
### Content
- [Function](/function/)    
    The task was to implement a function for getting the score in the game at a specified moment and write tests for it.
- [Channels](/channels/)    
    Console utility for finding prime numbers in a range and writing them to a file. Search and write are implemented using goroutines, synchronization goes through channels, context is used for timeout interrupt.
- [Client](/batch/)    
    The task was to implement a service to optimize interaction with the service. The service processes objects in batches (the maximum size of a batch is _n_). It is known that the service processes them for a certain time _p_ and is blocked for a long time if another batch arrives at this time.
- [Rest api service](https://github.com/wildmaus/backend_test)      
    Service for working with user balances.