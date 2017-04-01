# tws-ov-convert
This tool helps to convert exported tradelog files from Interactive Brokers TWS into optionvue format.

## Setup optionvue
Follow the instructions at http://help.capitaldiscussions.com/article/how-to-import-trades-from-interactive-brokers-to-optionvue to setup the optionvue importer.

## Setup tws tradelog export
![export column settings](https://cloud.githubusercontent.com/assets/9795022/24576303/793d4ab4-166e-11e7-8b8f-b38049796c33.png)

![export trade reports](https://cloud.githubusercontent.com/assets/9795022/24576305/7bc44166-166e-11e7-91b6-add124840ecc.png)

## Configuration
Copy config.json.example to config.json and adjust as needed.

- Use always "/" as path separator
- Keep always one multiplier entry

## Expected input file name
Currently the tool expect the default filename from TWS which is, "trades.YYYYMMdd.csv".

## Multiplier
Per default the application uses 100. I added 50 for ES futures to the config file. The multipliers array can be extended to support other Underlyings.

