# invindex
github.com/polisgo2020/search-senyast4745 implements inverted index to perform full-text search.


## Build application

```shell script
make build
```

## Usage

### Run build index

```shell script
./search build --soruces /path/to/folder/to/index --index /index/file/path
```

### Run search

The program can be launched in two ways

#### Console search

```shell script
export LISTEN=inteface-to-listen
export LOG_LEVEL=log-level
export TIMEOUT=server-timeout  
./search search --index /index/file/path
```

After it you can send request to server
```http request
POST /?search=`search-phrase` HTTP/1.1
Host: `interfase-to-listen`
```

#### Search in docker

Before all run **Elastic Stack**:

```shell script
docker-compose up -f --build docker-compose-logs.yml
```

After you can up search docker-compose:

```shell script
export LOG_LEVEL=`log-level`
export IND_FILE=index-file-name
docker-compose up --build
```

> To start the search in docker, you must have the output folder in the same directory as the docker-compose.yml file.
> Index-file should be in this folder.

You can go in your browser to `localhost` and start searching by web-interface.

To secure your logs [see it.](http://codingfundas.com/setting-up-elasticsearch-6-8-with-kibana-and-x-pack-security-enabled/index.html)