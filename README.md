# Traffic Tracker

Small utility to track statistics about network traffic going in and out of multiple interfaces.

The tool is based around a common data unit called record. The record defines a packet that went through an interface. It specifies the network interface, length of the packet, source, destination and the different protocols, among other information.

With `record` data from multiple interfaces can be recorded into a JSON file, which can be then rendered to an HTML file using `visualize`. `record` prompts the user for the interfaces to listen to, while `visualize` uses a flag (`-i <file_name>`) to specify the input log file.

`record` uses PCAP through [gopacket](https://github.com/google/gopacket) for traffic recording and `visualize` uses [go-echarts](https://github.com/go-echarts/go-echarts) to render the charts.

## Usage

First run `record`. This will follow with multiple promts for each of the detected network interfaces available. Once this happens, `record` will listen to those interfaces and store all the data on a json file under the folder `records`. The file name has a name of the form `traffic_<timestamp>.json`, where timestamp is the Unix timestamp. Each minute the command will print how much data has been written in bytes.

After a record is made, run `visualize`. `visualize` requires to pass the flag `-i`, followed by the record file for which data will be plotted. This will output an `html` file that can be opened using any browser to see the chart. Configuring the data visualization can only be done through code.

## Authors

[Juan Martinez](https://github.com/jmaralo)
