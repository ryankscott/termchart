# termchart

Termchart is a charting app that allows you to quickly visualise data from STDIN. It generates simple bar charts of data that are displayed directly in your terminal.
Termchart expects the data to be formatted in two columns, the first a string which is used to label the bar, the second an integer used as the value of the bar. For example:

| label | value |
| ----- | ----- |
| foo   | 123   |
| bar   | 456   |

This would produce a bar chart with two bars, one bar called "foo" with the value 123, and one bar called "bar" with the value 456. 

Please note that termchart expects no header data.

## Getting Started

### Installing

Use `go get` to get the latest version of termchart

```
go get github.com/ryankscott/termchart
```

### Usage
```
echo -e "foo,123\nbar,456" | termchart
cat test/foobar.tsv | termchart

etc.
```

Works well with https://github.com/MarianoGappa/sql


## Contributing

Any pull requests are welcome.

## Authors

* **Ryan Scott** - ryan@ryankscott.com


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

This project heavily leverages termui https://github.com/gizak/termui

