# Channels
Console utility for finding prime numbers in a range and writing them to a file. Search and write are implemented using goroutines, synchronization goes through channels, context is used for timeout interrupt. Found numbers are written to the specified file one per line.      
Flag `--filename` used for specify output file, if file already exist founded numbers will be just added in the end. Default `primes.txt`.      
Flag `--timeout` used for specify timeout in seconds for programm. Default 10.      
Flag `--range` can be used multiple times, for each input range starts it own gorutine.        
### Commands
```bash
go run channels.go --filename <filename> --timeout <timeout> --range <start:end>
```