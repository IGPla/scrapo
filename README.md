# Scrapo

A high performance, fully customizable web scraping framework

## Why Scrapo?

Scrapo has born as a high performance scraping framework, over a language that offers a nice trade-off between performance, concurrency and flexibility.

As you will see through this documentation, core ideas are easy to learn, and all them are already implemented and ready to use for your scraping projects to get the results you are waiting for

It will take just few lines of code to setup your first scrapo project :)

## Quick start

Download the project

```
go get github.com/IGPla/scrapo
```

Download all dependencies (if you have not done it yet)

```
cd scrapo/
go get ./...
```

scraposample.go is a good place to start from. Let's look inside

```
package main

import (
	"github.com/igpla/scrapo/config"
	"github.com/igpla/scrapo/crawler"
	"github.com/igpla/scrapo/processor"
	"github.com/igpla/scrapo/storage"
)

func main() {
	// Base settings
	config.PopulateMainConfig("sampleconfig.json")
	// Storage choice and configuration
	var fileStorage *storage.FileStorage = new(storage.FileStorage)
	fileStorage.Prefix = "/tmp/myscrapingproject"
	config.MainConfig.Storage = fileStorage
	// Base HTML class to parse content
	config.MainConfig.ResourceProcessor = new(processor.PlainResourceProcessor)
	// Start crawler
	crawler.Start()
}

```

This basic scrapo task will get the config file to set up all settings, start the crawling process using a PlainResourceProcessor (we will talk further about it in the (#Processors) section), and finally, store all results under the filesystem (/tmp/myscrapingproject)

Finally, to run your scrapo project, just type the following

```
go run myscrapoproject.go
```

We will take a look at (#Config) for all configuration options available. Available storages and all their properties will be discussed in (#Storage) section

## Config

Scrapo project configuration is done through a .json file. The main function to load all settings, and the one that should be the first function call in your main file, is the following:

```
config.PopulateMainConfig("/path/to/settings.json")
```

Inside your settings file, the following content should be typed:

```
{
    "user_agent": "<USER AGENT STRING>",
    "urls": [
        "<URL1>",
        "<URL2>",
        ...
        "<URLN>"
    ],
    "allowed_domains": [
        "<DOMAIN1>",
        "<DOMAIN2>",
        ...
        "<DOMAINM>"
    ],
    "num_workers": <NUMBER OF CONCURRENT WORKERS>
}
```

Let's take a look at each one

- user\_agent: this will be the user agent string that crawler process will use to identify itself to each resource. 
- urls: this will be the starting point for your scrapo project, being the first URLs to be crawled.
- allowed\_domains: this list of domains will be the responsible to restrict which resources will be crawled and which ones will be discarded. All URLs must match one of the provided domains here to be processed.
- num\_workers: the number of effective worker instances to perform crawling. More workers often lead to lower crawling times (this is not always true, due to the fact that num\_workers will increase the concurrent factor, but it depends on the underlying data graph for the final performance)


## Storage

Storage will be a piece of code that will define where your results will be stored. You can provide your own storage implementations, just implementing the storage.Storage interface. 

Additionally, you've already implemented several storages for common situations. It's up to you to use them of implement your own for a more customized experience.

After storage implementation, it must be configured for your project.

An important step that you must take after configure your storage is to assign that storage to scrapo project's configuration. 

Right now, we've the following storage implementations

- File: this storage will store all content to filesystem. It will take the following parameters:
  - Prefix: A filesystem base path, where all content will be downloaded

The following snippet is an example of File storage configuration
```
var fileStorage *storage.FileStorage = new(storage.FileStorage)
fileStorage.Prefix = "/tmp/myscrapingproject"
config.MainConfig.Storage = fileStorage
```

- Redis: this storage allows your scrapo project to store all resources in Redis database. It is a good choice for a data pipeline, using scrapo to feed redis and use it as a buffer, where another tool gets data from there for further processing.
  - Host: host where redis can be accessed
  - Port: port where redis is listening
  - Password: password to access redis host. Leave it empty if redis is open
  - Db: db number to use in the redis host
  - DefaultExpiration: default expiration (in seconds) for all keys
  
Note: as redis is an in memory database, be careful using this storage. If you download entire website, with all resources (say images, sounds and videos), it can lead to a crash job due to run in out of memory issues.

## Processors

Processors are the backbone of scrapo project. They are responsible of extracting meaningful knowledge for your project.

To implement your processor, just be sure to satisfy this interface

```
type ResourceProcessor interface {
	Parse(task *tasks.Task) (parsedResources []*tasks.Task, newRawResources []*tasks.Task)
}
```

Return parameters should contain the following

- parsedResources: a list of Task instances, ready to be stored. For example, in a minimal project where you want only the titles of the .html pages, parsedResources will contain a single Task with .html page title as the Content field. (Important note: if you want to return different resources to be stored, be sure that you provide different URLs, as they will be used as identifiers. One common pattern is to use the same URL for all them and append a suffix distinguish each resource coming from the same source)

- newRawResources: a list of all new resources that should be processed. These items will be the base for the next crawling steps. The most simple crawler will return here a slice of Task items, each one with URL field pointing to .html file links (href properties of anchor tags). It's important to say that URL is the only expected parameter in the Tasks retrieved; the rest of them will be overridden by the crawling process.

Finally, your processor must be assigned to scrapo project configuration, as follows

```
config.MainConfig.ResourceProcessor = new(YourProcessorImplementation)
```

There is a basic processor implemented, just to be useful as a sample of how to implement your own processors. Let's take a look on it

```
func (dhp DownloadResourceProcessor) Parse(task *tasks.Task) (parsedResources []*tasks.Task, newRawResources []*tasks.Task) {
	var links []string = GetHTMLElements(bytes.NewBuffer(task.Content), "a", "href")
	var images []string = GetHTMLElements(bytes.NewBuffer(task.Content), "img", "src")
	var scripts []string = GetHTMLElements(bytes.NewBuffer(task.Content), "script", "src")
	var csss []string = GetHTMLElements(bytes.NewBuffer(task.Content), "link", "href")

	processorLogger.Printf("Got %d links, %d images, %d scripts and %d css",
		len(links), len(images), len(scripts), len(csss))

	newRawResources = make([]*tasks.Task, 0, len(images)+len(links)+len(scripts)+len(csss))

	for _, link := range links {
		newRawResources = append(newRawResources, tasks.NewTask(link))
	}
	for _, image := range images {
		newRawResources = append(newRawResources, tasks.NewTask(image))
	}
	for _, script := range scripts {
		newRawResources = append(newRawResources, tasks.NewTask(script))
	}
	for _, css := range csss {
		newRawResources = append(newRawResources, tasks.NewTask(css))
	}
	parsedResources = make([]*tasks.Task, 0, 1)
	parsedResources = append(parsedResources, task)
	return parsedResources, newRawResources
}
```

This processor will download all basic assets. It will simply parse each html element, looking for new links, images, scrips and css stylesheets, and then, just download them with the provided storage.

As you can see, there are already implemented helper functions to get content tags and properties (GetHTMLElements), ready to be used for your own implementation.

The return parameters of this processor will be the following

- parsedResouces: a slice with a single Task. This task will contain the .html content
- newRawResources: a slice with as many tasks as all links, images, scripts and css style sheets in the parsed page.

To wrap up, a Processor will be the main function that you will be implementing to define your scrapo project target.

## Final notes

- To run tests, just type the following command

```
go test ./...
```

- Redis storage test relies on docker image of redis. Take it into account before running tests (it will require docker installed to be successful)
