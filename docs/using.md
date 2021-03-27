# Usage

This is how you can use it.

```bash
clatest --help
This command line tool can be used to download the latest Covid related 
statistics. The data is downloaded from disease.sh and is sourced from John 
Hopkins.

Usage:
  clatest [flags]

Flags:
      --format string   Output format (markdown, csv, tab) (default "markdown")
  -f, --from string     first date to download data for (default "2021-03-26")
  -h, --help            help for clatest
  -o, --on string       A single date to get
  -t, --to string       last date to download data for (default "2021-03-27")
```

```bash
./clatest united states
```

  DATE       | CASES    | DEATHS | RECOVERED  
-------------|----------|--------|------------
  2021-03-26 | 30156621 | 548087 | 0 


```bash
./clatest united states --from 2021-03-01
```

  DATE       | CASES    | DEATHS | RECOVERED  
-------------|----------|--------|------------
  2021-03-01 | 28705285 | 515524 | 0          
  2021-03-02 | 28762326 | 517467 | 0          
  2021-03-03 | 28829520 | 519957 | 0          
  2021-03-04 | 28897518 | 521888 | 0          
  2021-03-05 | 28963921 | 523663 | 0          
  2021-03-06 | 29022116 | 525162 | 0          
  2021-03-07 | 29063082 | 525844 | 0          
  2021-03-08 | 29108096 | 526574 | 0          
  2021-03-09 | 29165791 | 528369 | 0          
  2021-03-10 | 29223730 | 529926 | 0          
  2021-03-11 | 29286134 | 531487 | 0          
  2021-03-12 | 29347338 | 533050 | 0          
  2021-03-13 | 29400553 | 534780 | 0          
  2021-03-14 | 29438775 | 535356 | 0          
  2021-03-15 | 29495424 | 536098 | 0          
  2021-03-16 | 29549364 | 537259 | 0          
  2021-03-17 | 29608458 | 538434 | 0          
  2021-03-18 | 29668959 | 540050 | 0          
  2021-03-19 | 29730486 | 541144 | 0          
  2021-03-20 | 29785935 | 541920 | 0          
  2021-03-21 | 29819701 | 542363 | 0          
  2021-03-22 | 29871268 | 542922 | 0          
  2021-03-23 | 29924892 | 543810 | 0          
  2021-03-24 | 30011839 | 545264 | 0          
  2021-03-25 | 30079282 | 546822 | 0          
  2021-03-26 | 30156621 | 548087 | 0

```bash
./clatest united states --from 2021-03-01 --to 2021-03-03
```

  DATE       | CASES    | DEATHS | RECOVERED  
-------------|----------|--------|------------
  2021-03-01 | 28705285 | 515524 | 0          
  2021-03-02 | 28762326 | 517467 | 0          
  2021-03-03 | 28829520 | 519957 | 0   

